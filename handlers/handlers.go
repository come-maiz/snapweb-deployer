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

// Package handlers implements the service endpoints.
package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snapcore/snapd/client"

	"github.com/elopio/snapweb-deployer/repo"
	"github.com/elopio/snapweb-deployer/snappy"
)

// InitURLHandlers initialize the URL handlers for the deployer service.
func InitURLHandlers(log *log.Logger) {
	log.Println("Initializing HTTP handlers...")

	handler := NewHandler("https://github.com/snapcore/snapweb.git")
	http.Handle("/", handler.makeMuxer())
}

// A Handler listens and servers the requests to the service.
type Handler struct {
	snapdClient snappy.SnapdClient
	gitURL      string
	repo        repo.Snapper
}

// NewHandler returns a service handler.
func NewHandler(gitURL string) *Handler {
	return &Handler{
		snapdClient: client.New(nil),
		gitURL:      gitURL,
		repo:        &repo.Repo{},
	}
}

func (h *Handler) makeMuxer() http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/{prID}", h.deploy).Methods("PUT")
	m.HandleFunc("/{prID}", h.remove).Methods("DELETE")

	return m
}

func (h *Handler) deploy(w http.ResponseWriter, r *http.Request) {
	prID := mux.Vars(r)["prID"]

	snapPath, err := h.repo.GetSnap(h.gitURL, prID)
	if err != nil {

	}
	_, err = h.snapdClient.InstallPath(snapPath, nil)
	if err != nil {

	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	prID := mux.Vars(r)["prID"]
	snapName := fmt.Sprintf("snapweb-%s", prID)
	_, err := h.snapdClient.Remove(snapName, nil)
	if err != nil {

	}
	w.WriteHeader(http.StatusAccepted)
}
