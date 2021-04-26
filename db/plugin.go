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
)

// Plugins ...
type Plugins struct {
	Db *gorm.DB
}

// Plugin ...
type Plugin struct {
	Name    string `gorm:"primary_key"`
	Version string
	Enabled bool
	System  bool
}

// TableName ...
func (d Plugin) TableName() string {
	return "plugins"
}

// Add ...
func (n Plugins) Add(plugin Plugin) (err error) {
	if err = n.Db.Create(&plugin).Error; err != nil {
		return
	}
	return
}

// CreateOrUpdate ...
func (n Plugins) CreateOrUpdate(v Plugin) (err error) {
	err = n.Db.Model(&Plugin{}).
		Set("gorm:insert_option",
			fmt.Sprintf("ON CONFLICT (name) DO UPDATE SET version = '%s', enabled = '%t', system = '%t'", v.Version, v.Enabled, v.System)).
		Create(&v).Error
	if err != nil {
		log.Error(err.Error())
	}
	return
}

// Update ...
func (n Plugins) Update(m Plugin) (err error) {
	q := map[string]interface{}{
		"version":   m.Version,
		"installed": m.Enabled,
		"system":  m.System,
	}
	err = n.Db.Model(&Plugin{Name: m.Name}).Updates(q).Error
	return
}

// Delete ...
func (n Plugins) Delete(name string) (err error) {
	err = n.Db.Delete(&Plugin{Name: name}).Error
	return
}

// List ...
func (n Plugins) List(limit, offset int64, orderBy, sort string) (list []Plugin, total int64, err error) {

	if err = n.Db.Model(Plugin{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]Plugin, 0)
	q := n.Db.Model(&Plugin{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

// Search ...
func (n Plugins) Search(query string, limit, offset int) (list []Plugin, total int64, err error) {

	q := n.Db.Model(&Plugin{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]Plugin, 0)
	err = q.Find(&list).Error

	return
}

// GetByName ...
func (n Plugins) GetByName(name string) (plugin Plugin, err error) {

	plugin = Plugin{}
	err = n.Db.Model(plugin).
		Where("name = ?", name).
		First(&plugin).
		Error

	return
}
