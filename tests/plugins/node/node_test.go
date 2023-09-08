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

package node

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {

	Convey("node", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus) {

			// register plugins
			err := AddPlugin(adaptors, "node")
			ctx.So(err, ShouldBeNil)

			go mqttServer.Start()
			automation.Start()
			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			// add entity
			// ------------------------------------------------

			nodeEnt := GetNewNode("main")
			err = adaptors.Entity.Add(nodeEnt)
			ctx.So(err, ShouldBeNil)

			eventBus.Publish("system/entities/"+nodeEnt.Id.String(), events.EventCreatedEntity{
				EntityId: nodeEnt.Id,
			})

			time.Sleep(time.Second)

			// ------------------------------------------------

			now := time.Now()

			t.Run("ping", func(t *testing.T) {
				Convey("case", t, func(ctx C) {
					ch := make(chan struct{})
					fn := func(topic string, msg interface{}) {
						switch v := msg.(type) {
						case events.EventStateChanged:
							if v.PluginName != "node" {
								return
							}
							ctx.So(v.OldState.State, ShouldNotBeNil)
							ctx.So(v.OldState.State.Name, ShouldEqual, "wait")
							ctx.So(v.NewState.State, ShouldNotBeNil)
							ctx.So(v.NewState.State.Name, ShouldEqual, "connected")
							ctx.So(v.NewState.Attributes[node.AttrThread].Int64(), ShouldEqual, 1)
							ctx.So(v.NewState.Attributes[node.AttrRps].Int64(), ShouldEqual, 2)
							ctx.So(v.NewState.Attributes[node.AttrMin].Int64(), ShouldEqual, 3)
							ctx.So(v.NewState.Attributes[node.AttrMax].Int64(), ShouldEqual, 4)
							ctx.So(v.NewState.Attributes[node.AttrStartedAt].Time(), ShouldEqual, now)
							ch <- struct{}{}
						}
					}
					_ = eventBus.Subscribe("system/entities/+", fn)
					defer func() { _ = eventBus.Unsubscribe("system/entities/+", fn) }()

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

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
					case <-ticker.C:
					}

					ctx.So(ok, ShouldBeTrue)

				})
			})

			mqttCli := mqttServer.NewClient("cli")

			t.Run("request", func(t *testing.T) {
				Convey("case", t, func(ctx C) {
					ch := make(chan struct{})
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

					ticker := time.NewTimer(time.Second * 1)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
					case <-ticker.C:
					}

					ctx.So(ok, ShouldBeTrue)

				})
			})

			t.Run("response", func(t *testing.T) {
				Convey("case", t, func(ctx C) {

					ch := make(chan struct{})
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

					ticker := time.NewTimer(time.Second * 1)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
					case <-ticker.C:
					}

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Millisecond * 500)
				})
			})
		})
	})
}
