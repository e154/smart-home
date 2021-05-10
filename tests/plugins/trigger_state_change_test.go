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
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

func TestTriggerStateChange(t *testing.T) {

	const (
		zigbeeButtonId = "0x00158d00031c8ef3"

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

		task1SourceScript = `
automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    Done msg.new_state.state.name
    return false
`
	)

	Convey("trigger state change", t, func(ctx C) {
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
			AddPlugin(adaptors, "triggers")
			AddPlugin(adaptors, "zigbee2mqtt")

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
			}
			err = adaptors.Zigbee2mqttDevice.Add(butonDevice)
			So(err, ShouldBeNil)

			// add scripts
			// ------------------------------------------------
			buttonScript := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "button",
				Source:      buttonSourceScript,
				Description: "button sensor",
			}

			engine1, err := scriptService.NewEngine(buttonScript)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)

			buttonScript.Id, err = adaptors.Script.Add(buttonScript)
			So(err, ShouldBeNil)

			// task1 scripts
			task1Script := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "trigger_script_entityId_",
				Source:      task1SourceScript,
				Description: "",
			}

			task1ScriptEngine, err := scriptService.NewEngine(task1Script)
			So(err, ShouldBeNil)
			err = task1ScriptEngine.Compile()
			So(err, ShouldBeNil)

			task1Script.Id, err = adaptors.Script.Add(task1Script)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			buttonEnt := GetNewButton(fmt.Sprintf("zigbee2mqtt.%s", zigbeeButtonId), []m.Script{*buttonScript})
			err = adaptors.Entity.Add(buttonEnt)
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
			err = adaptors.Task.Add(task1)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			var counter atomic.Int32
			var lastStat atomic.String
			var wg sync.WaitGroup
			wg.Add(2)
			scriptService.PushFunctions("Done", func(state string) {
				lastStat.Store(state)
				counter.Inc()
				wg.Done()
			})

			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)
			go zigbee2mqtt.Start()

			defer func() {
				mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			mqttCli := mqttServer.NewClient("cli2")
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"action":"double","linkquality":134,"voltage":3042}`))
			time.Sleep(time.Millisecond * 100)
			err = mqttCli.Publish("zigbee2mqtt/"+zigbeeButtonId, []byte(`{"battery":100,"click":"double","linkquality":134,"voltage":3042}`))
			So(err, ShouldBeNil)

			wg.Wait()

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)
			So(lastStat.Load(), ShouldEqual, "DOUBLE_CLICK")
		})
	})
}
