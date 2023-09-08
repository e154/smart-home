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

package email

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/email"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmail(t *testing.T) {

	Convey("email", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			settings := email.NewSettings()
			settings[email.AttrAuth].Value = "XXX"
			settings[email.AttrPass].Value = "XXX"
			settings[email.AttrSmtp].Value = "XXX"
			settings[email.AttrPort].Value = 123
			settings[email.AttrSender].Value = "XXX"
			AddPlugin(adaptors, "email", settings.Serialize())

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

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
					eventBus.Subscribe("system/entities/+", fn)
					defer eventBus.Unsubscribe("system/entities/+", fn)

					eventBus.Publish(notify.TopicNotify, notify.Message{
						Type: email.Name,
						Attributes: map[string]interface{}{
							"addresses": "test@e154.ru,test2@e154.ru",
							"subject":   "subject",
							"body":      "body",
						},
					})

					ok := Wait(5, ch)

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Second * 2)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 2)

					for _, del := range list {
						ctx.So(del.Status, ShouldEqual, m.MessageStatusSucceed)
						ctx.So(del.Address, ShouldBeIn, []string{"test@e154.ru", "test2@e154.ru"})
						ctx.So(del.ErrorMessageBody, ShouldBeNil)
						ctx.So(del.ErrorMessageStatus, ShouldBeNil)
						ctx.So(del.Message.Type, ShouldEqual, email.Name)

						attr := email.NewMessageParams()
						_, _ = attr.Deserialize(del.Message.Attributes)
						ctx.So(attr[email.AttrSubject].String(), ShouldEqual, "subject")
						ctx.So(attr[email.AttrBody].String(), ShouldEqual, "body")
					}
				})
			})
		})
	})
}
