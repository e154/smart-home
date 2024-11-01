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

package email

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/internal/plugins/email"
	"github.com/e154/smart-home/internal/plugins/notify"
	notifyCommon "github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
)

func TestEmail(t *testing.T) {

	Convey("email", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "notify")
			AddPlugin(adaptors, "email")

			settings := email.NewSettings()
			settings[email.AttrAuth].Value = "XXX"
			settings[email.AttrPass].Value = "XXX"
			settings[email.AttrSmtp].Value = "XXX"
			settings[email.AttrPort].Value = 123
			settings[email.AttrSender].Value = "XXX"

			sensorEnt := &models.Entity{
				Id:         common.EntityId("email.email"),
				PluginName: "email",
				AutoLoad:   true,
			}
			sensorEnt.Settings = settings
			err := adaptors.Entity.Add(context.Background(), sensorEnt)
			ctx.So(err, ShouldBeNil)

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "notify", "email")
			supervisor.Start(context.Background())
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			t.Run("succeed", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					eventBus.Publish(notify.TopicNotify, notifyCommon.Message{
						EntityId: common.NewEntityId("email.email"),
						Attributes: map[string]interface{}{
							"addresses": "test@e154.ru,test2@e154.ru",
							"subject":   "subject",
							"body":      "body",
						},
					})

					//todo: fix
					time.Sleep(time.Millisecond * 500)

					list, total, err := adaptors.MessageDelivery.List(context.Background(), 10, 0, "", "", nil)
					ctx.So(err, ShouldBeNil)
					ctx.So(total, ShouldEqual, 2)

					for _, del := range list {
						ctx.So(del.Status, ShouldEqual, models.MessageStatusSucceed)
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
