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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common/apperr"
)

// Plugins ...
type Plugins struct {
	Db *gorm.DB
}

// Plugin ...
type Plugin struct {
	Name     string `gorm:"primary_key"`
	Version  string
	Enabled  bool
	System   bool
	Actor    bool
	Settings json.RawMessage `gorm:"type:jsonb;not null"`
}

// TableName ...
func (d Plugin) TableName() string {
	return "plugins"
}

// Add ...
func (n Plugins) Add(plugin *Plugin) (err error) {
	if err = n.Db.Create(&plugin).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginAdd, err.Error())
		return
	}
	return
}

// CreateOrUpdate ...
func (n Plugins) CreateOrUpdate(v *Plugin) (err error) {
	err = n.Db.Model(&Plugin{}).
		Set("gorm:insert_option",
			fmt.Sprintf("ON CONFLICT (name) DO UPDATE SET version = '%s', enabled = '%t', system = '%t', settings = '%s', actor = '%t'", v.Version, v.Enabled, v.System, v.Settings, v.Actor)).
		Create(&v).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrPluginUpdate, err.Error())
	}
	return
}

// Update ...
func (n Plugins) Update(m *Plugin) (err error) {
	if err = n.Db.Model(&Plugin{Name: m.Name}).Updates(m).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Plugins) Delete(name string) (err error) {
	if err = n.Db.Delete(&Plugin{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginDelete, err.Error())
	}
	return
}

// List ...
func (n Plugins) List(limit, offset int, orderBy, sort string) (list []*Plugin, total int64, err error) {

	if err = n.Db.Model(&Plugin{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginList, err.Error())
		return
	}

	list = make([]*Plugin, 0)
	q := n.Db.Model(&Plugin{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginList, err.Error())
	}

	return
}

// Search ...
func (n Plugins) Search(query string, limit, offset int) (list []*Plugin, total int64, err error) {

	q := n.Db.Model(&Plugin{}).
		Where("name LIKE ? and actor=true and enabled=true", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Plugin, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrPluginSearch, err.Error())
	}

	return
}

// GetByName ...
func (n Plugins) GetByName(name string) (plugin *Plugin, err error) {

	plugin = &Plugin{}
	err = n.Db.Model(plugin).
		Where("name = ?", name).
		First(&plugin).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrPluginNotFound, fmt.Sprintf("name \"%s\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrPluginGet, err.Error())
	}
	return
}
