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

func TestAuth(t *testing.T) {

	//Convey("auth", t, func(ctx C) {
	//	err := container.Invoke(func(adaptors *adaptors.Adaptors,
	//		migrations *migrations.Migrations,
	//		scriptService scripts.ScriptService,
	//		eventBus bus.Bus,
	//		controllers *controllers.Controllers,
	//		dialer *container2.Dialer) {
	//
	//		//eventBus.Restart()
	//		//scriptService.Restart()
	//		//
	//		//err := migrations.Restart()
	//		//ctx.So(err, ShouldBeNil)
	//		//
	//		//c := context.Background()
	//		//conn, err := grpc.DialContext(c, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer.Call()))
	//		//if err != nil {
	//		//	log.Fatal(err)
	//		//}
	//		//defer conn.Close()
	//		//
	//		//client := gw.NewAuthServiceClient(conn)
	//
	//		t.Run("signin", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("signout", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				ctx.Println("test not implemented")
	//			})
	//		})
	//
	//		t.Run("access_list", func(t *testing.T) {
	//			Convey("", t, func(ctx C) {
	//				ctx.Println("test not implemented")
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
