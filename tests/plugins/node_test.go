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
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/system/event_bus/events"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {

	Convey("node", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "node")
			ctx.So(err, ShouldBeNil)

			go mqttServer.Start()

			// add entity
			// ------------------------------------------------

			nodeEnt := GetNewNode("main")
			err = adaptors.Entity.Add(nodeEnt)
			ctx.So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.SetPluginManager(pluginManager)
			entityManager.LoadEntities()

			defer func() {
				mqttServer.Shutdown()
				zigbee2mqtt.Shutdown()
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

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
					eventBus.Subscribe(event_bus.TopicEntities, fn)
					defer eventBus.Unsubscribe(event_bus.TopicEntities, fn)

					b, err := json.Marshal(node.MessageStatus{
						Status:    "enabled",
						Thread:    1,
						Rps:       2,
						Min:       3,
						Max:       4,
						StartedAt: now,
					})
					ctx.So(err, ShouldBeNil)
					err = mqttServer.Publish("home/node/main/ping", b, 0, false)
					ctx.So(err, ShouldBeNil)

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

				})
			})

			mqttCli := mqttServer.NewClient("cli")

			t.Run("request", func(t *testing.T) {
				Convey("case", t, func(ctx C) {
					ch := make(chan struct{})
					mqttCli.Subscribe("home/node/main/req/#", func(client mqtt.MqttCli, message mqtt.Message) {
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
					eventBus.Publish(fmt.Sprintf("plugin.node/main/req/%s", nodeEnt.Id), req)

					ticker := time.NewTimer(time.Second * 1)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

				})
			})

			t.Run("response", func(t *testing.T) {
				Convey("case", t, func(ctx C) {

					ch := make(chan struct{})
					topic := fmt.Sprintf("plugin.node/main/resp/%s", "plugin.test")
					fn := func(topic string, resp node.MessageResponse) {
						ctx.So(topic, ShouldEqual, "plugin.node/main/resp/plugin.test")
						ctx.So(resp.EntityId, ShouldEqual, "plugin.test")
						ctx.So(resp.DeviceType, ShouldEqual, "test")
						ctx.So(resp.Status, ShouldEqual, "success")
						ch <- struct{}{}
					}
					eventBus.Subscribe(topic, fn)
					defer eventBus.Unsubscribe(topic, fn)
					b, err := json.Marshal(node.MessageResponse{
						EntityId:   "plugin.test",
						DeviceType: "test",
						Properties: nil,
						Response:   nil,
						Status:     "success",
					})
					ctx.So(err, ShouldBeNil)

					mqttCli.Publish(fmt.Sprintf("home/node/main/resp/%s", "plugin.test"), b)

					ticker := time.NewTimer(time.Second * 1)
					defer ticker.Stop()

					var ok bool
					select {
					case <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					time.Sleep(time.Millisecond * 500)
				})
			})
		})
	})
}
