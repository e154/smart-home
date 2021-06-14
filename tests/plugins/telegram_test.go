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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/telegram"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTelegram(t *testing.T) {

	Convey("telegram", t, func(ctx C) {
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
			err = AddPlugin(adaptors, "telegram")
			ctx.So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			tgEnt := GetNewTelegram("clavicus")
			err = adaptors.Entity.Add(tgEnt)
			ctx.So(err, ShouldBeNil)
			tgChan := m.TelegramChat{
				EntityId: tgEnt.Id,
				ChatId:   123,
				Username: "user",
			}
			adaptors.TelegramChat.Add(tgChan)

			pluginManager.Start()
			entityManager.LoadEntities(pluginManager)

			defer func() {
				entityManager.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					ch := make(chan interface{}, 2)
					fn := func(topic string, message interface{}) {
						switch v := message.(type) {
						case event_bus.EventStateChanged:
							ch <- v
						default:
						}

					}
					eventBus.Subscribe(event_bus.TopicEntities, fn)
					defer eventBus.Unsubscribe(event_bus.TopicEntities, fn)

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: telegram.Name,
						Attributes: map[string]interface{}{
							"name": "clavicus",
							"body": "body",
						},
					})

					ok := Wait(2, ch)

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Millisecond * 500)

					list, total, err := adaptors.MessageDelivery.List(10, 0, "", "")
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					ctx.So(list[0].Status, ShouldEqual, m.MessageStatusSucceed)
					ctx.So(list[0].Address, ShouldBeIn, []string{"clavicus"})
					ctx.So(list[0].ErrorMessageBody, ShouldBeNil)
					ctx.So(list[0].ErrorMessageStatus, ShouldBeNil)
					ctx.So(list[0].Message.Type, ShouldEqual, telegram.Name)

					attr := telegram.NewMessageParams()
					attr.Deserialize(list[0].Message.Attributes)
					ctx.So(attr[telegram.AttrBody].String(), ShouldEqual, "body")

				})
			})
		})
	})
}
