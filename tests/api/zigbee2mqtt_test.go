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

package api

import (
	"testing"
)

func TestZigbee2mqtt(t *testing.T) {

	//Convey("zigbee2mqtt", t, func(ctx C) {
	//	err := container.Invoke(func(adaptors *adaptors.Adaptors,
	//		migrations *migrations.Migrations,
	//		scriptService scripts.ScriptService,
	//		eventBus bus.Bus,
	//		controllers *controllers.Controllers,
	//		dialer *container2.Dialer) {
	//
	//		eventBus.Purge()
	//		scriptService.Restart()
	//
	//		err := migrations.Purge()
	//		ctx.So(err, ShouldBeNil)
	//
	//		c := context.Background()
	//		conn, err := grpc.DialContext(c, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer.Call()))
	//		ctx.So(err, ShouldBeNil)
	//		defer func() { _ = conn.Close() }()
	//
	//		client := gw.NewZigbee2MqttServiceClient(conn)
	//
	//		t.Run("success", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//
	//				request := &gw.NewZigbee2MqttRequest{
	//					Name:       "zigbee2mqtt",
	//					Login:      "zigbee2mqtt",
	//					Password:   common.String("zigbee2mqtt"),
	//					PermitJoin: true,
	//					BaseTopic:  "zigbee2mqtt",
	//				}
	//
	//				bridge, err := client.AddZigbee2MqttBridge(c, request)
	//				ctx.So(err, ShouldBeNil)
	//				ctx.So(bridge, ShouldNotBeNil)
	//
	//				t.Run("getById", func(t *testing.T) {
	//					Convey("", t, func(ctx C) {
	//
	//						request := &gw.GetBridgeRequest{Id: bridge.Id}
	//						bridge, err := client.GetZigbee2MqttBridge(c, request)
	//						ctx.So(err, ShouldBeNil)
	//						ctx.So(bridge, ShouldNotBeNil)
	//						ctx.So(bridge.Id, ShouldEqual, 1)
	//						ctx.So(bridge.Name, ShouldEqual, "zigbee2mqtt")
	//						ctx.So(bridge.Login, ShouldEqual, "zigbee2mqtt")
	//						ctx.So(bridge.BaseTopic, ShouldEqual, "zigbee2mqtt")
	//						ctx.So(bridge.PermitJoin, ShouldEqual, true)
	//					})
	//				})
	//
	//				t.Run("update", func(t *testing.T) {
	//					Convey("", t, func(ctx C) {
	//
	//						// failed
	//						request := &gw.UpdateBridgeRequest{
	//							Id:         bridge.Id + 1,
	//							Name:       "new_zigbee2mqtt",
	//							Login:      "new_zigbee2mqtt",
	//							Password:   common.String("new_zigbee2mqtt"),
	//							PermitJoin: false,
	//							BaseTopic:  "new_zigbee2mqtt",
	//						}
	//						_, err := client.UpdateBridgeById(c, request)
	//						ctx.So(err, ShouldNotBeNil)
	//						ctx.So(err.Error(), ShouldEqual, "rpc error: code = NotFound desc = not found")
	//
	//						// success
	//						request.Id = bridge.Id
	//						bridge, err := client.UpdateBridgeById(c, request)
	//						ctx.So(err, ShouldBeNil)
	//						ctx.So(bridge, ShouldNotBeNil)
	//						ctx.So(bridge.Id, ShouldEqual, 1)
	//						ctx.So(bridge.Name, ShouldEqual, "zigbee2mqtt")
	//						ctx.So(bridge.Login, ShouldEqual, "new_zigbee2mqtt")
	//						ctx.So(bridge.BaseTopic, ShouldEqual, "new_zigbee2mqtt")
	//						ctx.So(bridge.PermitJoin, ShouldEqual, false)
	//					})
	//				})
	//			})
	//		})
	//
	//		t.Run("update", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				_, _ = ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("list", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				_, _ = ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("delete", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				_, _ = ctx.Println("test not implemented")
	//			})
	//		})
	//
	//	})
	//
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//})
}
