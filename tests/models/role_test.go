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

package models

import (
	"context"
	"strings"
	"testing"

	"github.com/e154/smart-home/internal/system/access_list"
	"github.com/e154/smart-home/internal/system/migrations"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRole(t *testing.T) {

	Convey("add role", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			accessList access_list.AccessListService) {

			// clear database
			_ = migrations.Purge()

			demoRole := &models.Role{
				Name: "demo",
			}
			err := adaptors.Role.Add(context.Background(), demoRole)
			So(err, ShouldBeNil)

			userRole := &models.Role{
				Name:   "user",
				Parent: demoRole,
			}
			err = adaptors.Role.Add(context.Background(), userRole)
			So(err, ShouldBeNil)

			//debug.Println(accessList.List)
			//fmt.Println("----")

			var counter int
			for pack, item := range *accessList.List(context.Background()) {
				for right := range item {
					if strings.Contains(right, "read") ||
						strings.Contains(right, "view") {
						permission := &models.Permission{
							RoleName:    demoRole.Name,
							PackageName: pack,
							LevelName:   right,
						}

						counter++
						_, err = adaptors.Permission.Add(context.Background(), permission)
						So(err, ShouldBeNil)
					}
				}
			}

			for pack, item := range *accessList.List(context.Background()) {
				for right := range item {
					if !strings.Contains(right, "read") &&
						!strings.Contains(right, "view") {
						permission := &models.Permission{
							RoleName:    userRole.Name,
							PackageName: pack,
							LevelName:   right,
						}

						counter++
						_, err = adaptors.Permission.Add(context.Background(), permission)
						So(err, ShouldBeNil)
					}
				}
			}

			adminRole := &models.Role{
				Name:   "admin",
				Parent: userRole,
			}
			err = adaptors.Role.Add(context.Background(), adminRole)
			So(err, ShouldBeNil)

			// user
			userRole, err = adaptors.Role.GetByName(context.Background(), "user")
			So(err, ShouldBeNil)
			So(userRole.Name, ShouldEqual, "user")
			So(userRole.Parent, ShouldNotBeNil)
			So(userRole.Parent.Name, ShouldEqual, "demo")
			So(len(userRole.Children), ShouldEqual, 1)
			So(userRole.Children[0].Name, ShouldEqual, "admin")

			// demo
			demoRole, err = adaptors.Role.GetByName(context.Background(), "demo")
			So(err, ShouldBeNil)
			So(demoRole.Parent, ShouldBeNil)
			So(demoRole.Name, ShouldEqual, "demo")
			So(len(demoRole.Children), ShouldEqual, 1)
			So(demoRole.Children[0].Name, ShouldEqual, "user")

			// admin
			adminRole, err = adaptors.Role.GetByName(context.Background(), "admin")
			So(err, ShouldBeNil)
			So(adminRole.Parent, ShouldNotBeNil)
			So(adminRole.Parent.Name, ShouldEqual, "user")
			So(adminRole.Name, ShouldEqual, "admin")
			So(len(adminRole.Children), ShouldBeZeroValue)

			permissions, err := adaptors.Permission.GetAllPermissions(context.Background(), "user")
			So(err, ShouldBeNil)

			So(len(permissions), ShouldEqual, counter)

			err = adaptors.Role.Delete(context.Background(), "demo")
			So(err, ShouldNotBeNil)
			err = adaptors.Role.Delete(context.Background(), "user")
			So(err, ShouldNotBeNil)
			err = adaptors.Role.Delete(context.Background(), "admin")
			So(err, ShouldBeNil)
			err = adaptors.Role.Delete(context.Background(), "user")
			So(err, ShouldBeNil)
			err = adaptors.Role.Delete(context.Background(), "demo")
			So(err, ShouldBeNil)
		})
	})
}
