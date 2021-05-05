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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/modbus_rtu"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestModbusRtu(t *testing.T) {

	Convey("node", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "node")
			err = AddPlugin(adaptors, "modbus_rtu")
			So(err, ShouldBeNil)

			go mqttServer.Start()

			// add entity
			// ------------------------------------------------
			nodeEnt := GetNewNode()
			err = adaptors.Entity.Add(nodeEnt)
			So(err, ShouldBeNil)

			modbusEnt := GetNewModbusRtu()
			modbusEnt.ParentId = &nodeEnt.Id
			modbusEnt.Actions = []*m.EntityAction{
				{
					Name:        "ON",
					Description: "включить",
					//Script:      plugActionOnOffScript,
				},
				{
					Name:        "OFF",
					Description: "выключить",
					//Script:      plugActionOnOffScript,
				},
				{
					Name:        "CONDITION CHECK",
					Description: "condition check",
					//Script:      plugActionOnOffScript,
				},
			}
			modbusEnt.States = []*m.EntityState{
				{
					Name:        "ON",
					Description: "on state",
				},
				{
					Name:        "OFF",
					Description: "off state",
				},
				{
					Name:        "ERROR",
					Description: "error state",
				},
			}
			modbusEnt.Attributes[modbus_rtu.AttrSlaveId].Value = 1
			modbusEnt.Attributes[modbus_rtu.AttrBaud].Value = 19200
			modbusEnt.Attributes[modbus_rtu.AttrStopBits].Value = 1
			modbusEnt.Attributes[modbus_rtu.AttrTimeout].Value = 100
			modbusEnt.Attributes[modbus_rtu.AttrDataBits].Value = 8
			modbusEnt.Attributes[modbus_rtu.AttrParity].Value = "none"
			err = adaptors.Entity.Add(modbusEnt)
			So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(m.EntityStorage{
				EntityId:   modbusEnt.Id,
				Attributes: modbusEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)

			time.Sleep(time.Millisecond * 500)

			//...

			mqttServer.Shutdown()
			zigbee2mqtt.Shutdown()
			entityManager.Shutdown()
			automation.Shutdown()
			pluginManager.Shutdown()
		})
	})
}
