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
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts/bind"
	"github.com/e154/smart-home/system/storage"
)

var (
	log = logger.MustGetLogger("scripts")
)

// ScriptService ...
type ScriptService interface {
	NewEngine(s *m.Script) (*Engine, error)
	PushStruct(name string, s interface{})
	PushFunctions(name string, s interface{})
	Purge()
}

// scriptService ...
type scriptService struct {
	cfg        *m.AppConfig
	functions  *Pull
	structures *Pull
	storage    *storage.Storage
}

// NewScriptService ...
func NewScriptService(cfg *m.AppConfig, storage *storage.Storage) (service ScriptService) {

	service = &scriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
		storage:    storage,
	}

	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	service.PushStruct("Storage", bind.NewStorageBind(storage))
	service.PushStruct("http", &bind.HttpBind{})
	return service
}

// NewEngine ...
func (service *scriptService) NewEngine(s *m.Script) (*Engine, error) {
	return NewEngine(s, service.structures, service.functions)
}

// PushStruct ...
func (service *scriptService) PushStruct(name string, s interface{}) {
	log.Infof("register structure: '%s'", name)
	service.structures.Add(name, s)
}

// PushFunctions ...
func (service *scriptService) PushFunctions(name string, s interface{}) {
	log.Infof("register function: '%s'", name)
	service.functions.Add(name, s)
}

// Purge ...
func (service *scriptService) Purge() {
	service.functions.Purge()
	service.structures.Purge()
	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	service.PushStruct("Storage", bind.NewStorageBind(service.storage))
	service.PushStruct("http", &bind.HttpBind{})
}
