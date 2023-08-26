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
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
)

// Actions ...
type Actions struct {
	Db *gorm.DB
}

// Action ...
type Action struct {
	Id               int64 `gorm:"primary_key"`
	Name             string
	Script           *Script
	ScriptId         *int64
	Entity           *Entity
	EntityId         *common.EntityId
	EntityActionName *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName ...
func (*Action) TableName() string {
	return "actions"
}

// Add ...
func (t Actions) Add(action *Action) (id int64, err error) {
	if err = t.Db.Create(&action).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionAdd, err.Error())
		return
	}
	id = action.Id
	return
}

// GetById ...
func (t Actions) GetById(id int64) (action *Action, err error) {
	action = &Action{}
	err = t.Db.Model(action).
		Where("id = ?", id).
		Preload("Entity").
		Preload("Script").
		First(&action).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrActionNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrActionGet, err.Error())
	}
	return
}

// Update ...
func (t Actions) Update(m *Action) (err error) {
	if err = t.Db.Model(&Action{}).Where("id = ?", m.Id).Updates(m).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionUpdate, err.Error())
	}
	return
}

// Delete ...
func (t Actions) Delete(id int64) (err error) {
	if err = t.Db.Delete(&Action{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionDelete, err.Error())
	}
	return
}

// List ...
func (t *Actions) List(limit, offset int64, orderBy, sort string) (list []*Action, total int64, err error) {

	if err = t.Db.Model(Action{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionList, err.Error())
		return
	}

	list = make([]*Action, 0)
	q := t.Db.Model(&Action{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionList, err.Error())
	}
	return
}

// Search ...q
func (t *Actions) Search(query string, limit, offset int) (list []*Action, total int64, err error) {

	q := t.Db.Model(&Action{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrActionSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Action, 0)
	err = q.Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrActionSearch, err.Error())
	}
	return
}
