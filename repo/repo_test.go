/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package repo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/elopio/snapweb-deployer/repo"
)

type fakeSnapcrafter struct {
}

// The fake snap method returns the path to the patched repository before it's snapped.
func (s *fakeSnapcrafter) Snap(path string) (string, error) {
	return path, nil
}

func prepareTestRepo(prID string, snapName string, snapPort string, commitMsg string) (string, error) {
	repoDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir: %v", err)
	}
	err = exec.Command("git", "init", repoDir).Run()
	if err != nil {
		return "", fmt.Errorf("error initializing the test repo: %v", err)
	}
	branch := exec.Command("git", "checkout", "-b", fmt.Sprintf("pull/%s/merge", prID))
	branch.Dir = repoDir
	err = branch.Run()
	if err != nil {
		return "", fmt.Errorf("error creating pull request branch: %v", err)
	}

	yamlFile, err := os.Create(path.Join(repoDir, "snapcraft.yaml"))
	if err != nil {
		return "", fmt.Errorf("error creating yaml file: %v", err)
	}
	defer yamlFile.Close()
	yamlFile.WriteString(fmt.Sprintf("name: %s", snapName))
	snapwebCmdDir := path.Join(repoDir, "cmd", "snapweb")
	os.MkdirAll(snapwebCmdDir, os.ModePerm)
	mainFile, err := os.Create(path.Join(snapwebCmdDir, "main.go"))
	if err != nil {
		return "", fmt.Errorf("error creating main file: %v", err)
	}
	defer mainFile.Close()
	mainFile.WriteString(fmt.Sprintf(":%s", snapPort))

	add := exec.Command("git", "add", "*")
	add.Dir = repoDir
	err = add.Run()
	if err != nil {
		return "", fmt.Errorf("error adding files to the repository: %v", err)
	}

	commit := exec.Command("git", "commit", "--allow-empty", "-m", commitMsg)
	commit.Dir = repoDir
	err = commit.Run()
	if err != nil {
		return "", fmt.Errorf("error during commit: %v", err)
	}
	return repoDir, nil
}

func TestGetSnapClonesRepoPR(t *testing.T) {
	repo := &repo.Repo{}
	repo.Init(&fakeSnapcrafter{})

	testCommit := "test commit"
	repoDir, err := prepareTestRepo("test-prid", "snapweb", "dummy", testCommit)
	if err != nil {
		t.Fatal(err)
	}
	path, err := repo.GetSnap(repoDir, "test-prid")
	if err != nil {
		t.Fatalf("error on GetSnap: %v", err)
	}
	getOrigin := exec.Command("git", "config", "--get", "remote.origin.url")
	getOrigin.Dir = path
	origin, err := getOrigin.Output()
	if err != nil {
		t.Fatalf("error getting the origin: %v", err)
	}
	originStr := strings.TrimSpace(string(origin))
	if originStr != repoDir {
		t.Fatalf("wrong origin: expected %v, got %v", repoDir, originStr)
	}
	getCommitMsg := exec.Command("git", "log", "-1", "--pretty=%B")
	getCommitMsg.Dir = path
	commitMsg, err := getCommitMsg.Output()
	if err != nil {
		t.Fatalf("error getting the last commit: %v", err)
	}
	commitStr := strings.TrimSpace(string(commitMsg))
	if commitStr != testCommit {
		t.Fatalf("wrong branch: expected last commit %v, got %v", testCommit, commitStr)
	}
}

func TestGetSnapPatchesRepo(t *testing.T) {
	repo := &repo.Repo{}
	repo.Init(&fakeSnapcrafter{})

	repoDir, err := prepareTestRepo("test-prid", "snapweb", "4200", "dummy")
	if err != nil {
		t.Fatal(err)
	}

	patchedPath, err := repo.GetSnap(repoDir, "test-prid")
	if err != nil {
		t.Fatalf("error on GetSnap: %v", err)
	}

	yaml, err := ioutil.ReadFile(path.Join(patchedPath, "snapcraft.yaml"))
	if err != nil {
		t.Fatalf("error reading yaml file: %v", err)
	}
	expectedPatchedName := "name: snapweb-test-prid"
	yamlStr := string(yaml)
	if yamlStr != expectedPatchedName {
		t.Fatalf("wrong patched name: expected %v, got %v", expectedPatchedName, yamlStr)
	}

	main, err := ioutil.ReadFile(path.Join(patchedPath, "cmd", "snapweb", "main.go"))
	if err != nil {
		t.Fatalf("error reading main file: %v", err)
	}
	expectedPatchedPort := ":4test-prid"
	mainStr := string(main)
	if mainStr != expectedPatchedPort {
		t.Fatalf("wrong patched port: expected %v, got %v", expectedPatchedPort, mainStr)
	}
}
