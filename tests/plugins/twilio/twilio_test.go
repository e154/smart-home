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

package twilio

import (
	"context"
	"testing"
	"time"

	notifyCommon "github.com/e154/smart-home/plugins/notify/common"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/twilio"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestTwilio(t *testing.T) {

	Convey("twilio", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			AddPlugin(adaptors, "twilio")

			settings := twilio.NewSettings()
			settings[twilio.AttrAuthToken].Value = "XXXX"
			settings[twilio.AttrSid].Value = "YYYY"
			settings[twilio.AttrFrom].Value = "YYYY"

			sensorEnt := &m.Entity{
				Id:         common.EntityId("twilio.twilio"),
				PluginName: "twilio",
				AutoLoad:   true,
			}
			sensorEnt.Settings = settings
			err := adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "twilio", "notify")
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
						EntityId: common.NewEntityId("twilio.twilio"),
						Attributes: map[string]interface{}{
							twilio.AttrPhone: phone,
							twilio.AttrBody:  body,
						},
					})

					//todo: fix
					time.Sleep(time.Millisecond * 1000)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					del := list[0]
					//ctx.So(del.Status, ShouldEqual, m.MessageStatusInProgress)
					ctx.So(del.Address, ShouldEqual, phone)
					ctx.So(del.ErrorMessageBody, ShouldBeNil)
					ctx.So(del.ErrorMessageStatus, ShouldBeNil)
					ctx.So(del.Message.Type, ShouldEqual, twilio.Name)

					attr := twilio.NewMessageParams()
					_, _ = attr.Deserialize(del.Message.Attributes)
					ctx.So(attr[twilio.AttrPhone].String(), ShouldEqual, phone)
					ctx.So(attr[twilio.AttrBody].String(), ShouldEqual, body)

				})
			})

		})
	})
}
