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
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"
)

func TestTriggerSystem(t *testing.T) {

	const (
		task3SourceScript = `
automationTriggerSystem = (msg)->
    #print '---trigger---'
    Done msg.payload.event
    return false
`
	)

	Convey("trigger system", t, func(ctx C) {
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

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			ctx.So(err, ShouldBeNil)

			go mqttServer.Start()

			// task3 scripts
			task3Script := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "task3",
				Source:      task3SourceScript,
				Description: "",
			}

			task3ScriptEngine, err := scriptService.NewEngine(task3Script)
			So(err, ShouldBeNil)
			err = task3ScriptEngine.Compile()
			So(err, ShouldBeNil)

			task3Script.Id, err = adaptors.Script.Add(task3Script)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------

			//TASK3
			task3 := &m.Task{
				Name:      "Toggle plug OFF",
				Enabled:   true,
				Condition: common.ConditionAnd,
			}
			task3.AddTrigger(&m.Trigger{
				Name:       "",
				Script:     task3Script,
				PluginName: "system",
			})
			err = adaptors.Task.Add(task3)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			var counter atomic.Int32
			var lastEvent atomic.String
			scriptService.PushFunctions("Done", func(systemEvent string) {
				lastEvent.Store(systemEvent)
				counter.Inc()
			})

			pluginManager.Start()
			automation.Reload()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()
			go zigbee2mqtt.Start()

			defer func() {
				mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			eventBus.Publish(triggers.TopicSystemStart, "started")

			//
			// ------------------------------------------------

			automation.Reload()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()
			go zigbee2mqtt.Start()

			time.Sleep(time.Second)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)
			So(lastEvent.Load(), ShouldEqual, "START")

		})
	})
}
