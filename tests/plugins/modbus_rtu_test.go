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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestModbusRtu(t *testing.T) {

	const plugActionOnOffSourceScript = `

getStatus =(status)->
    if status == 1
        return 'ON'
    else
        return 'OFF'

writeRegisters =(d, c, r)->
    res = ModbusRtu 'WriteMultipleRegisters', d, c, r
    if res.error
        print 'error: ', res.error

checkStatus =->
    COMMAND = []
    FUNC = 'ReadHoldingRegisters'
    ADDRESS = 0
    COUNT = 16

    res = ModbusRtu FUNC, ADDRESS, COUNT, COMMAND
    if res.error
        print 'error: ', res.error
    else 
        print 'ok: ', res.result
        #newStatus = getStatus(res.Result[0])

doOnAction = ->
    writeRegisters(0, 1, [1])
doOffAction = ->
    writeRegisters(0, 1, [0])

entityAction = (entityId, actionName)->
    print '---action on/off--'
    switch actionName
        when 'ON' then doOnAction()
        when 'OFF' then doOffAction()
        when 'CHECK' then checkStatus()
`

	Convey("modbus_rtu", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "node")
			err = AddPlugin(adaptors, "triggers")
			err = AddPlugin(adaptors, "modbus_rtu")
			So(err, ShouldBeNil)

			go mqttServer.Start()

			// add scripts
			// ------------------------------------------------
			plugActionOnOffScript := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "plug script",
				Source:      plugActionOnOffSourceScript,
				Description: "on/off/check condition",
			}

			enginePlugOnOffScript, err := scriptService.NewEngine(plugActionOnOffScript)
			So(err, ShouldBeNil)
			err = enginePlugOnOffScript.Compile()
			So(err, ShouldBeNil)

			plugActionOnOffScript.Id, err = adaptors.Script.Add(plugActionOnOffScript)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			nodeEnt := GetNewNode()
			err = adaptors.Entity.Add(nodeEnt)
			So(err, ShouldBeNil)

			plugEnt := GetNewModbusRtu("plug")
			plugEnt.ParentId = &nodeEnt.Id
			plugEnt.Actions = []*m.EntityAction{
				{
					Name:        "ON",
					Description: "включить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "OFF",
					Description: "выключить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      plugActionOnOffScript,
				},
			}
			plugEnt.States = []*m.EntityState{
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
			plugEnt.Attributes[modbus_rtu.AttrSlaveId].Value = 1
			plugEnt.Attributes[modbus_rtu.AttrBaud].Value = 19200
			plugEnt.Attributes[modbus_rtu.AttrStopBits].Value = 1
			plugEnt.Attributes[modbus_rtu.AttrTimeout].Value = 100
			plugEnt.Attributes[modbus_rtu.AttrDataBits].Value = 8
			plugEnt.Attributes[modbus_rtu.AttrParity].Value = "none"
			err = adaptors.Entity.Add(plugEnt)
			So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(m.EntityStorage{
				EntityId:   plugEnt.Id,
				Attributes: plugEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)

			defer func() {
				mqttServer.Shutdown()
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			entityManager.CallAction(plugEnt.Id, "ON", nil)
			entityManager.CallAction(plugEnt.Id, "OFF", nil)
			entityManager.CallAction(plugEnt.Id, "CHECK", nil)
			entityManager.CallAction(plugEnt.Id, "NULL", nil)

			time.Sleep(time.Millisecond * 500)

			//...
		})
	})
}
