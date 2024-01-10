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
	"github.com/e154/smart-home/system/orm"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationUpdatePermissions struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	orm        *orm.Orm
}

func NewMigrationUpdatePermissions(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	orm *orm.Orm) *MigrationUpdatePermissions {
	return &MigrationUpdatePermissions{
		adaptors:   adaptors,
		accessList: accessList,
		orm:        orm,
	}
}

func (r *MigrationUpdatePermissions) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		r.adaptors = adaptors
	}

	if _, err = r.truncatePermissions(ctx); err != nil {
		return
	}

	if _, err = r.updateDemo(ctx); err != nil {
		return
	}

	if _, err = r.updateUser(ctx); err != nil {
		return
	}

	return
}

func (r *MigrationUpdatePermissions) truncatePermissions(ctx context.Context) (adminRole *m.Role, err error) {
	_, err = r.orm.DB().Exec("truncate table permissions")
	return
}

func (r *MigrationUpdatePermissions) updateUser(ctx context.Context) (userRole *m.Role, err error) {

	if userRole, err = r.adaptors.Role.GetByName(ctx, "user"); err == nil {
		for pack, item := range *r.accessList.List(ctx) {
			for levelName, right := range item {
				if right.Method == "put" || right.Method == "post" || right.Method == "delete" {
					permission := &m.Permission{
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

	return
}

func (r *MigrationUpdatePermissions) updateDemo(ctx context.Context) (demoRole *m.Role, err error) {

	if demoRole, err = r.adaptors.Role.GetByName(ctx, "demo"); err == nil {
		for pack, item := range *r.accessList.List(ctx) {
			for levelName, right := range item {
				if right.Method == "get" {
					permission := &m.Permission{
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

	return
}
