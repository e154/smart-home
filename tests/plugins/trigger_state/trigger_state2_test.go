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
	"github.com/e154/smart-home/system/migrations"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	auto "github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/e154/smart-home/tests/plugins/container"
)

func TestTriggerState2(t *testing.T) {

	t.Run("init", func(t *testing.T) {
		Convey("", t, func(ctx C) {
			err := BuildContainer().Invoke(func(adaptors *adaptors.Adaptors,
				scriptService scripts.ScriptService,
				supervisor supervisor.Supervisor,
				mqttServer mqtt.MqttServ,
				automation auto.Automation,
				eventBus bus.Bus,
				migrations *migrations.Migrations) {

				migrations.Purge()

				// register plugins
				AddPlugin(adaptors, "triggers")

				waitCh := WaitService(eventBus, time.Second*5, "Supervisor")
				pluginsCh := WaitPlugins(eventBus, time.Second*5, "triggers")

				supervisor.Start(context.Background())
				defer supervisor.Shutdown(context.Background())
				So(<-waitCh, ShouldBeTrue)
				So(<-pluginsCh, ShouldBeTrue)

				script := &m.Script{
					Id:       1,
					Compiled: `function automationTriggerStateChanged(msg) {return true}`,
				}

				model := &m.Trigger{
					Id:         1,
					Name:       "tr1",
					PluginName: "state_change",
					Script:     script,
					Enabled:    true,
				}

				plugin, err := supervisor.GetPlugin("triggers")
				So(err, ShouldBeNil)

				rawPlugin, ok := plugin.(triggers.IGetTrigger)
				So(ok, ShouldBeTrue)

				trigger, err := auto.NewTrigger(eventBus, scriptService, model, rawPlugin)
				So(err, ShouldBeNil)
				So(trigger, ShouldNotBeNil)

				trigger.Start()

				msg, ok := WaitMessage[events.EventTriggerLoaded](eventBus, time.Second*2, "system/automation/triggers/1", true)
				So(ok, ShouldBeTrue)
				So(msg, ShouldNotBeNil)

				So(msg.Id, ShouldEqual, 1)

				// test 2
				// ------------------------------------------------
				script3 := &m.Script{
					Id:   3,
					Name: "script3",
					Compiled: `
function automationTriggerStateChanged(msg) {
	//console.log(marshal(msg))
	return msg.payload.new_state.attributes.foo == "omega";
}`,
				}

				entityFoo := GetNewSensor("foo")
				entityBar := GetNewSensor("bar")

				model2 := &m.Trigger{
					Name:       "tr1",
					PluginName: "state_change",
					Id:         2,
					Script:     script3,
					Enabled:    true,
					Entities:   []*m.Entity{entityFoo, entityBar},
				}

				trigger2, err := auto.NewTrigger(eventBus, scriptService, model2, rawPlugin)
				So(err, ShouldBeNil)
				So(trigger2, ShouldNotBeNil)

				trigger2.Start()

				eventBus.Publish(fmt.Sprintf("system/entities/%s", entityFoo.Id), events.EventStateChanged{
					EntityId: entityFoo.Id,
					NewState: events.EventEntityState{
						Attributes: m.Attributes{
							"foo": &m.Attribute{
								Name:  "foo",
								Type:  common.AttributeString,
								Value: "bar",
							},
						},
					},
				})

				_, ok = WaitMessage[events.EventTriggerCompleted](eventBus, time.Second, "system/automation/triggers/2", false)
				So(ok, ShouldBeFalse)

				// update script
				script3.Compiled = `
				function automationTriggerStateChanged(msg) {
					//console.log(marshal(msg))
					return msg.payload.new_state.attributes.foo == "bar";
				}`
				eventBus.Publish("system/models/scripts/3", events.EventUpdatedScriptModel{
					ScriptId: 3,
					Script:   script3,
				})

				time.Sleep(time.Millisecond * 500)

				eventBus.Publish(fmt.Sprintf("system/entities/%s", entityFoo.Id), events.EventStateChanged{
					EntityId: entityFoo.Id,
					NewState: events.EventEntityState{
						Attributes: m.Attributes{
							"foo": &m.Attribute{
								Name:  "foo",
								Type:  common.AttributeString,
								Value: "bar",
							},
						},
					},
				})

				msg3, ok := WaitMessage[events.EventTriggerCompleted](eventBus, time.Second, "system/automation/triggers/2", false)
				So(ok, ShouldBeTrue)
				So(msg3, ShouldNotBeNil)

				//debug.Println(msg3)

				So(msg3.Args, ShouldNotBeNil)

				ev, ok := msg3.Args.Payload.(triggers.TriggerStateChangedMessage)
				So(ok, ShouldBeTrue)
				//debug.Println(ev)
				So(ev.NewState, ShouldNotBeNil)
				So(ev.NewState.Attributes["foo"], ShouldEqual, "bar")

				time.Sleep(time.Millisecond * 500)

			})
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	})
}
