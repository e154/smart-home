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
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/messagebird"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMessagebird(t *testing.T) {

	Convey("messagbird", t, func(ctx C) {
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
			settings := messagebird.NewSettings()
			settings[messagebird.AttrAccessKey].Value = "XXXX"
			settings[messagebird.AttrName].Value = "YYYY"
			err = AddPlugin(adaptors, "messagebird", settings)
			ctx.So(err, ShouldBeNil)

			pluginManager.Start()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()

			defer func() {
				entityManager.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					ch := make(chan interface{}, 1)
					fn := func(topic string, message interface{}) {
						switch v := message.(type) {
						case event_bus.EventStateChanged:
							ch <- v
						default:
						}

					}
					eventBus.Subscribe(event_bus.TopicEntities, fn)
					defer eventBus.Unsubscribe(event_bus.TopicEntities, fn)

					const (
						phone = "+79990000001"
						body  = "some text"
					)

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: messagebird.Name,
						Attributes: map[string]interface{}{
							messagebird.AttrPhone: phone,
							messagebird.AttrBody:  body,
						},
					})

					ok := Wait(5, ch)

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Second * 2)

					list, total, err := adaptors.MessageDelivery.List(10, 0, "", "")
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					del := list[0]
					ctx.So(del.Status, ShouldEqual, m.MessageStatusSucceed)
					ctx.So(del.Address, ShouldEqual, phone)
					ctx.So(del.ErrorMessageBody, ShouldBeNil)
					ctx.So(del.ErrorMessageStatus, ShouldBeNil)
					ctx.So(del.Message.Type, ShouldEqual, messagebird.Name)

					attr := messagebird.NewMessageParams()
					attr.Deserialize(del.Message.Attributes)
					ctx.So(attr[messagebird.AttrPhone].String(), ShouldEqual, phone)
					ctx.So(attr[messagebird.AttrBody].String(), ShouldEqual, body)

				})
			})

		})
	})
}
