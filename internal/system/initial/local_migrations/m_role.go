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

	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/internal/system/rbac/access_list"
	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/models"
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

func (r *MigrationRoles) Up(ctx context.Context) (err error) {

	if _, err = r.addAdmin(ctx); err != nil {
		return
	}
	var demo *models.Role
	if demo, err = r.addDemo(ctx); err != nil {
		return
	}
	_, err = r.addUser(ctx, demo)

	return
}

func (r *MigrationRoles) addAdmin(ctx context.Context) (adminRole *models.Role, err error) {

	if adminRole, err = r.adaptors.Role.GetByName(ctx, "admin"); err != nil {
		adminRole = &models.Role{
			Name: "admin",
		}
		err = r.adaptors.Role.Add(ctx, adminRole)
		So(err, ShouldBeNil)
	}
	if _, err = r.adaptors.User.GetByNickname(ctx, "admin"); err == nil {
		return
	}

	// add admin
	adminUser := &models.User{
		Nickname: "admin",
		RoleName: adminRole.Name,
		Email:    "admin@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = adminUser.SetPass("admin")
	So(err, ShouldBeNil)

	ok, _ := r.validation.Valid(adminUser)
	So(ok, ShouldEqual, true)

	adminUser.Id, err = r.adaptors.User.Add(ctx, adminUser)
	So(err, ShouldBeNil)

	return
}

func (r *MigrationRoles) addUser(ctx context.Context, demoRole *models.Role) (userRole *models.Role, err error) {

	if userRole, err = r.adaptors.Role.GetByName(ctx, "user"); err != nil {
		userRole = &models.Role{
			Name:   "user",
			Parent: demoRole,
		}
		err = r.adaptors.Role.Add(ctx, userRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List(ctx) {
			for levelName, right := range item {
				if right.Method == "put" || right.Method == "post" || right.Method == "delete" {
					permission := &models.Permission{
						RoleName:    userRole.Name,
						PackageName: pack,
						LevelName:   levelName,
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
	baseUser := &models.User{
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

func (r *MigrationRoles) addDemo(ctx context.Context) (demoRole *models.Role, err error) {

	if demoRole, err = r.adaptors.Role.GetByName(ctx, "demo"); err != nil {

		demoRole = &models.Role{
			Name: "demo",
		}
		err = r.adaptors.Role.Add(ctx, demoRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List(ctx) {
			for levelName, right := range item {
				if right.Method == "get" {
					permission := &models.Permission{
						RoleName:    demoRole.Name,
						PackageName: pack,
						LevelName:   levelName,
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
	demoUser := &models.User{
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
