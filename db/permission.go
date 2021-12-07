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

package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Permissions ...
type Permissions struct {
	Db *gorm.DB
}

// Permission ...
type Permission struct {
	Id          int64 `gorm:"primary_key"`
	Role        *Role `gorm:"foreignkey:RoleName"`
	RoleName    string
	PackageName string
	LevelName   string
}

// TableName ...
func (m *Permission) TableName() string {
	return "permissions"
}

// Add ...
func (n Permissions) Add(permission *Permission) (id int64, err error) {
	if err = n.Db.Create(&permission).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = permission.Id
	return
}

// Delete ...
func (n Permissions) Delete(roleName, packageName string, levelName []string) (err error) {

	err = n.Db.
		Delete(&Permission{}, "role_name = ? and package_name = ? and level_name in (?)", roleName, packageName, levelName).
		Error
	if err != nil {
		err = errors.Wrap(err, "delete failed")
	}

	return
}

// GetAllPermissions ...
func (n Permissions) GetAllPermissions(name string) (permissions []*Permission, err error) {

	permissions = make([]*Permission, 0)
	err = n.Db.Raw(`
WITH RECURSIVE r AS (
    SELECT name, description, parent, created_at, updated_at, 1 AS level
    FROM roles
    WHERE name = ?

        UNION

        SELECT roles.name, roles.description, roles.parent, roles.created_at, roles.updated_at, r.level + 1 AS level
        FROM roles
               JOIN r
                 ON roles.name = r.parent
    )

SELECT DISTINCT p.*
FROM r
left join permissions p on p.role_name = r.name
where p notnull
order by p.id;
`, name).
		Scan(&permissions).
		Error
	if err != nil {
		err = errors.Wrap(err, "getAllPermissions failed")
	}
	return
}
