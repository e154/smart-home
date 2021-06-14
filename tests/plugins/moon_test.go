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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	moonPlugin "github.com/e154/smart-home/plugins/moon"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestMoon(t *testing.T) {

	Convey("moon", t, func(ctx C) {
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
			err = AddPlugin(adaptors, "moon")
			ctx.So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			moonEnt := GetNewMoon("main")
			err = adaptors.Entity.Add(moonEnt)
			ctx.So(err, ShouldBeNil)

			ch := make(chan event_bus.EventStateChanged)
			eventBus.Subscribe(event_bus.TopicEntities, func(topic string, msg event_bus.EventStateChanged) {
				ch <- msg
			})

			moon := moonPlugin.NewActor(moonEnt, entityManager, eventBus)

			t.Run("entity", func(t *testing.T) {
				Convey("position", t, func(ctx C) {
					moonEnt, err = adaptors.Entity.GetById(moonEnt.Id)
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

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg event_bus.EventStateChanged
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
					ctx.So(msg.NewState.State.Name, ShouldEqual, moonPlugin.StateBelowHorizon)
					ctx.So(msg.NewState.State.Description, ShouldEqual, "below horizon")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrAzimuth].Float64(), ShouldEqual, 131.93217912435995)
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrElevation].Float64(), ShouldEqual, -1.2069581282932162)
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrHorizonState].String(), ShouldEqual, "belowHorizon")
					ctx.So(msg.NewState.Attributes[moonPlugin.AttrPhase].String(), ShouldEqual, "full_moon")
				})
			})
		})
	})
}
