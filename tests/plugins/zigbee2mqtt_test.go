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
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestZigbee2mqtt(t *testing.T) {

	const (
		zigbeeButtonId = "0x00158d00031c8ef3"
		zigbeePlugId   = "0x00158d0001bc27f2"

		buttonSourceScript = `# {"battery":100,"click":"long","duration":1515,"linkquality":126,"voltage":3042}
zigbee2mqttEvent = ->
  #print '---mqtt new event from button---'
  if !message
    return
  payload = JSON.parse(message.payload)
  attrs =
    'battery': payload.battery
    'linkquality': payload.linkquality
    'voltage': payload.voltage
  state = ''
  if payload.action
    attrs.action = payload.action
    state = payload.action + "_action"
  if payload.click
    attrs.click = payload.click
    attrs.action = ""
    state = payload.click + "_click"
  Actor.setState
    'new_state': state.toUpperCase()
    'attribute_values': attrs
`

		plugSourceScript = `# {"consumption":0.07,"linkquality":102,"power":0,"state":"ON","temperature":35,"voltage":240.5}
zigbee2mqttEvent = ->
  #print '---mqtt new event from plug---'
  if !message || message.topic.includes('/set')
    return
  payload = JSON.parse(message.payload)
  attrs =
    'consumption': payload.consumption
    'linkquality': payload.linkquality
    'power': payload.power
    'state': payload.state
    'temperature': payload.temperature
    'voltage': payload.voltage
  Actor.setState
    'new_state': payload.state
    'attribute_values': attrs
`

		plugActionOnOffSourceScript = `
entityAction = (entityId, actionName)->
    #print '---action on/off--'
    topic = entityId.replace(".", "/")+'/set'
    payload = JSON.stringify({"state": actionName})
    Mqtt.publish(topic, payload, 0, false)
`

		task1SourceScript = `
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    if !msg.payload.new_state || !msg.payload.new_state.state
        return false
    return msg.new_state.state.name == 'DOUBLE_CLICK'

automationCondition = (entityId)->
    #print '---condition---'
    entity = Condition.getEntityById('zigbee2mqtt.` + zigbeePlugId + `')
    if !entity
        return false
    if !entity.state || entity.state.name == 'OFF'
        return true
    return false

automationAction = (entityId)->
    #print '---action---'
    Action.callAction('zigbee2mqtt.` + zigbeePlugId + `', 'ON', {})
`

		task2SourceScript = `
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    if !msg.payload.new_state || !msg.payload.new_state.state
        return false
    return msg.new_state.state.name == 'DOUBLE_CLICK'

automationCondition = (entityId)->
    #print '---condition---'
    entity = Condition.getEntityById('zigbee2mqtt.` + zigbeePlugId + `')
    if !entity || !entity.state 
        return false
    if entity.state.name == 'ON'
        return true
    return false

automationAction = (entityId)->
    #print '---action---'
    Action.callAction('zigbee2mqtt.` + zigbeePlugId + `', 'OFF', {})
`
	)

	Convey("zigbee2mqtt", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			So(err, ShouldBeNil)

			err = AddPlugin(adaptors, "zigbee2mqtt")
			ctx.So(err, ShouldBeNil)
			ctx.So(err, ShouldBeNil)
			go mqttServer.Start()

			// add zigbee2mqtt
			zigbeeServer := &m.Zigbee2mqtt{
				Name:       "main",
				Devices:    nil,
				PermitJoin: true,
				BaseTopic:  "zigbee2mqtt",
			}
			zigbeeServer.Id, err = adaptors.Zigbee2mqtt.Add(zigbeeServer)
			So(err, ShouldBeNil)

			// add zigbee2mqtt_device
			butonDevice := &m.Zigbee2mqttDevice{
				Id:            zigbeeButtonId,
				Zigbee2mqttId: zigbeeServer.Id,
				Name:          zigbeeButtonId,
				Model:         "WXKG01LM",
				Description:   "MiJia wireless switch",
				Manufacturer:  "Xiaomi",
				Status:        "active",
				Payload:       []byte("{}"),
			}
			err = adaptors.Zigbee2mqttDevice.Add(butonDevice)
			So(err, ShouldBeNil)

			plugDevice := &m.Zigbee2mqttDevice{
				Id:            zigbeePlugId,
				Zigbee2mqttId: zigbeeServer.Id,
				Name:          zigbeePlugId,
				Model:         "ZNCZ02LM",
				Description:   "Mi power plug ZigBee",
				Manufacturer:  "Xiaomi",
				Status:        "active",
				Payload:       []byte("{}"),
			}
			err = adaptors.Zigbee2mqttDevice.Add(plugDevice)
			So(err, ShouldBeNil)

			// add scripts
			// ------------------------------------------------

			buttonScript, err := AddScript("button", buttonSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			plugScript, err := AddScript("plug", plugSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			plugActionOnOffScript, err := AddScript("plug on", plugActionOnOffSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			task1Script, err := AddScript("task1", task1SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			task2Script, err := AddScript("task2", task2SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			buttonEnt := GetNewButton(fmt.Sprintf("zigbee2mqtt.%s", zigbeeButtonId), []*m.Script{buttonScript})
			err = adaptors.Entity.Add(buttonEnt)
			So(err, ShouldBeNil)

			plugEnt := GetNewPlug(fmt.Sprintf("zigbee2mqtt.%s", zigbeePlugId), []*m.Script{plugScript})
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
			}
			err = adaptors.Entity.Add(plugEnt)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------
			//TASK1
			task1 := &m.Task{
				Name:      "Toggle plug ON",
				Enabled:   true,
				Condition: common.ConditionAnd,
			}
			task1.AddTrigger(&m.Trigger{
				Name:       "",
				EntityId:   &buttonEnt.Id,
				Script:     task1Script,
				PluginName: "state_change",
			})
			task1.AddCondition(&m.Condition{
				Name:   "check plug state",
				Script: task1Script,
			})
			task1.AddAction(&m.Action{
				Name:   "action on Plug",
				Script: task1Script,
			})
			err = adaptors.Task.Add(task1)
			So(err, ShouldBeNil)

			//TASK2
			task2 := &m.Task{
				Name:      "Toggle plug OFF",
				Enabled:   true,
				Condition: common.ConditionAnd,
			}
			task2.AddTrigger(&m.Trigger{
				Name:       "",
				EntityId:   &buttonEnt.Id,
				Script:     task2Script,
				PluginName: "state_change",
			})
			task2.AddCondition(&m.Condition{
				Name:   "check plug state",
				Script: task2Script,
			})
			task2.AddAction(&m.Action{
				Name:   "action on Plug",
				Script: task2Script,
			})
			err = adaptors.Task.Add(task2)
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()
			go zigbee2mqtt.Start()

			defer func() {
				_ = mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				_ = automation.Shutdown()
				pluginManager.Shutdown()
			}()

			//
			// ------------------------------------------------
			mqttCli := mqttServer.NewClient("cli3")
			_ = mqttCli.Subscribe("zigbee2mqtt/"+zigbeePlugId+"/set", func(client mqtt.MqttCli, message mqtt.Message) {
				if string(message.Payload) == `{"state":"ON"}` {
					err = mqttCli.Publish("zigbee2mqtt/"+zigbeePlugId, []byte(`{"consumption":2.87,"energy":2.87,"linkquality":147,"power":0,"state":"ON","temperature":41,"voltage":240.5}`))
				} else {
					err = mqttCli.Publish("zigbee2mqtt/"+zigbeePlugId, []byte(`{"consumption":2.87,"energy":2.87,"linkquality":147,"power":0,"state":"OFF","temperature":41,"voltage":240.5}`))
				}
			})

			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"single","linkquality":131,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"single","linkquality":131,"voltage":3042}`))
			So(err, ShouldBeNil)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"triple","linkquality":132,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"triple","linkquality":132,"voltage":3042}`))
			So(err, ShouldBeNil)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"quadruple","linkquality":133,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"quadruple","linkquality":133,"voltage":3042}`))
			So(err, ShouldBeNil)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"double","linkquality":134,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"double","linkquality":134,"voltage":3042}`))
			So(err, ShouldBeNil)
			time.Sleep(time.Millisecond * 200)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"many","linkquality":134,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"many","linkquality":134,"voltage":3042}`))
			So(err, ShouldBeNil)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"double","linkquality":135,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"double","linkquality":135,"voltage":3042}`))
			So(err, ShouldBeNil)
			time.Sleep(time.Millisecond * 200)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"double","linkquality":135,"voltage":3042}`))
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"double","linkquality":135,"voltage":3042}`))
			So(err, ShouldBeNil)
			time.Sleep(time.Millisecond * 200)

			//wg.Wait()

			time.Sleep(time.Second * 2)
		})
	})
}
