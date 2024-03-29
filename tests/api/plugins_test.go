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

func TestPlugins(t *testing.T) {

	//Convey("plugins", t, func(ctx C) {
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
	//		client := gw.NewPluginServiceClient(conn)
	//
	//		t.Run("list", func(t *testing.T) {
	//			Convey("list", t, func(ctx C) {
	//
	//				request := &gw.PaginationRequest{
	//					Limit: 100,
	//					Page:  1,
	//					Sort:  "+name",
	//				}
	//
	//				response, err := client.GetPluginList(c, request)
	//				ctx.So(err, ShouldBeNil)
	//
	//				if response != nil {
	//					ctx.So(len(response.Items), ShouldEqual, 0)
	//				}
	//
	//			})
	//		})
	//
	//		t.Run("enable", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				_, _ = ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("disabled", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				_, _ = ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("options", func(t *testing.T) {
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
