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
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/modbus_rtu"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/plugins/scene"
	"github.com/e154/smart-home/plugins/script"
	"github.com/e154/smart-home/plugins/zigbee2mqtt"
	"github.com/e154/smart-home/plugins/zone"
)

func GetNewButton(id string, scripts []m.Script) *m.Entity {
	return &m.Entity{
		Id:          common.EntityId(id),
		Description: "MiJia wireless switch",
		Type:        zigbee2mqtt.EntityZigbee2mqtt,
		Scripts:     scripts,
		AutoLoad:    true,
		Attributes: m.EntityAttributes{
			"click": &m.EntityAttribute{
				Name: "click",
				Type: common.EntityAttributeString,
			},
			"action": &m.EntityAttribute{
				Name: "action",
				Type: common.EntityAttributeString,
			},
			"battery": &m.EntityAttribute{
				Name: "battery",
				Type: common.EntityAttributeInt,
			},
			"voltage": &m.EntityAttribute{
				Name: "voltage",
				Type: common.EntityAttributeInt,
			},
			"linkquality": &m.EntityAttribute{
				Name: "linkquality",
				Type: common.EntityAttributeInt,
			},
		},
		States: []*m.EntityState{
			{
				Name:        "LONG_CLICK",
				Description: "long click",
			},
			{
				Name:        "LONG_ACTION",
				Description: "long action",
			},
			{
				Name:        "SINGLE_CLICK",
				Description: "single click",
			},
			{
				Name:        "SINGLE_ACTION",
				Description: "single action",
			},
			{
				Name:        "DOUBLE_CLICK",
				Description: "double click",
			},
			{
				Name:        "DOUBLE_ACTION",
				Description: "double action",
			},
			{
				Name:        "TRIPLE_CLICK",
				Description: "triple click",
			},
			{
				Name:        "TRIPLE_ACTION",
				Description: "triple action",
			},
			{
				Name:        "QUADRUPLE_CLICK",
				Description: "quadruple click",
			},
			{
				Name:        "QUADRUPLE_ACTION",
				Description: "quadruple action",
			},
			{
				Name:        "MANY_CLICK",
				Description: "many click",
			},
			{
				Name:        "MANY_ACTION",
				Description: "many action",
			},
			{
				Name:        "LONG_RELEASE_CLICK",
				Description: "long_release click",
			},
			{
				Name:        "HOLD_ACTION",
				Description: "hold action",
			},
			{
				Name:        "RELEASE_ACTION",
				Description: "release action",
			},
		},
	}
}

func GetNewPlug(id string, scrits []m.Script) *m.Entity {
	return &m.Entity{
		Id:          common.EntityId(id),
		Description: "MiJia power plug ZigBee",
		Type:        zigbee2mqtt.EntityZigbee2mqtt,
		Scripts:     scrits,
		AutoLoad:    true,
		Attributes: m.EntityAttributes{
			"power": &m.EntityAttribute{
				Name: "power",
				Type: common.EntityAttributeInt,
			},
			"state": &m.EntityAttribute{
				Name: "state",
				Type: common.EntityAttributeString,
			},
			"voltage": &m.EntityAttribute{
				Name: "voltage",
				Type: common.EntityAttributeInt,
			},
			"consumption": &m.EntityAttribute{
				Name: "consumption",
				Type: common.EntityAttributeString,
			},
			"linkquality": &m.EntityAttribute{
				Name: "linkquality",
				Type: common.EntityAttributeInt,
			},
			"temperature": &m.EntityAttribute{
				Name: "temperature",
				Type: common.EntityAttributeInt,
			},
		},
		States: []*m.EntityState{
			{
				Name:        "ON",
				Description: "on state",
			},
			{
				Name:        "OFF",
				Description: "off state",
			},
		},
	}
}

func GetNewScript(id string, scrits []m.Script) *m.Entity {
	return &m.Entity{
		Id:          common.EntityId(id),
		Description: "MiJia power plug ZigBee",
		Type:        script.EntityScript,
		Scripts:     scrits,
		Attributes:  m.EntityAttributes{},
		AutoLoad:    true,
		States: []*m.EntityState{
			{
				Name:        "ON",
				Description: "on state",
			},
			{
				Name:        "OFF",
				Description: "off state",
			},
		},
	}
}

func GetNewScene(id string, scrits []m.Script) *m.Entity {
	return &m.Entity{
		Id:          common.EntityId(id),
		Description: "scene",
		Type:        scene.EntityScene,
		Scripts:     scrits,
		AutoLoad:    true,
	}
}

func GetNewZone() *m.Entity {
	return &m.Entity{
		Id:          "zone.home",
		Description: "main geo zone",
		Type:        "zone",
		AutoLoad:    true,
		Attributes: m.EntityAttributes{
			zone.AttrLat: {
				Name:  zone.AttrLat,
				Type:  common.EntityAttributeFloat,
				Value: 54.9022,
			},
			zone.AttrLon: {
				Name:  zone.AttrLon,
				Type:  common.EntityAttributeFloat,
				Value: 83.0335,
			},
			zone.AttrElevation: {
				Name:  zone.AttrElevation,
				Type:  common.EntityAttributeFloat,
				Value: 150,
			},
			zone.AttrTimezone: {
				Name:  zone.AttrTimezone,
				Type:  common.EntityAttributeInt,
				Value: 7,
			},
		},
	}
}

func GetNewNode() *m.Entity {
	return &m.Entity{
		Id:          "node.main",
		Description: "main node",
		Type:        "node",
		AutoLoad:    true,
		Attributes:  node.NewAttr(),
	}
}

func GetNewModbusRtu(name string) *m.Entity {
	return &m.Entity{
		Id:          common.EntityId(fmt.Sprintf("modbus_rtu.%s", name)),
		Description: fmt.Sprintf("%s entity", name),
		Type:        "modbus_rtu",
		AutoLoad:    true,
		Attributes:  modbus_rtu.NewAttr(),
	}
}

func AddPlugin(adaptors *adaptors.Adaptors, name string) (err error) {
	err = adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: true,
		System:  true,
	})
	return
}
