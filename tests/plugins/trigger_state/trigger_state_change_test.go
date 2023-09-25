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

package trigger_state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"go.uber.org/atomic"

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

func TestTriggerStateChange(t *testing.T) {

	const (
		zigbeeButtonId = "0x00158d00031c8ef3"

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
    'storage_save': true
`

		task1SourceScript = `
automationTriggerStateChanged = (msg)->
  #print '---trigger---'
  p = unmarshal msg.payload
  Done p.new_state.state.name
  return false
`
	)

	Convey("trigger state change", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "triggers")
			AddPlugin(adaptors, "zigbee2mqtt")

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
			buttonDevice := &m.Zigbee2mqttDevice{
				Id:            zigbeeButtonId,
				Zigbee2mqttId: zigbeeServer.Id,
				Name:          zigbeeButtonId,
				Model:         "WXKG01LM",
				Description:   "MiJia wireless switch",
				Manufacturer:  "Xiaomi",
				Status:        "active",
				Payload:       []byte("{}"),
			}
			err = adaptors.Zigbee2mqttDevice.Add(context.Background(), buttonDevice)
			So(err, ShouldBeNil)

			automation.Start()
			mqttServer.Start()
			zigbee2mqtt.Start(context.Background())
			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			var counter atomic.Int32
			var lastStat atomic.String
			ch := make(chan struct{})
			scriptService.PushFunctions("Done", func(state string) {
				lastStat.Store(state)
				counter.Inc()
				if counter.Load() > 1 {
					close(ch)
				}
			})

			time.Sleep(time.Millisecond * 500)

			// add scripts
			// ------------------------------------------------

			buttonScript, err := AddScript("button", buttonSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			task1Script, err := AddScript("task", task1SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			buttonEnt := GetNewButton(fmt.Sprintf("zigbee2mqtt.%s", zigbeeButtonId), []*m.Script{buttonScript})
			err = adaptors.Entity.Add(context.Background(), buttonEnt)
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+buttonEnt.Id.String(), events.EventCreatedEntity{
				EntityId: buttonEnt.Id,
			})

			time.Sleep(time.Second)

			// automation
			// ------------------------------------------------
			trigger := &m.Trigger{
				Name:       "state_change",
				EntityId:   &buttonEnt.Id,
				Script:     task1Script,
				PluginName: "state_change",
			}
			err = AddTrigger(trigger, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Millisecond * 500)

			//TASK1
			newTask := &m.NewTask{
				Name:       "Toggle plug ON",
				Enabled:    true,
				TriggerIds: []int64{trigger.Id},
				Condition:  common.ConditionAnd,
			}
			err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Millisecond * 500)

			// ------------------------------------------------

			mqttCli := mqttServer.NewClient("cli2")
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"double","linkquality":134,"voltage":3042}`))
			So(err, ShouldBeNil)
			time.Sleep(time.Millisecond * 500)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"double","linkquality":134,"voltage":3042}`))
			So(err, ShouldBeNil)
			time.Sleep(time.Millisecond * 500)

			timer := time.NewTimer(time.Second * 4)
			defer timer.Stop()

			select {
			case <-timer.C:
			case <-ch:

			}

			time.Sleep(time.Second)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)
			So(lastStat.Load(), ShouldEqual, "DOUBLE_CLICK")
		})
	})
}
