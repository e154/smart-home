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

package moon

import (
	"context"
	"testing"
	"time"

	moon2 "github.com/e154/smart-home/internal/plugins/moon"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMoon(t *testing.T) {

	Convey("moon", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "moon")

			// add entity
			// ------------------------------------------------
			moonEnt := GetNewMoon("main")
			err := adaptors.Entity.Add(context.Background(), moonEnt)
			ctx.So(err, ShouldBeNil)

			// common
			// ------------------------------------------------
			ch := make(chan events.EventStateChanged)
			defer close(ch)
			fn := func(topic string, msg interface{}) {
				switch v := msg.(type) {
				case events.EventStateChanged:
					ch <- v
				}
			}

			eventBus.Subscribe("system/entities/"+moonEnt.Id.String(), fn)
			defer eventBus.Unsubscribe("system/entities/"+moonEnt.Id.String(), fn)
			// ------------------------------------------------

			moon := moon2.NewActor(moonEnt, supervisor.GetService())

			t.Run("entity", func(t *testing.T) {
				Convey("position", t, func(ctx C) {
					moonEnt, err = adaptors.Entity.GetById(context.Background(), moonEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(moonEnt.Settings[moon2.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(moonEnt.Settings[moon2.AttrLon].Float64(), ShouldEqual, 83.0335)
				})
			})

			t.Run("actor", func(t *testing.T) {
				Convey("attrs", t, func(ctx C) {

					loc, err := time.LoadLocation("Asia/Novosibirsk")
					ctx.So(err, ShouldBeNil)
					moon.UpdateMoonPosition(time.Date(2021, 5, 27, 23, 0, 0, 0, loc))

					// wait message
					msg, ok := WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, moon2.StateBelowHorizon)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "below horizon")
					ctx.So(msg.NewState.Attributes[moon2.AttrAzimuth].String(), ShouldContainSubstring, "131.932179124")
					ctx.So(msg.NewState.Attributes[moon2.AttrElevation].String(), ShouldContainSubstring, "-1.206958128")
					ctx.So(msg.NewState.Attributes[moon2.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[moon2.AttrPhase].String(), ShouldEqual, "full_moon")
				})
			})
		})
	})
}
