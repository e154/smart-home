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
	"github.com/e154/smart-home/common/apperr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Triggers ...
type Triggers struct {
	Db *gorm.DB
}

// Trigger ...
type Trigger struct {
	Id         int64 `gorm:"primary_key"`
	Name       string
	Entity     *Entity
	EntityId   *common.EntityId
	Script     *Script
	ScriptId   *int64
	PluginName string
	Payload    string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName ...
func (*Trigger) TableName() string {
	return "triggers"
}

// Add ...
func (t Triggers) Add(trigger *Trigger) (id int64, err error) {
	if err = t.Db.Create(&trigger).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerAdd, err.Error())
		return
	}
	id = trigger.Id
	return
}

// GetById ...
func (t Triggers) GetById(id int64) (trigger *Trigger, err error) {
	trigger = &Trigger{}
	err = t.Db.Model(trigger).
		Where("id = ?", id).
		Preload("Entity").
		Preload("Script").
		First(&trigger).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrTriggerNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrTriggerGet, err.Error())
	}

	return
}

// Update ...
func (t Triggers) Update(m *Trigger) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"plugin_name": m.PluginName,
		"payload":     m.Payload,
		"enabled":     m.Enabled,
		"script_id":   m.ScriptId,
		"entity_id":   m.EntityId,
	}
	if err = t.Db.Model(&Trigger{}).Where("id = ?", m.Id).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerUpdate, err.Error())
	}
	return
}

// Delete ...
func (t Triggers) Delete(id int64) (err error) {
	if err = t.Db.Delete(&Trigger{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerDelete, err.Error())
	}
	return
}

// List ...
func (t *Triggers) List(limit, offset int, orderBy, sort string, onlyEnabled bool) (list []*Trigger, total int64, err error) {

	if err = t.Db.Model(Trigger{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerList, err.Error())
		return
	}

	list = make([]*Trigger, 0)
	q := t.Db.Model(&Trigger{})

	if onlyEnabled {
		q = q.Where("enabled = ?", true)
	}

	q = q.Preload("Entity").
		Preload("Script").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerList, err.Error())
	}
	return
}

// Search ...q
func (t *Triggers) Search(query string, limit, offset int) (list []*Trigger, total int64, err error) {

	q := t.Db.Model(&Trigger{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Trigger, 0)
	err = q.Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrTriggerSearch, err.Error())
	}
	return
}

// Enable ...
func (t Triggers) Enable(id int64) (err error) {
	if err = t.Db.Model(&Trigger{Id: id}).Updates(map[string]interface{}{"enabled": true}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerUpdate, err.Error())
		return
	}
	return
}

// Disable ...
func (t Triggers) Disable(id int64) (err error) {
	if err = t.Db.Model(&Trigger{Id: id}).Updates(map[string]interface{}{"enabled": false}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerUpdate, err.Error())
		return
	}
	return
}
