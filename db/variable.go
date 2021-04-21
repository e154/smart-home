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

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// Variables ...
type Variables struct {
	Db *gorm.DB
}

// Variable ...
type Variable struct {
	Name      string `gorm:"primary_key"`
	Value     string
	Autoload  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName ...
func (d *Variable) TableName() string {
	return "variables"
}

// Add ...
func (n Variables) Add(variable Variable) (err error) {
	err = n.Db.Create(&variable).Error
	return
}

// CreateOrUpdate ...
func (n *Variables) CreateOrUpdate(v Variable) (err error) {
	err = n.Db.Model(&Variable{}).
		Set("gorm:insert_option",
			fmt.Sprintf("ON CONFLICT (name) DO UPDATE SET value = '%s', updated_at = '%s'", v.Value, time.Now().Format(time.RFC3339))).
		Create(&v).Error
	if err != nil {
		log.Error(err.Error())
	}
	return
}

// GetByName ...
func (n Variables) GetByName(name string) (variable Variable, err error) {
	variable = Variable{}
	err = n.Db.Model(&Variable{}).
		Where("name = ?", name).
		First(&variable).
		Error
	return
}

// GetAllEnabled ...
func (n Variables) GetAllEnabled() (list []Variable, err error) {
	list = make([]Variable, 0)
	err = n.Db.Where("autoload = ?", true).
		Find(&list).Error
	return
}

// Update ...
func (n Variables) Update(m Variable) (err error) {
	err = n.Db.Model(&Variable{Name: m.Name}).Updates(map[string]interface{}{
		"value":    m.Value,
		"autoload": m.Autoload,
	}).Error
	return
}

// Delete ...
func (n Variables) Delete(name string) (err error) {
	err = n.Db.Delete(&Variable{Name: name}).Error
	return
}

// List ...
func (n *Variables) List(limit, offset int64, orderBy, sort string) (list []Variable, total int64, err error) {

	if err = n.Db.Model(Variable{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]Variable, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
