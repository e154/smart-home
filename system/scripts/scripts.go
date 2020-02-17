// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/scripts/bind"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("scripts")
)

type ScriptService struct {
	cfg        *config.AppConfig
	functions  *Pull
	structures *Pull
}

func NewScriptService(cfg *config.AppConfig) (service *ScriptService) {

	service = &ScriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
	}

	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	return service
}

func (service ScriptService) NewEngine(s *m.Script) (*Engine, error) {
	return NewEngine(s, service.structures, service.functions)
}

func (service *ScriptService) PushStruct(name string, s interface{}) {
	service.structures.Add(name, s)
}

func (service *ScriptService) PushFunctions(name string, s interface{}) {
	service.functions.Add(name, s)
}
