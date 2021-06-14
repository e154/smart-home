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
	weatherPlugin "github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestWeather(t *testing.T) {

	Convey("weather", t, func(ctx C) {
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
			err = AddPlugin(adaptors, "weather")
			ctx.So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			weatherEnt := GetNewWeather("main")
			err = adaptors.Entity.Add(weatherEnt)
			ctx.So(err, ShouldBeNil)

			ch := make(chan weatherPlugin.EventStateChanged)
			eventBus.Subscribe(weatherPlugin.TopicPluginWeather, func(topic string, msg weatherPlugin.EventStateChanged) {
				ch <- msg
			})

			weather := weatherPlugin.NewActor(weatherEnt, entityManager, eventBus)

			t.Run("entity", func(t *testing.T) {
				Convey("position", t, func(ctx C) {
					weatherEnt, err = adaptors.Entity.GetById(weatherEnt.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(weatherEnt.Settings[weatherPlugin.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(weatherEnt.Settings[weatherPlugin.AttrLon].Float64(), ShouldEqual, 83.0335)
				})
			})

			t.Run("actor", func(t *testing.T) {
				Convey("attrs", t, func(ctx C) {

					weather.UpdatePosition(weatherEnt.Settings)

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg weatherPlugin.EventStateChanged
					var ok bool
					select {
					case msg = <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.Settings[weatherPlugin.AttrLat].Float64(), ShouldEqual, 54.9022)
					ctx.So(msg.Settings[weatherPlugin.AttrLon].Float64(), ShouldEqual, 83.0335)
				})
			})
		})
	})
}
