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

package twilio

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/twilio"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTwilio(t *testing.T) {

	Convey("twilio", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			settings := twilio.NewSettings()
			settings[twilio.AttrAuthToken].Value = "XXXX"
			settings[twilio.AttrSid].Value = "YYYY"
			settings[twilio.AttrFrom].Value = "YYYY"
			AddPlugin(adaptors, "twilio", settings.Serialize())


			supervisor.Start(context.Background())

			time.Sleep(time.Millisecond * 500)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					ch := make(chan interface{}, 1)
					fn := func(topic string, message interface{}) {
						switch v := message.(type) {
						case events.EventStateChanged:
							ch <- v
						default:
						}

					}
					_ = eventBus.Subscribe("system/entities/+", fn)
					defer func() { _ = eventBus.Unsubscribe("system/entities/+", fn) }()

					const (
						phone = "+79990000001"
						body  = "some text"
					)

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: twilio.Name,
						Attributes: map[string]interface{}{
							twilio.AttrPhone: phone,
							twilio.AttrBody:  body,
						},
					})

					ok := Wait(5, ch)

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Second * 2)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 1)

					del := list[0]
					ctx.So(del.Status, ShouldEqual, m.MessageStatusSucceed)
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
