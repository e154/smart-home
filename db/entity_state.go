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
	"fmt"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// EntityStates ...
type EntityStates struct {
	Db *gorm.DB
}

// EntityState ...
type EntityState struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Icon        *string
	Entity      *Entity
	EntityId    common.EntityId
	Image       *Image
	ImageId     *int64
	Style       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *EntityState) TableName() string {
	return "entity_states"
}

// Add ...
func (n EntityStates) Add(v *EntityState) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n EntityStates) GetById(mapId int64) (v *EntityState, err error) {
	v = &EntityState{Id: mapId}
	if err = n.Db.First(&v).Error; err != nil {
		err = errors.Wrap(err, "getById failed")
		return
	}
	return
}

// Update ...
func (n EntityStates) Update(m *EntityState) (err error) {
	err = n.Db.Model(&EntityState{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"icon":        m.Icon,
		//"device_state_id": m.DeviceStateId,
		"entity_id": m.EntityId,
		"image_id":  m.ImageId,
		"style":     m.Style,
	}).Error
	return
}

// DeleteByEntityId ...
func (n EntityStates) DeleteByEntityId(entityId common.EntityId) (err error) {
	if err = n.Db.Delete(&EntityState{}, "entity_id = ?", entityId).Error; err != nil {
		err = errors.Wrap(err, "deleteByEntityId failed")
		return
	}
	return
}

// List ...
func (n *EntityStates) List(limit, offset int64, orderBy, sort string) (list []*EntityState, total int64, err error) {

	if err = n.Db.Model(EntityState{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*EntityState, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(err, "list failed")
		return
	}
	return
}
