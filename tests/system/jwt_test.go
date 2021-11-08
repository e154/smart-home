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

package system

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/debug"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/jwt_manager"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJwt(t *testing.T) {

	t.Run("jwt", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				manager jwt_manager.JwtManager) {
				manager.Start()

				t.Run("generate", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						user := &m.User{
							Id:       1,
							Nickname: "John Doe",
							RoleName: "user",
						}
						accessToken, err := manager.Generate(user)
						fmt.Println(accessToken)
						ctx.So(err, ShouldBeNil)
						ctx.So(accessToken, ShouldNotBeBlank)
					})
				})

				t.Run("verify", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzg5ODA2NjcsImlhdCI6MTYzNjM4ODY2NywiaXNzIjoic2VydmVyIiwibmJmIjoxNjM2Mzg4NjY3LCJpIjoxLCJuIjoiSm9obiBEb2UiLCJyIjoidXNlciJ9.RxHgi86tgXJg_I_1ZCxDYZOdmldgDWnR5wGi1pgF4ig"

						claims, err := manager.Verify(accessToken)
						debug.Println(claims)
						ctx.So(err, ShouldBeNil)
						ctx.So(claims, ShouldNotBeNil)
						ctx.So(claims.ExpiresAt, ShouldEqual, 1638980667)
						ctx.So(claims.IssuedAt, ShouldEqual, 1636388667)
						ctx.So(claims.Issuer, ShouldEqual, "server")
						ctx.So(claims.NotBefore, ShouldEqual, 1636388667)
						ctx.So(claims.UserId, ShouldEqual, 1)
						ctx.So(claims.Username, ShouldEqual, "John Doe")
						ctx.So(claims.RoleName, ShouldEqual, "user")
					})
				})

				t.Run("generate + verify", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						user := &m.User{
							Id:       1,
							Nickname: "John Doe",
							RoleName: "user",
						}
						accessToken, err := manager.Generate(user)
						ctx.So(err, ShouldBeNil)
						ctx.So(accessToken, ShouldNotBeBlank)

						claims, err := manager.Verify(accessToken)
						ctx.So(err, ShouldBeNil)
						ctx.So(claims, ShouldNotBeNil)
						//ctx.So(claims.ExpiresAt, ShouldEqual, 1626345821)
						//ctx.So(claims.IssuedAt, ShouldEqual, 1623753821)
						ctx.So(claims.Issuer, ShouldEqual, "server")
						//ctx.So(claims.NotBefore, ShouldEqual, 1623753821)
						ctx.So(claims.UserId, ShouldEqual, 1)
						ctx.So(claims.Username, ShouldEqual, "John Doe")
						ctx.So(claims.RoleName, ShouldEqual, "user")

						err = claims.Valid()
						ctx.So(err, ShouldBeNil)
					})
				})

				t.Run("invalid signature", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						const accessToken1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjYzNDU4MjEsImlhdCI6MTYyMzc1MzgyMSwiaXNzIjoic2VydmVyIiwibmJmIjoxNjIzNzUzODIxLCJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6InVzZXIifQ.PlnyM928_KJaeBseB5IpCrphu4T7O4y-2oK0SeUgv8Qq"
						claims, err := manager.Verify(accessToken1)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token: signature is invalid")

						const accessToken2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjYzNDU4MjEsImlhdCI6MTYyMzc1MzgyMSwiaXNzIjoic2VydmVyIiwibmJmIjoxNjIzNzUzODIxLCJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6InVzZXIifQq.PlnyM928_KJaeBseB5IpCrphu4T7O4y-2oK0SeUgv8Q"
						claims, err = manager.Verify(accessToken2)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token: signature is invalid")

						const accessToken3 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9q.eyJleHAiOjE2MjYzNDU4MjEsImlhdCI6MTYyMzc1MzgyMSwiaXNzIjoic2VydmVyIiwibmJmIjoxNjIzNzUzODIxLCJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6InVzZXIifQ.PlnyM928_KJaeBseB5IpCrphu4T7O4y-2oK0SeUgv8Q"
						claims, err = manager.Verify(accessToken3)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token: illegal base64 data at input byte 37")
					})
				})
			})
			ctx.So(err, ShouldBeNil)
		})
	})
}
