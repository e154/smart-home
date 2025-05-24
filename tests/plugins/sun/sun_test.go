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

	sun2 "github.com/e154/smart-home/internal/plugins/sun"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSun(t *testing.T) {

	Convey("sun", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
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
			sun := sun2.NewActor(sunEnt, supervisor.GetService())

			t.Run("entity", func(t *testing.T) {
				Convey("phase", t, func(ctx C) {
					sunEnt, err = adaptors.Entity.GetById(context.Background(), sunEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(sunEnt.Settings[sun2.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(sunEnt.Settings[sun2.AttrLon].Float64(), ShouldEqual, 83.0335)
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
					ctx.So(msg.NewState.State.Name, ShouldEqual, sun2.AttrDusk)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "dusk (evening nautical twilight starts)")
					ctx.So(msg.NewState.Attributes[sun2.AttrAzimuth].String(), ShouldContainSubstring, "326.127522661")
					ctx.So(msg.NewState.Attributes[sun2.AttrElevation].String(), ShouldContainSubstring, "-7.663125133")
					ctx.So(msg.NewState.Attributes[sun2.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[sun2.AttrPhase].String(), ShouldEqual, "dusk")
				})
			})
		})
	})
}
