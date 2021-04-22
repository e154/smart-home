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

package models

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	Convey("add user", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			accessList access_list.AccessListService) {

			// clear database
			migrations.Purge()

			// add role
			userRole := &m.Role{
				Name: "user_role",
			}
			err := adaptors.Role.Add(userRole)
			So(err, ShouldBeNil)

			var counter int
			for pack, item := range *accessList.List {
				for right := range item {
					if strings.Contains(right, "read") ||
						strings.Contains(right, "view") ||
						strings.Contains(right, "show") {
						permission := &m.Permission{
							RoleName:    userRole.Name,
							PackageName: pack,
							LevelName:   right,
						}

						counter++
						_, err = adaptors.Permission.Add(permission)
						So(err, ShouldBeNil)
					}
				}
			}

			const (
				nickname = "user"
				email    = "email@mail.com"
			)

			// add user
			user := &m.User{
				Nickname: nickname,
				RoleName: "user_role",
				Email:    email,
				Lang:     "en",
				Meta: []*m.UserMeta{
					&m.UserMeta{
						Key:   "phone1",
						Value: "+18004001234",
					},
				},
			}
			err = user.SetPass("123456")
			So(err, ShouldBeNil)

			ok, _ := user.Valid()
			So(ok, ShouldEqual, true)

			user.Id, err = adaptors.User.Add(user)
			So(err, ShouldBeNil)

			user, err = adaptors.User.GetById(user.Id)
			So(err, ShouldBeNil)

			//debug.Println(user)
			//fmt.Println("----")

			So(user.Nickname, ShouldEqual, nickname)
			So(user.Email, ShouldEqual, email)
			So(user.Role, ShouldNotBeNil)
			So(user.Role.Name, ShouldEqual, "user_role")
		})
	})
}
