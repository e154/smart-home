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

	timeTrigger "github.com/e154/smart-home/internal/plugins/time"
	"github.com/e154/smart-home/internal/system/automation"
	"github.com/e154/smart-home/pkg/adaptors"
	commonPkg "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scheduler"
	"github.com/e154/smart-home/pkg/scripts"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/bus"
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
			supervisor plugins.Supervisor,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler scheduler.Scheduler,
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
			sensorEnt.Actions = []*models.EntityAction{
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
			trigger := &models.NewTrigger{
				Enabled:    true,
				Name:       "trigger1",
				PluginName: "time",
				Payload: models.Attributes{
					timeTrigger.AttrCron: {
						Name:  timeTrigger.AttrCron,
						Type:  commonPkg.AttributeString,
						Value: "* * * * * *", //every seconds
					},
				},
			}
			triggerId, err := AddTrigger(trigger, adaptors, eventBus)
			So(err, ShouldBeNil)

			action := &models.Action{
				Name:             "action1",
				EntityId:         commonPkg.NewEntityId(string(sensorEnt.Id)),
				EntityActionName: commonPkg.String(sensorEnt.Actions[0].Name),
			}
			action.Id, err = adaptors.Action.Add(context.Background(), action)
			So(err, ShouldBeNil)

			//TASK3
			newTask := &models.NewTask{
				Name:       "Toggle plug OFF",
				Enabled:    true,
				Condition:  commonPkg.ConditionAnd,
				TriggerIds: []int64{triggerId},
				ActionIds:  []int64{action.Id},
			}
			_, err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 2)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)

		})
	})
}
