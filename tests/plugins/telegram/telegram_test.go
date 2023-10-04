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

package telegram

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/telegram"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTelegram(t *testing.T) {

	const sourceScript = `
checkStatus =->
    #print '----------------1'

telegramAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()
`

	Convey("telegram", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			AddPlugin(adaptors, "telegram")

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			// add scripts
			// ------------------------------------------------

			plugScript, err := AddScript("telegram script", sourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			tgEnt := GetNewTelegram("clavicus")
			tgEnt.Actions = []*m.EntityAction{
				{
					Name:        "CHECK",
					Description: "check status",
					Script:      plugScript,
				},
			}
			err = adaptors.Entity.Add(context.Background(), tgEnt)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &m.EntityStorage{
				EntityId:   tgEnt.Id,
				Attributes: tgEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+tgEnt.Id.String(), events.EventCreatedEntity{
				EntityId: tgEnt.Id,
			})

			time.Sleep(time.Second)

			// add chat
			tgChan := m.TelegramChat{
				EntityId: tgEnt.Id,
				ChatId:   123,
				Username: "user",
			}
			_ = adaptors.TelegramChat.Add(context.Background(), tgChan)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					ch := make(chan interface{}, 2)
					fn := func(topic string, message interface{}) {
						switch v := message.(type) {
						case events.EventStateChanged:
							ch <- v
						default:
						}

					}
					_ = eventBus.Subscribe("system/entities/+", fn)

					time.Sleep(time.Millisecond * 500)

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: telegram.Name,
						Attributes: map[string]interface{}{
							"name": "clavicus",
							"body": "body",
						},
					})

					ok := Wait(3, ch)
					ctx.So(ok, ShouldBeTrue)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					ctx.So(list[0].Status, ShouldEqual, m.MessageStatusSucceed)
					ctx.So(list[0].Address, ShouldBeIn, []string{"clavicus"})
					ctx.So(list[0].ErrorMessageBody, ShouldBeNil)
					ctx.So(list[0].ErrorMessageStatus, ShouldBeNil)
					ctx.So(list[0].Message.Type, ShouldEqual, telegram.Name)

					attr := telegram.NewMessageParams()
					_, _ = attr.Deserialize(list[0].Message.Attributes)
					ctx.So(attr[telegram.AttrBody].String(), ShouldEqual, "body")

				})
			})

			t.Run("call actions", func(t *testing.T) {
				Convey("call actions", t, func(ctx C) {
					supervisor.CallAction(tgEnt.Id, "CHECK", nil)
					time.Sleep(time.Second)
				})
			})
		})
	})
}
