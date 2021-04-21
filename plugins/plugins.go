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

package plugins

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/plugins/cpuspeed"
	"github.com/e154/smart-home/plugins/moon"
	"github.com/e154/smart-home/plugins/scene"
	"github.com/e154/smart-home/plugins/script"
	"github.com/e154/smart-home/plugins/sun"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/plugins/updater"
	"github.com/e154/smart-home/plugins/uptime"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/plugins/weather_met"
	"github.com/e154/smart-home/plugins/zigbee2mqtt"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/plugin_manager"
	"github.com/e154/smart-home/system/scripts"
)

type Loader struct {
	pluginManager *plugin_manager.PluginManager
	entityManager *entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	eventBus      *event_bus.EventBus
	mqtt          *mqtt.Mqtt
	scriptService *scripts.ScriptService
}

func NewPluginLoader(
	pluginManager *plugin_manager.PluginManager,
	entityManager *entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	eventBus *event_bus.EventBus,
	mqtt *mqtt.Mqtt,
	scriptService *scripts.ScriptService) *Loader {
	plugins := &Loader{
		pluginManager: pluginManager,
		entityManager: entityManager,
		adaptors:      adaptors,
		eventBus:      eventBus,
		mqtt:          mqtt,
		scriptService: scriptService,
	}
	return plugins
}

// Register ...
func (p *Loader) Register() {

	triggers.Register(p.pluginManager, p.eventBus)
	zigbee2mqtt.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors, p.mqtt, p.scriptService)
	script.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors, p.scriptService)
	scene.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors, p.scriptService)
	updater.Register(p.pluginManager, p.entityManager, 24)
	uptime.Register(p.pluginManager, p.entityManager, p.adaptors, 60)
	zone.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors)
	sun.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors, 240)
	moon.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors, 240)
	weather.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors)
	weather_met.Register(p.pluginManager, p.entityManager, p.eventBus, p.adaptors)
	cpuspeed.Register(p.pluginManager, p.entityManager, p.adaptors, 5)
}
