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

package api

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestRole(t *testing.T) {

	Convey("POST /role", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /role/{name}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /role/{name}/access_list", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("PUT /role/{name}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("PUT /role/{name}/access_list", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("DELETE /role/{name}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /roles", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /roles/search", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})
}

func AddRoles(adaptors *adaptors.Adaptors,
	accessList *access_list.AccessListService,
	ctx C) {

	// admin role
	adminRole := &m.Role{
		Name: "admin",
	}
	err := adaptors.Role.Add(adminRole)
	ctx.So(err, ShouldBeNil)

	// demo role
	demoRole := &m.Role{
		Name: "demo",
	}
	err = adaptors.Role.Add(demoRole)
	ctx.So(err, ShouldBeNil)

	for pack, item := range *accessList.List {
		for right := range item {
			if strings.Contains(right, "read") ||
				strings.Contains(right, "view") ||
				strings.Contains(right, "preview") {
				permission := &m.Permission{
					RoleName:    demoRole.Name,
					PackageName: pack,
					LevelName:   right,
				}

				_, err = adaptors.Permission.Add(permission)
				ctx.So(err, ShouldBeNil)
			}
		}
	}

	// user role
	userRole := &m.Role{
		Name:   "user",
		Parent: demoRole,
	}
	err = adaptors.Role.Add(userRole)
	ctx.So(err, ShouldBeNil)

	// add admin
	adminUser := &m.User{
		Nickname: "admin",
		RoleName: adminRole.Name,
		Email:    "admin@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = adminUser.SetPass("admin")
	ctx.So(err, ShouldBeNil)

	ok, _ := adminUser.Valid()
	ctx.So(ok, ShouldEqual, true)

	adminUser.Id, err = adaptors.User.Add(adminUser)
	ctx.So(err, ShouldBeNil)

	// add demo user
	demoUser := &m.User{
		Nickname: "demo",
		RoleName: demoRole.Name,
		Email:    "demo@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = demoUser.SetPass("demo")
	ctx.So(err, ShouldBeNil)

	ok, _ = demoUser.Valid()
	ctx.So(ok, ShouldEqual, true)

	demoUser.Id, err = adaptors.User.Add(demoUser)
	ctx.So(err, ShouldBeNil)

	// add base user
	baseUser := &m.User{
		Nickname: "user",
		RoleName: userRole.Name,
		Email:    "user@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = baseUser.SetPass("user")
	ctx.So(err, ShouldBeNil)

	ok, _ = baseUser.Valid()
	ctx.So(ok, ShouldEqual, true)

	baseUser.Id, err = adaptors.User.Add(baseUser)
	ctx.So(err, ShouldBeNil)
}
