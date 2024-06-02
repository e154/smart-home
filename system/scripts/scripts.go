// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

	"github.com/e154/smart-home/common/encryptor"

	"github.com/e154/bus"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts/bind"
	"github.com/e154/smart-home/system/storage"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("scripts")
)

// ScriptService ...
type ScriptService interface {
	NewEngine(s *m.Script) (*Engine, error)
	NewEngineWatcher(*m.Script) (*EngineWatcher, error)
	NewEnginesWatcher([]*m.Script) (*EnginesWatcher, error)
	PushStruct(name string, s interface{})
	PopStruct(name string)
	PushFunctions(name string, s interface{})
	PopFunction(name string)
	Restart()
}

// scriptService ...
type scriptService struct {
	cfg        *m.AppConfig
	functions  *Pull
	structures *Pull
	storage    *storage.Storage
	eventBus   bus.Bus
}

// NewScriptService ...
func NewScriptService(lc fx.Lifecycle,
	cfg *m.AppConfig,
	storage *storage.Storage,
	eventBus bus.Bus) ScriptService {

	s := &scriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
		storage:    storage,
		eventBus:   eventBus,
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
func (s *scriptService) NewEngine(scr *m.Script) (*Engine, error) {
	return NewEngine(scr, s.structures, s.functions)
}

// NewEngineWatcher ...
func (s *scriptService) NewEngineWatcher(script *m.Script) (*EngineWatcher, error) {
	return NewEngineWatcher(script, s, s.eventBus), nil
}

// NewEnginesWatcher ...
func (s *scriptService) NewEnginesWatcher(scripts []*m.Script) (*EnginesWatcher, error) {
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
	s.PushStruct("Log", &bind.LogBind{})
	s.PushFunctions("ExecuteSync", bind.ExecuteSync)
	s.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	s.PushFunctions("Encrypt", encryptor.EncryptBind)
	s.PushFunctions("Decrypt", encryptor.DecryptBind)
	s.PushStruct("Storage", bind.NewStorageBind(s.storage))
	s.PushStruct("http", bind.NewHttpBind())
	s.PushStruct("HTTP", bind.NewHttpBind())
}
