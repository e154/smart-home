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
	Task       *Task
	TaskId     int64
	Entity     *Entity
	EntityId   *common.EntityId
	Script     *Script
	ScriptId   int64
	PluginName string
	Payload    string
}

// TableName ...
func (d *Trigger) TableName() string {
	return "triggers"
}

// DeleteByTaskId ...
func (n Triggers) DeleteByTaskId(id int64) (err error) {
	if err = n.Db.Delete(&Trigger{}, "task_id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerDelete, err.Error())
	}
	return
}

// AddMultiple ...
func (n *Triggers) AddMultiple(triggers []*Trigger) (err error) {
	if err = n.Db.Create(&triggers).Error; err != nil {
		err = errors.Wrap(apperr.ErrTriggerAdd, err.Error())
	}
	return
}
