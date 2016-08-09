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

package handlers

import "net/http"

import (
	"github.com/elopio/snapweb-deployer/repo"
	"github.com/elopio/snapweb-deployer/snappy"
)

func (h *Handler) Init(c snappy.SnapdClient, gitURL string, r repo.Snapper) {
	h.snapdClient = c
	h.gitURL = gitURL
	h.repo = r
}

func (h *Handler) MakeMuxer() http.Handler {
	return h.makeMuxer()
}
