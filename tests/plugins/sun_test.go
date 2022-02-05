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
	"github.com/e154/smart-home/system/event_bus/events"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	sunPlugin "github.com/e154/smart-home/plugins/sun"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSun(t *testing.T) {

	Convey("sun", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "sun")
			ctx.So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			sunEnt := GetNewSun("main")
			err = adaptors.Entity.Add(sunEnt)
			ctx.So(err, ShouldBeNil)

			ch := make(chan events.EventStateChanged)
			eventBus.Subscribe(event_bus.TopicEntities, func(topic string, msg events.EventStateChanged) {
				ch <- msg
			})

			sun := sunPlugin.NewActor(sunEnt, entityManager, eventBus)

			t.Run("entity", func(t *testing.T) {
				Convey("phase", t, func(ctx C) {
					sunEnt, err = adaptors.Entity.GetById(sunEnt.Id)
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

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg events.EventStateChanged
					var ok bool
					select {
					case msg = <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, sunPlugin.AttrDusk)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "dusk (evening nautical twilight starts)")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrAzimuth].Float64(), ShouldEqual, 326.12752266195463)
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrElevation].Float64(), ShouldEqual, -7.663125133580264)
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[sunPlugin.AttrPhase].String(), ShouldEqual, "dusk")
				})
			})
		})
	})
}
