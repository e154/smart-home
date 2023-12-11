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

package supervisor

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/web"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
)

type service struct {
	bus           bus.Bus
	adaptors      *adaptors.Adaptors
	supervisor    Supervisor
	scriptService scripts.ScriptService
	mqttServ      mqtt.MqttServ
	appConfig     *models.AppConfig
	scheduler     *scheduler.Scheduler
	crawler       web.Crawler
}

// Plugins ...
func (s service) Plugins() map[string]Pluggable {
	list := make(map[string]Pluggable)
	pluginList.Range(func(key, value interface{}) bool {
		name := key.(string)
		list[name] = value.(Pluggable)
		return true
	})
	return list
}

// EventBus ...
func (s service) EventBus() bus.Bus {
	return s.bus
}

// Supervisor ...
func (s service) Supervisor() Supervisor {
	return s.supervisor
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

// Scheduler ...
func (s service) Scheduler() *scheduler.Scheduler {
	return s.scheduler
}

// Crawler ...
func (s service) Crawler() web.Crawler {
	return s.crawler
}
