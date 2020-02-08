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

package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Role struct {
	table *db.Roles
	db    *gorm.DB
}

func GetRoleAdaptor(d *gorm.DB) *Role {
	return &Role{
		table: &db.Roles{Db: d},
		db:    d,
	}
}

func (n *Role) Add(role *m.Role) (err error) {

	dbRole := n.toDb(role)
	err = n.table.Add(dbRole)

	return
}

func (n *Role) GetByName(name string) (role *m.Role, err error) {

	var dbRole *db.Role
	if dbRole, err = n.table.GetByName(name); err != nil {
		return
	}

	role = n.fromDb(dbRole)

	err = n.GetAccessList(role)

	return
}

func (n *Role) Update(role *m.Role) (err error) {
	dbRole := n.toDb(role)
	err = n.table.Update(dbRole)
	return
}

func (n *Role) Delete(name string) (err error) {
	err = n.table.Delete(name)
	return
}

func (n *Role) List(limit, offset int64, orderBy, sort string) (list []*m.Role, total int64, err error) {
	var dbList []*db.Role
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Role, 0)
	for _, dbRole := range dbList {
		role := n.fromDb(dbRole)
		list = append(list, role)
	}

	return
}

func (n *Role) Search(query string, limit, offset int) (list []*m.Role, total int64, err error) {
	var dbList []*db.Role
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Role, 0)
	for _, dbRole := range dbList {
		role := n.fromDb(dbRole)
		list = append(list, role)
	}

	return
}

func (n *Role) GetAccessList(role *m.Role) (err error) {

	role.AccessList = make(map[string][]string)
	permissionAdaptor := GetPermissionAdaptor(n.db)
	var permissions []*m.Permission
	if permissions, err = permissionAdaptor.GetAllPermissions(role.Name); err != nil {
		return
	}

	for _, perm := range permissions {
		if _, ok := role.AccessList[perm.PackageName]; !ok {
			role.AccessList[perm.PackageName] = []string{}
		}
		role.AccessList[perm.PackageName] = append(role.AccessList[perm.PackageName], perm.LevelName)
	}

	return
}

func (n *Role) fromDb(dbRole *db.Role) (role *m.Role) {
	role = &m.Role{
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt,
		UpdatedAt:   dbRole.UpdatedAt,
		Children:    []*m.Role{},
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

func (n *Role) toDb(role *m.Role) (dbRole *db.Role) {
	dbRole = &db.Role{
		Name:        role.Name,
		Description: role.Description,
	}

	if role.Parent !=nil {
		dbRole.RoleName.Scan(role.Parent.Name)
	}
	return
}
