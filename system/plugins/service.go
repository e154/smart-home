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

package plugins

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
)

type service struct {
	pluginManager common.PluginManager
	bus           event_bus.EventBus
	adaptors      *adaptors.Adaptors
	entityManager entity_manager.EntityManager
	scriptService scripts.ScriptService
	mqttServ      mqtt.MqttServ
	appConfig     *models.AppConfig
	gateClient    *gate_client.GateClient
}

// Plugins ...
func (s service) Plugins() map[string]Plugable {
	return pluginList
}

// PluginManager ...
func (s service) PluginManager() common.PluginManager {
	return s.pluginManager
}

// EventBus ...
func (s service) EventBus() event_bus.EventBus {
	return s.bus
}

// EntityManager ...
func (s service) EntityManager() entity_manager.EntityManager {
	return s.entityManager
}

// Adaptors ...
func (s service) Adaptors() *adaptors.Adaptors {
	return s.adaptors
}

// ScriptService ...
func (s service) ScriptService() scripts.ScriptService {
	return s.scriptService
}

// MqttServ ...
func (s service) MqttServ() mqtt.MqttServ {
	return s.mqttServ
}

// AppConfig ...
func (s service) AppConfig() *models.AppConfig {
	return s.appConfig
}

// GateClient ...
func (s service) GateClient() *gate_client.GateClient {
	return s.gateClient
}
