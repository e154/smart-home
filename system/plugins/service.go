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
	"github.com/e154/smart-home/common/web"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
)

type service struct {
	pluginManager common.PluginManager
	bus           bus.Bus
	adaptors      *adaptors.Adaptors
	entityManager entity_manager.EntityManager
	scriptService scripts.ScriptService
	mqttServ      mqtt.MqttServ
	appConfig     *models.AppConfig
	gateClient    *gate_client.GateClient
	scheduler     *scheduler.Scheduler
	crawler       web.Crawler
}

// Plugins ...
func (s service) Plugins() map[string]Plugable {
	list := make(map[string]Plugable)
	pluginList.Range(func(key, value interface{}) bool {
		name := key.(string)
		list[name] = value.(Plugable)
		return true
	})
	return list
}

// PluginManager ...
func (s service) PluginManager() common.PluginManager {
	return s.pluginManager
}

// EventBus ...
func (s service) EventBus() bus.Bus {
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

// Scheduler ...
func (s service) Scheduler() *scheduler.Scheduler {
	return s.scheduler
}

// Crawler ...
func (s service) Crawler() web.Crawler {
	return s.crawler
}
