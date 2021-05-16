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
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// Roles ...
type Roles struct {
	Db *gorm.DB
}

// Role ...
type Role struct {
	Name        string `gorm:"primary_key"`
	Description string
	Role        *Role
	RoleName    sql.NullString `gorm:"column:parent"`
	Children    []*Role
	Permissions []*Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (m *Role) TableName() string {
	return "roles"
}

// Add ...
func (n Roles) Add(role *Role) (err error) {
	if err = n.Db.Create(&role).Error; err != nil {
		return
	}
	return
}

// GetByName ...
func (n Roles) GetByName(name string) (role *Role, err error) {

	role = &Role{Name: name}
	err = n.Db.First(&role).Error
	if err != nil {
		return
	}

	err = n.RelData(role)

	return
}

// Update ...
func (n Roles) Update(m *Role) (err error) {
	err = n.Db.Model(&Role{Name: m.Name}).Updates(map[string]interface{}{
		"description": m.Description,
		"parent":      m.RoleName,
	}).Error
	return
}

// Delete ...
func (n Roles) Delete(name string) (err error) {
	err = n.Db.Delete(&Role{Name: name}).Error
	return
}

// List ...
func (n *Roles) List(limit, offset int64, orderBy, sort string) (list []*Role, total int64, err error) {

	if err = n.Db.Model(Role{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Role, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		return
	}

	for _, role := range list {
		n.RelData(role)
	}

	return
}

// Search ...
func (n *Roles) Search(query string, limit, offset int) (list []*Role, total int64, err error) {

	fmt.Println(query)
	q := n.Db.Model(&Role{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Role, 0)
	err = q.Find(&list).Error

	return
}

// RelData ...
func (n *Roles) RelData(role *Role) (err error) {

	// get parent
	if role.RoleName.Valid {
		role.Role = &Role{}
		err = n.Db.Model(role).
			Where("name = ?", role.RoleName.String).
			Find(&role.Role).
			Error
	}

	// get children
	role.Children = make([]*Role, 0)
	err = n.Db.Model(role).
		Where("parent = ?", role.Name).
		Find(&role.Children).
		Error

	return
}
