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

package _default

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

// PluginManager ...
type PluginManager struct {
	adaptors *adaptors.Adaptors
}

// NewPluginManager ...
func NewPluginManager(adaptors *adaptors.Adaptors) *PluginManager {
	return &PluginManager{
		adaptors: adaptors,
	}
}

func (n PluginManager) addPlugin(name string, enabled, system, actor bool) (node *m.Plugin) {
	_ = n.adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
	return
}

// Create ...
func (n PluginManager) Create() (home *m.Plugin) {
	n.addPlugin("alexa", false, false, false)
	n.addPlugin("cgminer", true, false, true)
	n.addPlugin("cpuspeed", false, false, false)
	n.addPlugin("email", true, false, false)
	n.addPlugin("messagebird", false, false, false)
	n.addPlugin("modbus_rtu", false, false, true)
	n.addPlugin("modbus_tcp", false, false, true)
	n.addPlugin("moon", false, false, true)
	n.addPlugin("node", true, true, true)
	n.addPlugin("notify", true, true, false)
	n.addPlugin("scene", true, false, true)
	n.addPlugin("script", true, false, true)
	n.addPlugin("sensor", true, false, true)
	n.addPlugin("slack", false, false, false)
	n.addPlugin("sun", false, false, true)
	n.addPlugin("telegram", true, false, true)
	n.addPlugin("triggers", true, true, false)
	n.addPlugin("twilio", false, false, false)
	n.addPlugin("updater", false, false, false)
	n.addPlugin("uptime", false, false, false)
	n.addPlugin("weather", false, false, false)
	n.addPlugin("weather_met", false, false, true)
	n.addPlugin("weather_owm", false, false, true)
	n.addPlugin("zigbee2mqtt", false, false, true)
	n.addPlugin("zone", false, false, true)
	return
}
