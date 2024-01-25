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
	"path/filepath"
	"time"

	"github.com/spf13/afero"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
)

type Scripts struct {
	adaptors *adaptors.Adaptors
	eventBus bus.Bus
	*FS
	rootDir string
}

func NewScripts(fs *FS) *Scripts {
	return &Scripts{
		FS:      fs,
		rootDir: "scripts",
	}
}

func (s *Scripts) Start(adaptors *adaptors.Adaptors, eventBus bus.Bus) {

	s.adaptors = adaptors
	s.eventBus = eventBus

	s.preload()

	_ = eventBus.Subscribe("system/models/scripts/+", s.eventHandler)
}

func (s *Scripts) Shutdown() {
	_ = s.eventBus.Unsubscribe("system/models/scripts/+", s.eventHandler)
}

// eventHandler ...
func (s *Scripts) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventUpdatedScriptModel:
		go s.updateScript(msg)
	case events.EventRemovedScriptModel:
		go s.removeScript(msg)
	case events.EventCreatedScriptModel:
		go s.addScript(msg)
	}
}

func (s *Scripts) addScript(msg events.EventCreatedScriptModel) {
	if msg.Owner == events.OwnerSystem {
		return
	}
	filePath := s.getFilePath(msg.Script)
	_ = afero.WriteFile(s.Fs, filePath, []byte(msg.Script.Source), 0644)
	now := time.Now()
	_ = s.Fs.Chtimes(filePath, now, now)
}

func (s *Scripts) removeScript(msg events.EventRemovedScriptModel) {
	if msg.Owner == events.OwnerSystem {
		return
	}
	filePath := s.getFilePath(msg.Script)
	_ = s.Fs.RemoveAll(filePath)
}

func (s *Scripts) updateScript(msg events.EventUpdatedScriptModel) {
	if msg.Owner == events.OwnerSystem {
		return
	}
	now := time.Now()
	filePath := s.getFilePath(msg.Script)
	_ = afero.WriteFile(s.Fs, filePath, []byte(msg.Script.Source), 0644)
	_ = s.Fs.Chtimes(filePath, now, now)
}

func (s *Scripts) preload() {
	log.Info("Load script list")

	var recordDir = filepath.Join(rootDir, rootDir)

	_ = s.Fs.MkdirAll(recordDir, 0755)

	var page int64
	var scripts []*m.Script
	const perPage = 500
	var err error

LOOP:

	if scripts, _, err = s.adaptors.Script.List(context.Background(), perPage, perPage*page, "desc", "id", nil, nil); err != nil {
		log.Error(err.Error())
		return
	}

	for _, script := range scripts {
		filePath := s.getFilePath(script)
		if err = afero.WriteFile(s.Fs, filePath, []byte(script.Source), 0644); err != nil {
			log.Error(err.Error())
		}
		if err = s.Fs.Chtimes(filePath, script.CreatedAt, script.UpdatedAt); err != nil {
			log.Error(err.Error())
		}
	}

	if len(scripts) != 0 {
		page++
		goto LOOP
	}
}

func (s *Scripts) getFilePath(script *m.Script) string {
	return filepath.Join(rootDir, s.rootDir, getFileName(script))
}
