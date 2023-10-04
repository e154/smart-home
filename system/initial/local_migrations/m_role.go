// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package local_migrations

import (
	"context"
	"strings"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/validation"
)

type MigrationRoles struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	validation *validation.Validate
}

func NewMigrationRoles(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate) *MigrationRoles {
	return &MigrationRoles{
		adaptors:   adaptors,
		accessList: accessList,
		validation: validation,
	}
}

func (r *MigrationRoles) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		r.adaptors = adaptors
	}

	if _, err = r.addAdmin(ctx); err != nil {
		return
	}
	var demo *m.Role
	if demo, err = r.addDemo(ctx); err != nil {
		return
	}
	_, err = r.addUser(ctx, demo)

	return
}

func (r *MigrationRoles) addAdmin(ctx context.Context) (adminRole *m.Role, err error) {

	if adminRole, err = r.adaptors.Role.GetByName(ctx, "admin"); err != nil {
		adminRole = &m.Role{
			Name: "admin",
		}
		err = r.adaptors.Role.Add(ctx, adminRole)
		So(err, ShouldBeNil)
	}
	if _, err = r.adaptors.User.GetByNickname(ctx, "admin"); err == nil {
		return
	}

	// add admin
	adminUser := &m.User{
		Nickname: "admin",
		RoleName: adminRole.Name,
		Email:    "admin@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = adminUser.SetPass("admin")
	So(err, ShouldBeNil)

	for pack, item := range *r.accessList.List(ctx) {
		for right := range item {
			permission := &m.Permission{
				RoleName:    adminUser.Nickname,
				PackageName: pack,
				LevelName:   right,
			}

			_, err = r.adaptors.Permission.Add(ctx, permission)
			So(err, ShouldBeNil)
		}
	}

	ok, _ := r.validation.Valid(adminUser)
	So(ok, ShouldEqual, true)

	adminUser.Id, err = r.adaptors.User.Add(ctx, adminUser)
	So(err, ShouldBeNil)

	return
}

func (r *MigrationRoles) addUser(ctx context.Context, demoRole *m.Role) (userRole *m.Role, err error) {

	if userRole, err = r.adaptors.Role.GetByName(ctx, "user"); err != nil {
		userRole = &m.Role{
			Name:   "user",
			Parent: demoRole,
		}
		err = r.adaptors.Role.Add(ctx, userRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List(ctx) {
			for right := range item {
				if strings.Contains(right, "create") ||
					strings.Contains(right, "update") ||
					strings.Contains(right, "search") ||
					strings.Contains(right, "signout") ||
					strings.Contains(right, "read_access_list") ||
					strings.Contains(right, "delete") {
					permission := &m.Permission{
						RoleName:    userRole.Name,
						PackageName: pack,
						LevelName:   right,
					}

					_, err = r.adaptors.Permission.Add(ctx, permission)
					So(err, ShouldBeNil)
				}
			}
		}
	}

	if _, err = r.adaptors.User.GetByNickname(ctx, "user"); err == nil {
		return
	}

	// add base user
	baseUser := &m.User{
		Nickname: "user",
		RoleName: userRole.Name,
		Email:    "user@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = baseUser.SetPass("user")
	So(err, ShouldBeNil)

	ok, _ := r.validation.Valid(baseUser)
	So(ok, ShouldEqual, true)

	baseUser.Id, err = r.adaptors.User.Add(ctx, baseUser)
	So(err, ShouldBeNil)

	return
}

func (r *MigrationRoles) addDemo(ctx context.Context) (demoRole *m.Role, err error) {

	if demoRole, err = r.adaptors.Role.GetByName(ctx, "demo"); err != nil {

		demoRole = &m.Role{
			Name: "demo",
		}
		err = r.adaptors.Role.Add(ctx, demoRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List(ctx) {
			for right := range item {
				if strings.Contains(right, "read") ||
					strings.Contains(right, "signout") ||
					strings.Contains(right, "read_access_list") ||
					strings.Contains(right, "search") {
					permission := &m.Permission{
						RoleName:    demoRole.Name,
						PackageName: pack,
						LevelName:   right,
					}

					_, err = r.adaptors.Permission.Add(ctx, permission)
					So(err, ShouldBeNil)
				}
			}
		}
	}

	if _, err = r.adaptors.User.GetByNickname(ctx, "demo"); err == nil {
		return
	}

	// add demo user
	demoUser := &m.User{
		Nickname: "demo",
		RoleName: demoRole.Name,
		Email:    "demo@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = demoUser.SetPass("demo")
	So(err, ShouldBeNil)

	ok, _ := r.validation.Valid(demoUser)
	So(ok, ShouldEqual, true)

	demoUser.Id, err = r.adaptors.User.Add(ctx, demoUser)
	So(err, ShouldBeNil)

	return
}
