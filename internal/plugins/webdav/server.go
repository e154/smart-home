// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package webdav

import (
	"context"
	"net/http"
	"strings"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/scripts"

	"go.uber.org/atomic"
	"golang.org/x/net/webdav"

	"github.com/e154/bus"
)

const rootDir = "/webdav"

type Server struct {
	adaptors  *adaptors.Adaptors
	eventBus  bus.Bus
	isStarted *atomic.Bool
	scripts   *Scripts
	*FS
	handler *webdav.Handler
}

func NewServer() *Server {
	server := &Server{
		isStarted: atomic.NewBool(false),
	}

	return server
}

func (s *Server) Start(adaptors *adaptors.Adaptors, scriptService scripts.ScriptService, eventBus bus.Bus) {
	if !s.isStarted.CompareAndSwap(false, true) {
		return
	}

	s.adaptors = adaptors
	s.eventBus = eventBus

	s.FS = NewFS(s.onRemoveHandler)
	s.handler = &webdav.Handler{
		FileSystem: s,
		LockSystem: webdav.NewMemLS(),
	}

	_ = s.MkdirAll("/webdav/scripts", 0755)

	s.scripts = NewScripts(s.FS)
	s.scripts.Start(adaptors, scriptService, eventBus)

}

func (s *Server) Shutdown() {
	if !s.isStarted.CompareAndSwap(true, false) {
		return
	}

	s.scripts.Shutdown()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.isStarted.Load() {
		return
	}
	s.handler.ServeHTTP(w, r)
}

func (s *Server) onRemoveHandler(ctx context.Context, filePath string) (err error) {
	path := strings.Split(filePath, "/")
	if path[2] == "scripts" {
		err = s.scripts.onRemoveHandler(ctx, filePath)
	}
	return
}
