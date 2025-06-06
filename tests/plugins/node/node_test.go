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

package node

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/internal/plugins/node"
	"github.com/e154/smart-home/internal/system/zigbee2mqtt"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {

	Convey("node", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			eventBus bus.Bus) {

			// register plugins
			err := AddPlugin(adaptors, "node")
			ctx.So(err, ShouldBeNil)

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor", "Mqtt")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "node")
			go mqttServer.Start()
			supervisor.Start(context.Background())
			defer mqttServer.Shutdown()
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			// add entity
			// ------------------------------------------------

			nodeEnt := GetNewNode("main")
			err = adaptors.Entity.Add(context.Background(), nodeEnt)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+nodeEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: nodeEnt.Id,
			})

			time.Sleep(time.Second)

			// common
			// ------------------------------------------------
			var closed = false
			ch := make(chan events.EventStateChanged)
			defer close(ch)
			fn := func(topic string, msg interface{}) {
				switch v := msg.(type) {
				case events.EventStateChanged:
					if !closed {
						ch <- v
					}
				}
			}
			eventBus.Subscribe("system/entities/+", fn, false)
			defer eventBus.Unsubscribe("system/entities/+", fn)
			defer func() {
				closed = true
			}()
			// ------------------------------------------------

			// wait message
			_, ok := WaitT[events.EventStateChanged](time.Second*2, ch)
			ctx.So(ok, ShouldBeTrue)

			now := time.Now()

			t.Run("ping", func(t *testing.T) {
				Convey("case", t, func(ctx C) {

					b, err := json.Marshal(node.MessageStatus{
						Status:    "enabled",
						Thread:    1,
						Rps:       2,
						Min:       3,
						Max:       4,
						StartedAt: now,
					})
					ctx.So(err, ShouldBeNil)
					err = mqttServer.Publish("system/plugins/node/main/ping", b, 0, false)
					ctx.So(err, ShouldBeNil)

					// wait message
					msg, ok := WaitT[events.EventStateChanged](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					ctx.So(msg.OldState.State, ShouldNotBeNil)
					ctx.So(msg.OldState.State.Name, ShouldEqual, "wait")
					ctx.So(msg.NewState.State, ShouldNotBeNil)
					ctx.So(msg.NewState.State.Name, ShouldEqual, "connected")
					ctx.So(msg.NewState.Attributes[node.AttrThread].Int64(), ShouldEqual, 1)
					ctx.So(msg.NewState.Attributes[node.AttrRps].Int64(), ShouldEqual, 2)
					ctx.So(msg.NewState.Attributes[node.AttrMin].Int64(), ShouldEqual, 3)
					ctx.So(msg.NewState.Attributes[node.AttrMax].Int64(), ShouldEqual, 4)
					ctx.So(msg.NewState.Attributes[node.AttrStartedAt].Time(), ShouldEqual, now)
					ctx.So(ok, ShouldBeTrue)
				})
			})

			mqttCli := mqttServer.NewClient("cli")

			t.Run("request", func(t *testing.T) {
				Convey("case", t, func(ctx C) {
					ch := make(chan struct{})
					defer close(ch)
					_ = mqttCli.Subscribe("system/plugins/node/main/req/#", func(client mqtt.MqttCli, message mqtt.Message) {
						req := node.MessageRequest{}
						err = json.Unmarshal(message.Payload, &req)
						ctx.So(err, ShouldBeNil)
						ctx.So(req.EntityId, ShouldEqual, "plugin.test")
						ctx.So(req.DeviceType, ShouldEqual, "test")
						ch <- struct{}{}
					})
					defer mqttCli.Unsubscribe("home/node/main/req/#")

					req := node.MessageRequest{
						EntityId:   "plugin.test",
						DeviceType: "test",
						Properties: nil,
						Command:    nil,
					}
					eventBus.Publish(fmt.Sprintf("system/plugins/node/main/req/%s", nodeEnt.Id), req)

					// wait message
					_, ok := WaitT[struct{}](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

				})
			})

			t.Run("response", func(t *testing.T) {
				Convey("case", t, func(ctx C) {

					ch := make(chan struct{})
					defer close(ch)
					topic := fmt.Sprintf("system/plugins/node/main/resp/%s", "plugin.test")
					fn := func(topic string, resp node.MessageResponse) {
						ctx.So(topic, ShouldEqual, "system/plugins/node/main/resp/plugin.test")
						ctx.So(resp.EntityId, ShouldEqual, "plugin.test")
						ctx.So(resp.DeviceType, ShouldEqual, "test")
						ctx.So(resp.Status, ShouldEqual, "success")
						ch <- struct{}{}
					}
					_ = eventBus.Subscribe(topic, fn)
					defer func() { _ = eventBus.Unsubscribe(topic, fn) }()

					time.Sleep(time.Millisecond * 500)

					b, err := json.Marshal(node.MessageResponse{
						EntityId:   "plugin.test",
						DeviceType: "test",
						Properties: nil,
						Response:   nil,
						Status:     "success",
					})
					ctx.So(err, ShouldBeNil)

					_ = mqttCli.Publish(fmt.Sprintf("system/plugins/node/main/resp/%s", "plugin.test"), b)

					// wait message
					_, ok := WaitT[struct{}](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Millisecond * 500)
				})
			})
		})
	})
}
