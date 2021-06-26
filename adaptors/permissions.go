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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type IPermission interface {
	Add(permission *m.Permission) (id int64, err error)
	Delete(roleName, packageName string, levelName []string) (err error)
	GetAllPermissions(roleName string) (permissions []*m.Permission, err error)
	fromDb(dbPermission *db.Permission) (permission *m.Permission)
	toDb(permission *m.Permission) (dbPermission *db.Permission)
}

// Permission ...
type Permission struct {
	IPermission
	table *db.Permissions
	db    *gorm.DB
}

// GetPermissionAdaptor ...
func GetPermissionAdaptor(d *gorm.DB) IPermission {
	return &Permission{
		table: &db.Permissions{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Permission) Add(permission *m.Permission) (id int64, err error) {

	dbPermission := n.toDb(permission)
	if id, err = n.table.Add(dbPermission); err != nil {
		return
	}

	return
}

// Delete ...
func (n *Permission) Delete(roleName, packageName string, levelName []string) (err error) {

	err = n.table.Delete(roleName, packageName, levelName)

	return
}

// GetAllPermissions ...
func (n *Permission) GetAllPermissions(roleName string) (permissions []*m.Permission, err error) {

	var dbPermissions []*db.Permission
	if dbPermissions, err = n.table.GetAllPermissions(roleName); err != nil {
		return
	}

	for _, dbVer := range dbPermissions {
		ver := n.fromDb(dbVer)
		permissions = append(permissions, ver)
	}

	return
}

func (n *Permission) fromDb(dbPermission *db.Permission) (permission *m.Permission) {
	permission = &m.Permission{
		Id:          dbPermission.Id,
		RoleName:    dbPermission.RoleName,
		PackageName: dbPermission.PackageName,
		LevelName:   dbPermission.LevelName,
	}

	return
}

func (n *Permission) toDb(permission *m.Permission) (dbPermission *db.Permission) {
	dbPermission = &db.Permission{
		RoleName:    permission.RoleName,
		LevelName:   permission.LevelName,
		PackageName: permission.PackageName,
	}
	return
}
