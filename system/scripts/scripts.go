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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/models/devices"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/scripts/bind"
)

var (
	log = common.MustGetLogger("scripts")
)

// ScriptService ...
type ScriptService struct {
	cfg        *config.AppConfig
	functions  *Pull
	structures *Pull
}

// NewScriptService ...
func NewScriptService(cfg *config.AppConfig) (service *ScriptService) {

	service = &ScriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
	}

	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	service.PushFunctions("RunCommand", devices.NewRunCommandBind)
	service.PushFunctions("Zigbee2mqtt", devices.NewZigbee2mqttBind)
	service.PushFunctions("ModBus", devices.NewModBusBind)
	service.PushFunctions("Mqtt", devices.NewMqttBind)
	return service
}

// NewEngine ...
func (service ScriptService) NewEngine(s *m.Script) (*Engine, error) {
	return NewEngine(s, service.structures, service.functions)
}

// PushStruct ...
func (service *ScriptService) PushStruct(name string, s interface{}) {
	service.structures.Add(name, s)
}

// PushFunctions ...
func (service *ScriptService) PushFunctions(name string, s interface{}) {
	service.functions.Add(name, s)
}
