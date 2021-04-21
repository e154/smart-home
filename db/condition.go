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
	"github.com/jinzhu/gorm"
)

// Conditions ...
type Conditions struct {
	Db *gorm.DB
}

// Condition ...
type Condition struct {
	Id       int64 `gorm:"primary_key"`
	Name     string
	Task     *Task
	TaskId   int64
	Script   *Script
	ScriptId int64
}

// TableName ...
func (d *Condition) TableName() string {
	return "conditions"
}

// DeleteByTaskId ...
func (n Conditions) DeleteByTaskId(id int64) (err error) {
	err = n.Db.Delete(&Condition{}, "task_id = ?", id).Error
	return
}
