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

func (n PluginManager) addPlugin(name string, enabled bool) (node *m.Plugin) {
	n.adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: enabled,
		System:  true,
	})
	return
}

// Create ...
func (n PluginManager) Create() (home *m.Plugin) {
	n.addPlugin("cpuspeed", false)
	n.addPlugin("moon", true)
	n.addPlugin("scene", true)
	n.addPlugin("script", true)
	n.addPlugin("sun", true)
	n.addPlugin("zone", true)
	n.addPlugin("triggers", true)
	n.addPlugin("update", true)
	n.addPlugin("uptime", true)
	n.addPlugin("weather", true)
	n.addPlugin("weather_met", true)
	n.addPlugin("weather_owm", true)
	n.addPlugin("zigbee2mqtt", true)
	n.addPlugin("zone", true)
	n.addPlugin("node", true)
	n.addPlugin("modbus_rtu", true)
	n.addPlugin("modbus_tcp", true)
	n.addPlugin("alexa", true)
	n.addPlugin("notify", true)
	n.addPlugin("email", true)
	n.addPlugin("slack", true)
	n.addPlugin("cgminer", true)
	n.addPlugin("telegram", true)
	n.addPlugin("sensor", true)
	return
}
