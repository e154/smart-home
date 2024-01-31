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

package messagebird

import (
	"context"
	"testing"
	"time"

	notifyCommon "github.com/e154/smart-home/plugins/notify/common"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/messagebird"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestMessagebird(t *testing.T) {

	Convey("messagbird", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			AddPlugin(adaptors, "messagebird")

			settings := messagebird.NewSettings()
			settings[messagebird.AttrAccessKey].Value = "XXXX"
			settings[messagebird.AttrName].Value = "YYYY"

			sensorEnt := &m.Entity{
				Id:         common.EntityId("messagebird.messagebird"),
				PluginName: "messagebird",
				AutoLoad:   true,
			}
			sensorEnt.Settings = settings
			err := adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "notify", "messagebird")
			supervisor.Start(context.Background())
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					const (
						phone = "+79990000001"
						body  = "some text"
					)

					eventBus.Publish(notify.TopicNotify, notifyCommon.Message{
						EntityId: common.NewEntityId("messagebird.messagebird"),
						Attributes: map[string]interface{}{
							messagebird.AttrPhone: phone,
							messagebird.AttrBody:  body,
						},
					})

					//todo: fix
					time.Sleep(time.Millisecond * 100)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					del := list[0]
					ctx.So(del.Status, ShouldEqual, m.MessageStatusInProgress)
					ctx.So(del.Address, ShouldEqual, phone)
					ctx.So(del.ErrorMessageBody, ShouldBeNil)
					ctx.So(del.ErrorMessageStatus, ShouldBeNil)
					ctx.So(del.Message.Type, ShouldEqual, messagebird.Name)

					attr := messagebird.NewMessageParams()
					_, _ = attr.Deserialize(del.Message.Attributes)
					ctx.So(attr[messagebird.AttrPhone].String(), ShouldEqual, phone)
					ctx.So(attr[messagebird.AttrBody].String(), ShouldEqual, body)

				})
			})

		})
	})
}
