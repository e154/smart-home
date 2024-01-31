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

package trigger_empty

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggerEmpty(t *testing.T) {

	Convey("trigger empty", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler *scheduler.Scheduler,
		) {

			// register plugins
			AddPlugin(adaptors, "triggers")

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor", "Automation", "Scheduler")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "triggers")
			scheduler.Start(context.Background())
			automation.Start()
			supervisor.Start(context.Background())
			defer scheduler.Shutdown(context.Background())
			defer automation.Shutdown()
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			// automation
			// ------------------------------------------------
			trigger := &m.NewTrigger{
				Enabled:    true,
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

			//TASK3
			newTask := &m.NewTask{
				Name:       "Toggle plug OFF",
				Enabled:    true,
				Condition:  common.ConditionAnd,
				TriggerIds: []int64{triggerId},
			}
			_, err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)

			time.Sleep(time.Millisecond * 500)

		})
	})
}
