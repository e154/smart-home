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

package scene

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/initial/local_migrations"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestScene(t *testing.T) {

	const (
		sceneSourceScript = `
sceneEvent = (args)->
  #print '---action on/off--'
  Done args
`
	)

	Convey("scene", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			automation automation.Automation,
			eventBus bus.Bus,
			scheduler *scheduler.Scheduler) {

			// register plugins
			err := local_migrations.NewMigrationPlugins(adaptors).Up(context.TODO(), nil)
			So(err, ShouldBeNil)

			scheduler.Start(context.TODO())
			automation.Start()
			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			var counter atomic.Int32
			scriptService.PushFunctions("Done", func(args string) {
				counter.Inc()
			})

			// add scripts
			// ------------------------------------------------

			sceneScript, err := AddScript("romantic script", sceneSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			romanticEnt := GetNewScene("scene.romantic", []*m.Script{sceneScript})
			err = adaptors.Entity.Add(context.Background(), romanticEnt)
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+romanticEnt.Id.String(), events.EventCreatedEntity{
				EntityId: romanticEnt.Id,
			})

			time.Sleep(time.Millisecond * 500)

			t.Run("call scene", func(t *testing.T) {
				Convey("case", t, func(ctx C) {

					supervisor.CallScene(romanticEnt.Id, nil)
					time.Sleep(time.Millisecond * 500)

					So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 1)
				})
			})
		})
	})
}
