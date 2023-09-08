// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
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
	crawler    web.Crawler
	eventBus   bus.Bus
}

// NewScriptService ...
func NewScriptService(lc fx.Lifecycle,
	cfg *m.AppConfig,
	storage *storage.Storage,
	crawler web.Crawler,
	eventBus bus.Bus) ScriptService {

	s := &scriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
		storage:    storage,
		crawler:    crawler,
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
	s.PushStruct("Storage", bind.NewStorageBind(s.storage))
	s.PushStruct("http", bind.NewHttpBind(s.crawler))
	s.PushStruct("HTTP", bind.NewHttpBind(s.crawler))
}
