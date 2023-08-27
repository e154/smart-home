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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTriggerEmpty(t *testing.T) {

	Convey("trigger time", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler *scheduler.Scheduler,
		) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			ctx.So(err, ShouldBeNil)

			eventBus.Purge()
			scriptService.Restart()
			scheduler.Start(context.Background())
			automation.Restart()
			supervisor.Restart(context.Background())
			defer func() {
				scheduler.Shutdown(context.Background())
			}()

			time.Sleep(time.Millisecond * 500)

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
			err = adaptors.Task.Import(task3)
			So(err, ShouldBeNil)

			eventBus.Publish(bus.TopicAutomation, events.EventAddedTask{
				Id: task3.Id,
			})

			time.Sleep(time.Millisecond * 500)

		})
	})
}
