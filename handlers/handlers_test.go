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

package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elopio/snapweb-deployer/handlers"

	"github.com/snapcore/snapd/client"
)

type fakeSnapdClient struct {
	installed []string
	err       error
}

func (c *fakeSnapdClient) InstallPath(path string, options *client.SnapOptions) (string, error) {
	c.installed = append(c.installed, path)
	return "", nil
}

func (c *fakeSnapdClient) Remove(name string, options *client.SnapOptions) (string, error) {
	for i, v := range c.installed {
		if v == name {
			c.installed = append(c.installed[:i], c.installed[i+1:]...)
			break
		}
	}
	return "", nil
}

type fakeRepo struct {
}

func (r *fakeRepo) GetSnap(gitURL string, prID string) (string, error) {
	return prID, nil
}

func TestAcceptedPut(t *testing.T) {
	handler := &handlers.Handler{}
	client := &fakeSnapdClient{}
	handler.Init(client, "dummy", &fakeRepo{})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/test-pr-id", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	handler.MakeMuxer().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("wrong return code: expected %v, got %v", http.StatusAccepted, rec.Code)
	}
	if len(client.installed) != 1 || client.installed[0] != "test-pr-id" {
		t.Fatalf("snap from test-pr-id was not installed: installed snaps: %v", client.installed)
	}
}

func TestAcceptedDelete(t *testing.T) {
	handler := &handlers.Handler{}
	client := &fakeSnapdClient{}
	client.installed = append(client.installed, "snapweb-test-pr-id")
	handler.Init(client, "dummy", &fakeRepo{})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/test-pr-id", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	handler.MakeMuxer().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("wrong return code: expected %v, got %v", http.StatusAccepted, rec.Code)
	}
	if len(client.installed) != 0 {
		t.Fatalf("snap from test-pr-id was not removed: installed snaps: %v", client.installed)
	}
}
