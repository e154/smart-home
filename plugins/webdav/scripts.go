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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/net/webdav"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

type Scripts struct {
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	afero.Fs
	rootDir string
	handler *webdav.Handler
	done    chan struct{}
}

func NewScripts() *Scripts {
	return &Scripts{
		Fs:      afero.NewMemMapFs(),
		rootDir: "scripts",
	}
}

func (s *Scripts) Start(adaptors *adaptors.Adaptors, scriptService scripts.ScriptService, eventBus bus.Bus) {

	s.adaptors = adaptors
	s.eventBus = eventBus
	s.scriptService = scriptService

	s.handler = &webdav.Handler{
		FileSystem: s,
		LockSystem: webdav.NewMemLS(),
	}

	s.preload()

	s.done = make(chan struct{})
	go func() {
		for {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()

			select {
			case <-ticker.C:
				s.syncFiles()
			case <-s.done:
				return
			}
		}
	}()

	_ = eventBus.Subscribe("system/models/scripts/+", s.eventHandler)
}

func (s *Scripts) Shutdown() {
	close(s.done)
	_ = s.eventBus.Unsubscribe("system/models/scripts/+", s.eventHandler)
}

func (s *Scripts) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return errors.New("operation not allowed")
}

func (s *Scripts) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	_, err := s.Fs.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Infof("created new file %s", name)
		} else {
			return nil, err
		}
	}
	return s.Fs.OpenFile(name, flag, perm)
}

func (s *Scripts) RemoveAll(ctx context.Context, name string) error {
	return s.Fs.RemoveAll(name)
}

func (s *Scripts) Rename(ctx context.Context, oldName, newName string) error {
	return s.Fs.Rename(oldName, newName)
}

func (s *Scripts) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fileInfo, err := s.Fs.Stat(name)
	if err != nil {
		return nil, err
	}
	return fileInfo, err
}

// eventHandler ...
func (s *Scripts) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventUpdatedScriptModel:
		go s.eventUpdateScript(msg)
	case events.EventRemovedScriptModel:
		go s.eventRemoveScript(msg)
	case events.EventCreatedScriptModel:
		go s.eventAddScript(msg)
	}
}

func (s *Scripts) eventAddScript(msg events.EventCreatedScriptModel) {
	if msg.Owner == events.OwnerSystem {
		return
	}
	filePath := s.getFilePath(msg.Script)
	_ = afero.WriteFile(s.Fs, filePath, []byte(msg.Script.Source), 0644)
	now := time.Now()
	_ = s.Fs.Chtimes(filePath, now, now)
}

func (s *Scripts) eventRemoveScript(msg events.EventRemovedScriptModel) {
	if msg.Owner == events.OwnerSystem {
		return
	}
	filePath := s.getFilePath(msg.Script)
	_ = s.Fs.RemoveAll(filePath)
}

func (s *Scripts) eventUpdateScript(msg events.EventUpdatedScriptModel) {
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

	var recordDir = filepath.Join(rootDir, s.rootDir)

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

func (s *Scripts) extractFileName(path string) string {
	return filepath.Base(path)
}

func (s *Scripts) extractScriptName(path string) string {
	res := strings.Split(path, ".")
	if len(res) > 0 {
		return res[0]
	}
	return path
}

func (s *Scripts) extractScriptLang(path string) common.ScriptLang {
	res := strings.Split(path, ".")
	if len(res) > 1 {
		switch strings.ToLower(res[1]) {
		case "ts":
			return "ts"
		case "js":
			return "javascript"
		case "coffee":
			return "coffeescript"
		}
	}
	return ""
}

func (s *Scripts) CreateOrUpdateFile(ctx context.Context, name string, fileInfo os.FileInfo) (err error) {
	scriptName := s.extractScriptName(fileInfo.Name())
	lang := s.extractScriptLang(fileInfo.Name())

	if lang == "" {
		err = errors.New("bad extension")
		return
	}

	var source []byte
	source, err = afero.ReadFile(s.Fs, name)
	if err != nil {
		return
	}

	var script *m.Script
	script, err = s.adaptors.Script.GetByName(ctx, scriptName)
	if err == nil {
		script.Source = string(source)
		script.Lang = lang

		var engine *scripts.Engine
		engine, err = s.scriptService.NewEngine(script)
		if err != nil {
			return
		}
		if err = engine.Compile(); err != nil {
			return
		}
		err = s.adaptors.Script.Update(ctx, script)
		return
	}

	script = &m.Script{
		Name:   scriptName,
		Lang:   lang,
		Source: string(source),
	}
	engine, err := s.scriptService.NewEngine(script)
	if err != nil {
		return
	}
	if err = engine.Compile(); err != nil {
		return
	}

	if _, err = s.adaptors.Script.Add(ctx, script); err != nil {
		return
	}
	return
}

func (s *Scripts) syncFiles() {
	afero.Walk(s.Fs, "/webdav/scripts", func(path string, fileInfo os.FileInfo, err error) error {
		if time.Now().Sub(fileInfo.ModTime()).Seconds() > 5 {
			return err
		}
		if err := s.CreateOrUpdateFile(context.Background(), path, fileInfo); err != nil {
			log.Error(err.Error())
		}
		return nil
	})
}
