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

package zigbee2mqtt

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/e154/smart-home/tests/plugins"
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
  payload = unmarshal message.payload
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
  SetState ENTITY_ID,
    'new_state': state.toUpperCase()
    'attribute_values': attrs
`

		plugSourceScript = `# {"consumption":0.07,"linkquality":102,"power":0,"state":"ON","temperature":35,"voltage":240.5}
zigbee2mqttEvent = ->
  #print '---mqtt new event from plug---'
  if !message || message.topic.includes('/set')
    return
  payload = unmarshal message.payload
  attrs =
    'consumption': payload.consumption
    'linkquality': payload.linkquality
    'power': payload.power
    'state': payload.state
    'temperature': payload.temperature
    'voltage': payload.voltage
  SetState ENTITY_ID,
    'new_state': payload.state
    'attribute_values': attrs
`

		plugActionOnOffSourceScript = `
entityAction = (entityId, actionName)->
    #print '---action on/off--'
    topic = entityId.replace(".", "/")+'/set'
    payload = marshal {"state": actionName}
    Mqtt.publish(topic, payload, 0, false)
`

		task1SourceScript = `
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    p = unmarshal msg.payload
    return p.new_state.state.name == 'DOUBLE_CLICK'

automationCondition = (entityId)->
    #print '---condition---'
    entity = GetEntity('zigbee2mqtt.` + zigbeePlugId + `')
    if !entity
        return false
    if !entity.state || entity.state.name == 'OFF'
        return true
    return false

automationAction = (entityId)->
    #print '---action---'
    CallAction('zigbee2mqtt.` + zigbeePlugId + `', 'ON', {})
`

		task2SourceScript = `
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    p = unmarshal msg.payload
    return p.new_state.state.name == 'DOUBLE_CLICK'

automationCondition = (entityId)->
    print '---condition---'
    entity = Condition.getEntity('zigbee2mqtt.` + zigbeePlugId + `')
    if !entity || !entity.state 
        return false
    if entity.state.name == 'ON'
        return true
    return false

automationAction = (entityId)->
    #print '---action---'
    CallAction('zigbee2mqtt.` + zigbeePlugId + `', 'OFF', {})
`
	)

	Convey("zigbee2mqtt", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus) {

			AddPlugin(adaptors, "zigbee2mqtt")
			AddPlugin(adaptors, "triggers")

			go mqttServer.Start()

			time.Sleep(time.Millisecond * 500)

			// add zigbee2mqtt
			zigbeeServer := &m.Zigbee2mqtt{
				Name:       "main",
				Devices:    nil,
				PermitJoin: true,
				BaseTopic:  "zigbee2mqtt",
			}
			var err error
			zigbeeServer.Id, err = adaptors.Zigbee2mqtt.Add(context.Background(), zigbeeServer)
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
			err = adaptors.Zigbee2mqttDevice.Add(context.Background(), butonDevice)
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
			err = adaptors.Zigbee2mqttDevice.Add(context.Background(), plugDevice)
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
			err = adaptors.Entity.Add(context.Background(), buttonEnt)
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
			err = adaptors.Entity.Add(context.Background(), plugEnt)
			So(err, ShouldBeNil)

			time.Sleep(time.Second)

			// automation
			// ------------------------------------------------
			trigger1 := &m.Trigger{
				Name:       "state_change",
				EntityId:   &buttonEnt.Id,
				Script:     task1Script,
				PluginName: "state_change",
			}
			err = AddTrigger(trigger1, adaptors, eventBus)
			So(err, ShouldBeNil)

			trigger2 := &m.Trigger{
				Name:       "",
				EntityId:   &buttonEnt.Id,
				Script:     task2Script,
				PluginName: "state_change",
			}
			err = AddTrigger(trigger2, adaptors, eventBus)
			So(err, ShouldBeNil)

			// conditions
			// -----------------------
			condition1 := &m.Condition{
				Name:   "check plug state",
				Script: task1Script,
			}
			condition1.Id, err = adaptors.Condition.Add(context.Background(), condition1)
			So(err, ShouldBeNil)

			condition2 := &m.Condition{
				Name:   "check plug state",
				Script: task2Script,
			}
			condition2.Id, err = adaptors.Condition.Add(context.Background(), condition2)
			So(err, ShouldBeNil)

			// actions
			// -----------------------
			action1 := &m.Action{
				Name:   "action on Plug",
				Script: task1Script,
			}
			action1.Id, err = adaptors.Action.Add(context.Background(), action1)
			So(err, ShouldBeNil)

			action2 := &m.Action{
				Name:   "action on Plug",
				Script: task2Script,
			}
			action2.Id, err = adaptors.Action.Add(context.Background(), action2)
			So(err, ShouldBeNil)

			// tasks
			// -----------------------
			newTask1 := &m.NewTask{
				Name:         "Toggle plug ON",
				Enabled:      true,
				Condition:    common.ConditionAnd,
				TriggerIds:   []int64{trigger1.Id},
				ConditionIds: []int64{condition1.Id},
				ActionIds:    []int64{action1.Id},
			}

			err = AddTask(newTask1, adaptors, eventBus)
			So(err, ShouldBeNil)

			//TASK2
			newTask2 := &m.NewTask{
				Name:         "Toggle plug OFF",
				Enabled:      true,
				Condition:    common.ConditionAnd,
				TriggerIds:   []int64{trigger2.Id},
				ConditionIds: []int64{condition2.Id},
				ActionIds:    []int64{action2.Id},
			}
			err = AddTask(newTask2, adaptors, eventBus)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			automation.Start()
			go zigbee2mqtt.Start(context.Background())

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

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

			time.Sleep(time.Second)

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
