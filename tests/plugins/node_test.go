// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
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

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = adaptors.Plugin.CreateOrUpdate(m.Plugin{
				Name:    "node",
				Version: "0.0.1",
				Enabled: true,
				System:  true,
			})
			So(err, ShouldBeNil)

			go mqttServer.Start()

			// add entity
			// ------------------------------------------------

			nodeEnt := GetNewNode()
			err = adaptors.Entity.Add(nodeEnt)
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)

			time.Sleep(time.Millisecond * 500)

			now := time.Now()

			// ping
			// ------------------------------------------------
			wgPing := sync.WaitGroup{}
			wgPing.Add(1)
			eventBus.Subscribe(event_bus.TopicEntities, func(msg interface{}) {
				switch v := msg.(type) {
				case event_bus.EventStateChanged:
					if v.Type != "node" {
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
					wgPing.Done()
				}
			})

			b, err := json.Marshal(node.MessageStatus{
				Status:    "enabled",
				Thread:    1,
				Rps:       2,
				Min:       3,
				Max:       4,
				StartedAt: now,
			})
			So(err, ShouldBeNil)
			err = mqttServer.Publish("home/node/main/ping", b, 0, false)
			So(err, ShouldBeNil)

			wgPing.Wait()

			// request
			// ------------------------------------------------
			wgRequest := sync.WaitGroup{}
			wgRequest.Add(1)
			mqttCli := mqttServer.NewClient("cli")
			mqttCli.Subscribe("home/node/main/req/#", func(client mqtt.MqttCli, message mqtt.Message) {
				req := node.MessageRequest{}
				err = json.Unmarshal(message.Payload, &req)
				ctx.So(err, ShouldBeNil)
				ctx.So(req.EntityId, ShouldEqual, "test.test")
				ctx.So(req.DeviceType, ShouldEqual, "test")
				wgRequest.Done()
			})

			req := node.MessageRequest{
				EntityId:   "test.test",
				DeviceType: "test",
				Properties: nil,
				Command:    nil,
			}
			So(err, ShouldBeNil)
			eventBus.Publish("plugin.node/main/req", req)

			wgRequest.Wait()

			// response
			// ------------------------------------------------
			wgResp := sync.WaitGroup{}
			wgResp.Add(1)
			eventBus.Subscribe("plugin.node/main/resp", func(resp node.MessageResponse) {
				ctx.So(resp.EntityId, ShouldEqual, "test.test")
				ctx.So(resp.DeviceType, ShouldEqual, "test")
				ctx.So(resp.Status, ShouldEqual, "success")
				wgResp.Done()
			})
			b, err = json.Marshal(node.MessageResponse{
				EntityId:   "test.test",
				DeviceType: "test",
				Properties: nil,
				Response:   nil,
				Status:     "success",
			})
			So(err, ShouldBeNil)
			mqttCli.Publish("home/node/main/resp", b)

			wgResp.Wait()

			time.Sleep(time.Millisecond * 500)

			mqttServer.Shutdown()
			zigbee2mqtt.Shutdown()
			entityManager.Shutdown()
			automation.Shutdown()
			pluginManager.Shutdown()
		})
	})
}
