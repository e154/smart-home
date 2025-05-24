// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

package adaptors

import (
	"context"
	"fmt"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.RoleRepo = (*Role)(nil)

// RoleName ...
type Role struct {
	table *db.Roles
	db    *gorm.DB
}

// GetRoleAdaptor ...
func GetRoleAdaptor(d *gorm.DB) *Role {
	return &Role{
		table: &db.Roles{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *Role) Add(ctx context.Context, role *models.Role) (err error) {

	dbRole := n.toDb(role)
	err = n.table.Add(ctx, dbRole)

	return
}

// GetByName ...
func (n *Role) GetByName(ctx context.Context, name string) (role *models.Role, err error) {

	var dbRole *db.Role
	if dbRole, err = n.table.GetByName(ctx, name); err != nil {
		return
	}

	role = n.fromDb(dbRole)

	err = n.GetAccessList(ctx, role)

	return
}

// Update ...
func (n *Role) Update(ctx context.Context, role *models.Role) (err error) {
	dbRole := n.toDb(role)
	err = n.table.Update(ctx, dbRole)
	return
}

// Delete ...
func (n *Role) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// List ...
func (n *Role) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*models.Role, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Role
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*models.Role, 0)
	for _, dbRole := range dbList {
		role := n.fromDb(dbRole)
		list = append(list, role)
	}

	return
}

// Search ...
func (n *Role) Search(ctx context.Context, query string, limit, offset int64) (list []*models.Role, total int64, err error) {
	var dbList []*db.Role
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*models.Role, 0)
	for _, dbRole := range dbList {
		role := n.fromDb(dbRole)
		list = append(list, role)
	}

	return
}

// GetAccessList ...
func (n *Role) GetAccessList(ctx context.Context, role *models.Role) (err error) {
	if role == nil {
		err = fmt.Errorf("role is nil: %w", apperr.ErrPermissionGet)
		return
	}
	role.AccessList = make(map[string][]string)
	permissionAdaptor := GetPermissionAdaptor(n.db)
	var permissions []*models.Permission
	if permissions, err = permissionAdaptor.GetAllPermissions(ctx, role.Name); err != nil {
		return
	}

	for _, perm := range permissions {
		if _, ok := role.AccessList[perm.PackageName]; !ok {
			role.AccessList[perm.PackageName] = []string{}
		}
		var exist bool
		for _, level := range role.AccessList[perm.PackageName] {
			if level == perm.LevelName {
				exist = true
			}
		}
		if !exist {
			role.AccessList[perm.PackageName] = append(role.AccessList[perm.PackageName], perm.LevelName)
		}
	}

	return
}

func (n *Role) fromDb(dbRole *db.Role) (role *models.Role) {
	role = &models.Role{
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt,
		UpdatedAt:   dbRole.UpdatedAt,
		Children:    []*models.Role{},
	}

	if dbRole.Role != nil {
		role.Parent = n.fromDb(dbRole.Role)
	}

	if len(dbRole.Children) > 0 {
		for _, dbCh := range dbRole.Children {
			ch := n.fromDb(dbCh)
			role.Children = append(role.Children, ch)
		}
	}

	return
}

func (n *Role) toDb(role *models.Role) (dbRole *db.Role) {
	dbRole = &db.Role{
		Name:        role.Name,
		Description: role.Description,
	}

	if role.Parent != nil {
		_ = dbRole.RoleName.Scan(role.Parent.Name)
	}
	return
}
