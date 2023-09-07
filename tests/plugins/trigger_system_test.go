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
	"fmt"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/mqtt"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"
)

func TestTriggerSystem(t *testing.T) {

	const (
		task3SourceScript = `
automationTriggerSystem = (msg)->
    #print '---trigger---'
    p = unmarshal msg.payload
    Done p.event
    return false
`
	)

	Convey("trigger system", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			ctx.So(err, ShouldBeNil)

			eventBus.Purge()
			scriptService.Restart()

			go mqttServer.Start()
			supervisor.Restart(context.Background())
			automation.Restart()
			go zigbee2mqtt.Start()
			defer func() {
				zigbee2mqtt.Shutdown()
			}()

			time.Sleep(time.Second)

			var counter atomic.Int32
			var lastEvent atomic.String
			scriptService.PushFunctions("Done", func(systemEvent string) {
				lastEvent.Store(systemEvent)
				counter.Inc()
			})

			// add scripts
			// ------------------------------------------------

			task3Script, err := AddScript("task3", task3SourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// automation
			// ------------------------------------------------
			trigger := &m.Trigger{
				Name:       "tr1",
				Script:     task3Script,
				PluginName: "system",
			}
			err = AddTrigger(trigger, adaptors, eventBus)
			So(err, ShouldBeNil)

			//TASK3
			newTask := &m.NewTask{
				Name:      "Toggle plug OFF",
				Enabled:   true,
				Condition: common.ConditionAnd,
				TriggerIds: []int64{trigger.Id},
			}
			err = AddTask(newTask, adaptors, eventBus)
			So(err, ShouldBeNil)



			// ------------------------------------------------

			time.Sleep(time.Second)

			eventBus.Publish(triggers.TopicSystemStart, "started")

			time.Sleep(time.Second)

			fmt.Println(counter.Load())

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)
			So(lastEvent.Load(), ShouldEqual, "START")

		})
	})
}
