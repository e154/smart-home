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

package script

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/script"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestScript(t *testing.T) {

	const (
		plugActionOnOffSourceScript = `
entityAction = (entityId, actionName)->
  #print '---action on/off--'
  Done entityId, actionName
`
	)

	Convey("script", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			err := AddPlugin(adaptors, "script")
			So(err, ShouldBeNil)

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			// add scripts
			// ------------------------------------------------

			plugActionOnOffScript, err := AddScript("plug script", plugActionOnOffSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func(entityId, action string) {
				counter.Inc()
			})

			time.Sleep(time.Second)

			// add entity
			// ------------------------------------------------
			plugEnt := GetNewScript("script.1", []*m.Script{})
			plugEnt.PluginName = script.EntityScript
			plugEnt.Actions = []*m.EntityAction{
				{
					Name:        "ON",
					Description: "включить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "OFF",
					Description: "выключить",
					Script:      plugActionOnOffScript,
				},
			}
			err = adaptors.Entity.Add(plugEnt)
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+plugEnt.Id.String(), events.EventCreatedEntity{
				EntityId: plugEnt.Id,
			})

			time.Sleep(time.Second)

			// automation
			// ------------------------------------------------

			supervisor.CallAction(plugEnt.Id, "ON", nil)
			supervisor.CallAction(plugEnt.Id, "OFF", nil)
			supervisor.CallAction(plugEnt.Id, "NULL", nil)
			time.Sleep(time.Millisecond * 500)

			So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 2)

		})
	})
}
