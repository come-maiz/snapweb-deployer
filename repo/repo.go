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

package repo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/elopio/snapweb-deployer/snapcraft"
)

// Snapper builds snaps from a git repository pull request.
type Snapper interface {
	GetSnap(gitURL string, prID string) (string, error)
}

// Repo is a git repository.
type Repo struct {
	snapcrafter snapcraft.Snapcrafter
}

// GetSnap returns a snapweb snap built from a pull request.
// It patches the snap to have the pull request ID in the name and the
// port to make it unique.
func (r *Repo) GetSnap(gitURL string, prID string) (string, error) {
	repoPath, err := clone(gitURL, prID)
	if err != nil {
		return "", err
	}
	patchedPath, err := patch(repoPath, prID)
	if err != nil {
		return "", err
	}
	return r.snapcrafter.Snap(patchedPath)
}

func clone(gitURL string, prID string) (string, error) {
	repoDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	err = exec.Command("git", "clone", gitURL, repoDir).Run()
	if err != nil {
		return "", err
	}
	fetch := exec.Command("git", "fetch", "origin", fmt.Sprintf("pull/%s/merge:pull_%s", prID, prID))
	fetch.Dir = repoDir
	err = fetch.Run()
	if err != nil {
		return "", err
	}
	checkout := exec.Command("git", "checkout", fmt.Sprintf("pull_%s", prID))
	checkout.Dir = repoDir
	err = checkout.Run()
	if err != nil {
		return "", err
	}
	return repoDir, nil
}

func patch(repoPath string, prID string) (string, error) {
	yamlPath := path.Join(repoPath, "snapcraft.yaml")
	yamlSed := exec.Command("sed", "-i", fmt.Sprintf("s/name: snapweb/name: snapweb-%s/", prID), yamlPath)
	yamlSed.Dir = repoPath
	err := yamlSed.Run()
	if err != nil {
		return "", err
	}
	mainPath := path.Join(repoPath, "cmd", "snapweb", "main.go")
	mainSed := exec.Command("sed", "-i", fmt.Sprintf("s/:4200/:4%s/", prID), mainPath)
	mainSed.Dir = repoPath
	err = mainSed.Run()
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
