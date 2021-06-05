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
	"github.com/e154/smart-home/plugins/email"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestEmail(t *testing.T) {

	Convey("email", t, func(ctx C) {
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
			err = AddPlugin(adaptors, "notify")
			err = AddPlugin(adaptors, "email")
			ctx.So(err, ShouldBeNil)

			pluginManager.Start()
			entityManager.LoadEntities(pluginManager)

			defer func() {
				entityManager.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: email.Name,
						Attributes: map[string]interface{}{
							"addresses": "test@e154.ru,test2@e154.ru",
							"subject":   "subject",
							"body":      "body",
						},
					})

					time.Sleep(time.Millisecond * 500)

				})
			})

		})
	})
}
