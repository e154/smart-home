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

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
		err = errors.Wrap(apperr.ErrEntityActionAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n EntityActions) GetById(id int64) (v *EntityAction, err error) {
	v = &EntityAction{Id: id}
	if err = n.Db.First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrEntityActionNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrEntityActionGet, err.Error())
	}
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

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityActionUpdate, err.Error())
	}
	return
}

// DeleteByEntityId ...
func (n EntityActions) DeleteByEntityId(deviceId common.EntityId) (err error) {
	if err = n.Db.Delete(&EntityAction{}, "entity_id = ?", deviceId).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityActionDelete, err.Error())
	}
	return
}

// List ...
func (n *EntityActions) List(limit, offset int, orderBy, sort string) (list []*EntityAction, total int64, err error) {

	if err = n.Db.Model(EntityAction{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityActionList, err.Error())
		return
	}

	list = make([]*EntityAction, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityActionList, err.Error())
	}
	return
}

// AddMultiple ...
func (n *EntityActions) AddMultiple(actions []*EntityAction) (err error) {
	if err = n.Db.Create(&actions).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityActionAdd, err.Error())
	}
	return
}

