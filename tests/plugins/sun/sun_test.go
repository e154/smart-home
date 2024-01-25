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

package sun

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	sunPlugin "github.com/e154/smart-home/plugins/sun"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSun(t *testing.T) {

	Convey("sun", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			AddPlugin(adaptors, "sun")

			// common
			// ------------------------------------------------
			sunEnt := GetNewSun("main")
			err := adaptors.Entity.Add(context.Background(), sunEnt)
			ctx.So(err, ShouldBeNil)

			ch := make(chan events.EventStateChanged, 2)
			defer close(ch)
			fn := func(topic string, msg interface{}) {
				switch v := msg.(type) {
				case events.EventStateChanged:
					ch <- v
				}
			}
			eventBus.Subscribe("system/entities/"+sunEnt.Id.String(), fn)
			defer eventBus.Unsubscribe("system/entities/"+sunEnt.Id.String(), fn)

			// add entity
			// ------------------------------------------------
			sun := sunPlugin.NewActor(sunEnt, supervisor.GetService())

			t.Run("entity", func(t *testing.T) {
				Convey("phase", t, func(ctx C) {
					sunEnt, err = adaptors.Entity.GetById(context.Background(), sunEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(sunEnt.Settings[sunPlugin.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(sunEnt.Settings[sunPlugin.AttrLon].Float64(), ShouldEqual, 83.0335)
				})
			})

			t.Run("actor", func(t *testing.T) {
				Convey("attrs", t, func(ctx C) {

					loc, err := time.LoadLocation("Asia/Novosibirsk")
					ctx.So(err, ShouldBeNil)
					sun.UpdateSunPosition(time.Date(2021, 5, 27, 23, 0, 0, 0, loc))

					// wait message
					msg, ok := WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, sunPlugin.AttrDusk)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "dusk (evening nautical twilight starts)")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrAzimuth].String(), ShouldContainSubstring, "326.127522661")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrElevation].String(), ShouldContainSubstring, "-7.663125133")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrPhase].String(), ShouldEqual, "dusk")
				})
			})
		})
	})
}
