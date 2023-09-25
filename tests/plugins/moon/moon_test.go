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

package moon

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	moonPlugin "github.com/e154/smart-home/plugins/moon"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMoon(t *testing.T) {

	Convey("moon", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "moon")

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			// add entity
			// ------------------------------------------------
			moonEnt := GetNewMoon("main")
			err := adaptors.Entity.Add(context.Background(), moonEnt)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+moonEnt.Id.String(), events.EventCreatedEntity{
				EntityId: moonEnt.Id,
			})

			time.Sleep(time.Second)

			ch := make(chan events.EventStateChanged, 2)
			_ = eventBus.Subscribe("system/entities/"+moonEnt.Id.String(), func(topic string, msg events.EventStateChanged) {
				ch <- msg
			})

			moon := moonPlugin.NewActor(moonEnt, supervisor.GetService())

			t.Run("entity", func(t *testing.T) {
				Convey("position", t, func(ctx C) {
					moonEnt, err = adaptors.Entity.GetById(context.Background(), moonEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(moonEnt.Settings[moonPlugin.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(moonEnt.Settings[moonPlugin.AttrLon].Float64(), ShouldEqual, 83.0335)
				})
			})

			t.Run("actor", func(t *testing.T) {
				Convey("attrs", t, func(ctx C) {

					loc, err := time.LoadLocation("Asia/Novosibirsk")
					ctx.So(err, ShouldBeNil)
					moon.UpdateMoonPosition(time.Date(2021, 5, 27, 23, 0, 0, 0, loc))

					time.Sleep(time.Second * 2)

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg events.EventStateChanged
					var ok bool
					select {
					case msg = <-ch:
						ok = true
					case <-ticker.C:
					}

					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, moonPlugin.StateBelowHorizon)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "below horizon")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrAzimuth].String(), ShouldContainSubstring, "131.932179124")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrElevation].String(), ShouldContainSubstring, "-1.206958128")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrPhase].String(), ShouldEqual, "full_moon")
				})
			})
		})
	})
}
