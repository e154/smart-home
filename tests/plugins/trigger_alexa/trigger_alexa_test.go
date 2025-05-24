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

package trigger_alexa

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/internal/plugins/alexa"
	"github.com/e154/smart-home/internal/system/automation"
	"github.com/e154/smart-home/pkg/adaptors"
	commonPkg "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggerAlexa(t *testing.T) {

	const (
		task3SourceScript = `
automationTriggerAlexa = (msg)->
    #print '---trigger---'
    p = msg.payload
    Done p
    #msg.trigger_name
    #msg.task_name
    #msg.entity_id
    return false
`
	)

	Convey("trigger alexa", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "triggers")
			AddPlugin(adaptors, "alexa")

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor", "Automation")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "triggers", "alexa")
			automation.Start()
			supervisor.Start(context.Background())
			defer automation.Shutdown()
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			var ch = make(chan alexa.EventAlexaAction)
			defer close(ch)
			scriptService.PushFunctions("Done", func(msg alexa.EventAlexaAction) {
				ch <- msg
			})

			// add scripts
			// ------------------------------------------------

			task3Script, err := AddScript("task3", task3SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------
			trigger := &models.NewTrigger{
				Enabled:    true,
				Name:       "alexa",
				ScriptId:   commonPkg.Int64(task3Script.Id),
				PluginName: "alexa",
				Payload: models.Attributes{
					alexa.TriggerOptionSkillId: {
						Name:  alexa.TriggerOptionSkillId,
						Type:  commonPkg.AttributeInt,
						Value: 1,
					},
				},
			}
			triggerId, err := AddTrigger(trigger, adaptors, eventBus)
			So(err, ShouldBeNil)

			//TASK3
			newTask := &models.NewTask{
				Name:       "Toggle plug ON",
				Enabled:    true,
				TriggerIds: []int64{triggerId},
				Condition:  commonPkg.ConditionAnd,
			}
			_, err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Millisecond * 500)

			// ------------------------------------------------

			eventBus.Publish(alexa.TopicPluginAlexa, alexa.EventAlexaAction{
				SkillId:    1,
				IntentName: "FlatLights",
				Payload:    "kitchen_on",
			})

			time.Sleep(time.Millisecond * 500)

			// wait message
			msg, ok := WaitT[alexa.EventAlexaAction](time.Second*2, ch)
			ctx.So(ok, ShouldBeTrue)

			ctx.So(msg.Payload, ShouldEqual, "kitchen_on")
			ctx.So(msg.SkillId, ShouldEqual, 1)
			ctx.So(msg.IntentName, ShouldEqual, "FlatLights")

			time.Sleep(time.Second)

		})
	})
}
