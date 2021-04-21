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
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

// EntityActions ...
type EntityActions struct {
	Db *gorm.DB
}

// EntityAction ...
type EntityAction struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Icon        *string
	Entity      *Entity
	EntityId    common.EntityId
	Image       *Image
	ImageId     *int64
	Script      *Script
	ScriptId    *int64
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *EntityAction) TableName() string {
	return "entity_actions"
}

// Add ...
func (n EntityActions) Add(v *EntityAction) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n EntityActions) GetById(mapId int64) (v *EntityAction, err error) {
	v = &EntityAction{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

// Update ...
func (n EntityActions) Update(m *EntityAction) (err error) {
	err = n.Db.Model(&EntityAction{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"icon":        m.Icon,
		"entity_id":   m.EntityId,
		"image_id":    m.ImageId,
		"script_id":   m.ScriptId,
		"type":        m.Type,
	}).Error
	return
}

// DeleteByEntityId ...
func (n EntityActions) DeleteByEntityId(deviceId common.EntityId) (err error) {
	err = n.Db.Delete(&EntityAction{}, "entity_id = ?", deviceId).Error
	return
}

// List ...
func (n *EntityActions) List(limit, offset int64, orderBy, sort string) (list []*EntityAction, total int64, err error) {

	if err = n.Db.Model(EntityAction{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*EntityAction, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
