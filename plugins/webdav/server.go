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
	"fmt"
	"net/http"
	"os"

	"go.uber.org/atomic"
	"golang.org/x/net/webdav"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/bus"
)

const rootDir = "/webdav"

type Server struct {
	adaptors *adaptors.Adaptors
	eventBus bus.Bus
	*FS
	handler   *webdav.Handler
	isStarted *atomic.Bool
	scripts   *Scripts
}

func NewServer() *Server {
	server := &Server{
		isStarted: atomic.NewBool(false),
	}
	server.FS = NewFS(server.onFileCreated,
		server.onFileUpdated,
		server.onFileRemoved,
		server.onFileRenamed,
		server.onOpenFile,
	)
	return server
}

func (s *Server) Start(adaptors *adaptors.Adaptors, eventBus bus.Bus) {
	if !s.isStarted.CompareAndSwap(false, true) {
		return
	}

	s.adaptors = adaptors
	s.eventBus = eventBus

	s.scripts = NewScripts(s.FS)
	s.scripts.Start(adaptors, eventBus)

	s.handler = &webdav.Handler{
		FileSystem: s.FS,
		LockSystem: webdav.NewMemLS(),
	}
}

func (s *Server) Shutdown() {
	if !s.isStarted.CompareAndSwap(true, false) {
		return
	}

	s.handler = nil
	s.scripts.Shutdown()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.isStarted.Load() {
		return
	}
	s.handler.ServeHTTP(w, r)
}

func (s *Server) onFileCreated(fileInfo os.FileInfo) {
	fmt.Println("Файл был недавно создан.", fileInfo.Name())
}

func (s *Server) onFileUpdated(fileInfo os.FileInfo) {
	fmt.Println("Файл был изменен.", fileInfo.Name())
}

func (s *Server) onFileRenamed(oldName, newName string) {
	fmt.Println("Файл был переименован.", oldName, newName)
}

func (s *Server) onFileRemoved(name string) {
	fmt.Println("Файл был удален.", name)
}

func (s *Server) onOpenFile(name string, flag int, perm os.FileMode) {
	fmt.Println("Файл был открыт.", name, flag, perm)
}
