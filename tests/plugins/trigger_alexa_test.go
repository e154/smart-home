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
	"github.com/e154/smart-home/plugins/alexa"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggerAlexa(t *testing.T) {

	const (
		task3SourceScript = `
automationTriggerAlexa = (msg)->
    #print '---trigger---'
    p = unmarshal msg.payload
    Done p
    #msg.trigger_name
    #msg.task_name
    #msg.entity_id
    return false
`
	)

	Convey("trigger alexa", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			ctx.So(err, ShouldBeNil)
			err = AddPlugin(adaptors, "alexa")
			ctx.So(err, ShouldBeNil)

			go mqttServer.Start()

			// add scripts
			// ------------------------------------------------

			task3Script, err := AddScript("task3", task3SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------

			//TASK3
			task3 := &m.Task{
				Name:      "Toggle plug ON",
				Enabled:   true,
				Condition: common.ConditionAnd,
			}
			task3.AddTrigger(&m.Trigger{
				Name:       "alexa",
				Script:     task3Script,
				PluginName: "alexa",
				Payload: m.Attributes{
					alexa.TriggerOptionSkillId: {
						Name:  alexa.TriggerOptionSkillId,
						Type:  common.AttributeInt,
						Value: 1,
					},
				},
			})
			err = adaptors.Task.Add(task3)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			var ch = make(chan alexa.EventAlexaAction)
			scriptService.PushFunctions("Done", func(msg alexa.EventAlexaAction) {
				ch <- msg
			})

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

			eventBus.Publish(alexa.TopicPluginAlexa, alexa.EventAlexaAction{
				SkillId:    1,
				IntentName: "FlatLights",
				Payload:    "kitchen_on",
			})

			ticker := time.NewTimer(time.Second * 2)
			defer ticker.Stop()

			var msg alexa.EventAlexaAction
			select {
			case v := <-ch:
				msg = v
				break
			case <-ticker.C:
				break
			}

			ctx.So(msg.Payload, ShouldEqual, "kitchen_on")
			ctx.So(msg.SkillId, ShouldEqual, 1)
			ctx.So(msg.IntentName, ShouldEqual, "FlatLights")

		})
	})
}
