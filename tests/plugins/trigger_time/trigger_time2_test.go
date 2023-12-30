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

package trigger_time

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestTriggerTime2(t *testing.T) {

	const (
		task3SourceScript = `
entityAction = (entityId, actionName)->
    #print '---action---'
    Done()
`
	)

	Convey("trigger time", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler *scheduler.Scheduler,
		) {

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func() {
				counter.Inc()
			})
			defer scriptService.PopStruct("Done")

			// add scripts
			// ------------------------------------------------

			task3Script, err := AddScript("task3", task3SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------

			sensorEnt := GetNewSensor("device1")
			sensorEnt.Actions = []*m.EntityAction{
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      task3Script,
				},
			}
			err = adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)

			sensorEnt, err = adaptors.Entity.GetById(context.Background(), sensorEnt.Id)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+sensorEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: sensorEnt.Id,
			})

			// automation
			// ------------------------------------------------
			trigger := &m.NewTrigger{
				Name:       "trigger1",
				PluginName: "time",
				Payload: m.Attributes{
					triggers.CronOptionTrigger: {
						Name:  triggers.CronOptionTrigger,
						Type:  common.AttributeString,
						Value: "* * * * * *", //every seconds
					},
				},
			}
			triggerId, err := AddTrigger(trigger, adaptors, eventBus)
			So(err, ShouldBeNil)

			action := &m.Action{
				Name:             "action1",
				EntityId:         common.NewEntityId(string(sensorEnt.Id)),
				EntityActionName: common.String(sensorEnt.Actions[0].Name),
			}
			action.Id, err = adaptors.Action.Add(context.Background(), action)
			So(err, ShouldBeNil)

			//TASK3
			newTask := &m.NewTask{
				Name:       "Toggle plug OFF",
				Enabled:    true,
				Condition:  common.ConditionAnd,
				TriggerIds: []int64{triggerId},
				ActionIds:  []int64{action.Id},
			}
			err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 2)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)

		})
	})
}
