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
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
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
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler *scheduler.Scheduler,
		) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			ctx.So(err, ShouldBeNil)
			err = AddPlugin(adaptors, "sensor")
			ctx.So(err, ShouldBeNil)

			eventBus.Purge()
			scriptService.Restart()
			scheduler.Start(context.TODO())
			defer func() {
				scheduler.Shutdown(context.Background())
			}()

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
			err = adaptors.Entity.Add(sensorEnt)
			ctx.So(err, ShouldBeNil)

			sensorEnt, err = adaptors.Entity.GetById(sensorEnt.Id)
			ctx.So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------

			//TASK3
			task3 := &m.Task{
				Name:      "Toggle plug OFF",
				Enabled:   true,
				Condition: common.ConditionAnd,
			}
			task3.AddTrigger(&m.Trigger{
				Name:       "trigger1",
				PluginName: "time",
				Payload: m.Attributes{
					triggers.CronOptionTrigger: {
						Name:  triggers.CronOptionTrigger,
						Type:  common.AttributeString,
						Value: "* * * * * *", //every seconds
					},
				},
			})
			task3.AddAction(&m.Action{
				Name:             "action1",
				EntityId:         common.NewEntityId(string(sensorEnt.Id)),
				EntityActionName: common.String(sensorEnt.Actions[0].Name),
			})
			err = adaptors.Task.Add(task3)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func() {
				counter.Inc()
			})

			supervisor.Restart(context.Background())
			automation.Restart()

			time.Sleep(time.Second)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)

		})
	})
}
