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

	"github.com/e154/smart-home/internal/plugins/notify"
	notifyCommon "github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/internal/plugins/telegram"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
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
			supervisor plugins.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			AddPlugin(adaptors, "telegram")

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "notify", "telegram")
			supervisor.Start(context.Background())
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			// add scripts
			// ------------------------------------------------

			plugScript, err := AddScript("telegram script", sourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			tgEnt := GetNewTelegram("clavicus")
			tgEnt.Actions = []*models.EntityAction{
				{
					Name:        "CHECK",
					Description: "check status",
					Script:      plugScript,
				},
			}
			err = adaptors.Entity.Add(context.Background(), tgEnt)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &models.EntityStorage{
				EntityId:   tgEnt.Id,
				Attributes: tgEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+tgEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: tgEnt.Id,
			})

			time.Sleep(time.Second)

			// add chat
			tgChan := models.TelegramChat{
				EntityId: tgEnt.Id,
				ChatId:   123,
				Username: "user",
			}
			_ = adaptors.TelegramChat.Add(context.Background(), tgChan)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					time.Sleep(time.Millisecond * 500)

					eventBus.Publish(notify.TopicNotify, notifyCommon.Message{
						EntityId: common.NewEntityId("telegram.clavicus"),
						Attributes: map[string]interface{}{
							"chat_id": 123,
							"body":    "body",
						},
					})

					//todo: fix
					time.Sleep(time.Millisecond * 500)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					ctx.So(list[0].Status, ShouldEqual, models.MessageStatusSucceed)
					ctx.So(list[0].Address, ShouldEqual, "123")
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
