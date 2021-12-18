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
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/debug"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/jwt_manager"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJwt(t *testing.T) {

	t.Run("jwt", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			const hmac = "0a84d2275e242d99ceeaa93820bf64232435ab297d2d513d07596e661cc01957"

			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				jwtManager jwt_manager.JwtManager) {
				jwtManager.Start()
				b, _ := hex.DecodeString(hmac)
				jwtManager.SetHmacKey(b)

				t.Run("generate", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						user := &m.User{
							Id:       1,
							Nickname: "John Doe",
							RoleName: "user",
						}
						accessToken, err := jwtManager.Generate(user)
						fmt.Println(accessToken)
						ctx.So(err, ShouldBeNil)
						ctx.So(accessToken, ShouldNotBeBlank)
					})
				})

				t.Run("verify", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0Mjg1MDUyOTksImkiOjEsImlhdCI6MTYzOTU4Njg5OSwiaXNzIjoic2VydmVyIiwibiI6IkpvaG4gRG9lIiwibmJmIjoxNjM5NTg2ODk5LCJyIjoidXNlciJ9.QI35ZYqSsQkP7SW2gp4VCPXfs9jh6DOHfP7WKlPU71A"

						claims, err := jwtManager.Verify(accessToken)
						debug.Println(claims)
						ctx.So(err, ShouldBeNil)
						ctx.So(claims, ShouldNotBeNil)
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
						accessToken, err := jwtManager.Generate(user)
						ctx.So(err, ShouldBeNil)
						ctx.So(accessToken, ShouldNotBeBlank)

						claims, err := jwtManager.Verify(accessToken)
						ctx.So(err, ShouldBeNil)
						ctx.So(claims, ShouldNotBeNil)
						ctx.So(claims.UserId, ShouldEqual, 1)
						ctx.So(claims.Username, ShouldEqual, "John Doe")
						ctx.So(claims.RoleName, ShouldEqual, "user")
					})
				})

				t.Run("invalid signature", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						const accessToken1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0MjU2OTQ1MTgsImkiOjEsImlhdCI6MTYzNjc3NjExOCwiaXNzIjoic2VydmVyIiwibiI6IkpvaG4gRG9lIiwibmJmIjoxNjM2Nzc2MTE4LCJyIjoidXNlciJ9.gxLi_hKQvAdkZtydyMRCje228u3Y8Xiad-iJM-U8E38q"
						claims, err := jwtManager.Verify(accessToken1)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token claims")

						const accessToken2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0MjU2OTQ1MTgsImkiOjEsImlhdCI6MTYzNjc3NjExOCwiaXNzIjoic2VydmVyIiwibiI6IkpvaG4gRG9lIiwibmJmIjoxNjM2Nzc2MTE4LCJyIjoidXNlciJ9q.gxLi_hKQvAdkZtydyMRCje228u3Y8Xiad-iJM-U8E38"
						claims, err = jwtManager.Verify(accessToken2)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token claims")

						const accessToken3 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9q.eyJleHAiOjI0MjU2OTQ1MTgsImkiOjEsImlhdCI6MTYzNjc3NjExOCwiaXNzIjoic2VydmVyIiwibiI6IkpvaG4gRG9lIiwibmJmIjoxNjM2Nzc2MTE4LCJyIjoidXNlciJ9.gxLi_hKQvAdkZtydyMRCje228u3Y8Xiad-iJM-U8E38"
						claims, err = jwtManager.Verify(accessToken3)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid token claims")

						const accessToken4 = "sometext"
						claims, err = jwtManager.Verify(accessToken4)
						ctx.So(err, ShouldNotBeNil)
						ctx.So(claims, ShouldBeNil)
						ctx.So(err.Error(), ShouldEqual, "invalid access token")
					})
				})
			})
			ctx.So(err, ShouldBeNil)
		})
	})
}
