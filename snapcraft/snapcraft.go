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

package snapcraft

import (
	"os/exec"
	"path"
	"path/filepath"
)

// Snapcrafter is an interface to execute snapcraft commands.
type Snapcrafter interface {
	Snap(path string) (string, error)
}

// Snapcraft executes snapcraft commands.
type Snapcraft struct {
}

// Snap builds a snap using the snapcraft tool.
// It returns the path to the snap file.
func (s *Snapcraft) Snap(repoPath string) (string, error) {
	snapcraft := exec.Command("snapcraft")
	snapcraft.Dir = repoPath
	err := snapcraft.Run()
	if err != nil {
		return "", err
	}
	matches, err := filepath.Glob(path.Join(repoPath, "*.snap"))
	if err != nil {
		return "", err
	}
	return matches[0], nil
}
