// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

package scripts

import (
	"context"

	"github.com/e154/smart-home/internal/system/scripts/bind"
	"github.com/e154/smart-home/internal/system/storage"
	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common/encryptor"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("scripts")
)

// scriptService ...
type scriptService struct {
	cfg        *models.AppConfig
	functions  *Pull
	structures *Pull
	storage    *storage.Storage
	eventBus   bus.Bus
	adaptors   *adaptors.Adaptors
	validation *validation.Validate
}

// NewScriptService ...
func NewScriptService(lc fx.Lifecycle,
	cfg *models.AppConfig,
	storage *storage.Storage,
	eventBus bus.Bus,
	adaptors *adaptors.Adaptors,
	validation *validation.Validate) scripts.ScriptService {

	s := &scriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
		storage:    storage,
		eventBus:   eventBus,
		adaptors:   adaptors,
		validation: validation,
	}

	s.bind()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			eventBus.Publish("system/services/scripts", events.EventServiceStarted{Service: "Scripts"})
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			eventBus.Publish("system/services/scripts", events.EventServiceStopped{Service: "Scripts"})
			return
		},
	})

	return s
}

// NewEngine ...
func (s *scriptService) NewEngine(scr *models.Script) (scripts.Engine, error) {
	return NewEngine(scr, s.structures, s.functions, s.SourceLoader)
}

// NewEngineWatcher ...
func (s *scriptService) NewEngineWatcher(script *models.Script) (scripts.EngineWatcher, error) {
	return NewEngineWatcher(script, s, s.eventBus), nil
}

// NewEnginesWatcher ...
func (s *scriptService) NewEnginesWatcher(scripts []*models.Script) (scripts.EnginesWatcher, error) {
	return NewEnginesWatcher(scripts, s, s.eventBus), nil
}

// PushStruct ...
func (s *scriptService) PushStruct(name string, str interface{}) {
	log.Infof("register structure: '%s'", name)
	s.structures.Push(name, str)
}

// PopStruct ...
func (s *scriptService) PopStruct(name string) {
	log.Infof("unregister structure: '%s'", name)
	s.structures.Pop(name)
}

// PushFunctions ...
func (s *scriptService) PushFunctions(name string, f interface{}) {
	log.Infof("register function: '%s'", name)
	s.functions.Push(name, f)
}

// PopFunction ...
func (s *scriptService) PopFunction(name string) {
	log.Infof("unregister function: '%s'", name)
	s.functions.Pop(name)
}

// Restart ...
func (s *scriptService) Restart() {
	s.functions.Purge()
	s.structures.Purge()
	s.bind()
	s.eventBus.Publish("system/services/scripts", events.EventServiceRestarted{})
}

func (s *scriptService) bind() {
	s.PushFunctions("ExecuteSync", bind.ExecuteSync)
	s.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	s.PushFunctions("Encrypt", encryptor.EncryptBind)
	s.PushFunctions("Decrypt", encryptor.DecryptBind)
	s.PushStruct("Storage", bind.NewStorageBind(s.storage))
	s.PushStruct("Variable", bind.NewVariable(s.adaptors, s.validation, s.eventBus))
	s.PushStruct("http", bind.NewHttpBind())
	s.PushStruct("HTTP", bind.NewHttpBind())
}

func (s *scriptService) SourceLoader(path string) ([]byte, error) {
	script, err := s.adaptors.Script.GetByName(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return []byte(script.Compiled), nil
}
