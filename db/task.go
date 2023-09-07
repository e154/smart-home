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

// Tasks ...
type Tasks struct {
	Db *gorm.DB
}

// Task ...
type Task struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Enabled     bool
	Condition   common.ConditionType
	Conditions  []*Condition `gorm:"many2many:task_conditions;"`
	Actions     []*Action    `gorm:"many2many:task_actions;"`
	Triggers    []*Trigger   `gorm:"many2many:task_triggers;"`
	AreaId      *int64
	Area        *Area
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Task) TableName() string {
	return "tasks"
}

// Add ...
func (n Tasks) Add(task *Task) (id int64, err error) {
	if err = n.Db.Create(&task).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskAdd, err.Error())
		return
	}
	id = task.Id
	return
}

// GetById ...
func (n Tasks) GetById(taskId int64) (task *Task, err error) {
	task = &Task{}
	err = n.Db.Model(task).
		Where("id = ?", taskId).
		Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Triggers.Entity").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Actions.Entity").
		Preload("Area").
		First(task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrTaskNotFound, fmt.Sprintf("id \"%d\"", taskId))
			return
		}
		err = errors.Wrap(apperr.ErrTaskGet, err.Error())
		return
	}

	return
}

// Update ...
func (n Tasks) Update(m *Task) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"condition":   m.Condition,
		"area_id":     m.AreaId,
		"enabled":     m.Enabled,
	}

	if err = n.Db.Model(&Task{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskUpdate, err.Error())
		return
	}
	return
}

// Delete ...
func (n Tasks) Delete(id int64) (err error) {
	if err = n.Db.Delete(&Task{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskDelete, err.Error())
		return
	}
	return
}

// Enable ...
func (n Tasks) Enable(id int64) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Updates(map[string]interface{}{"enabled": true}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskUpdate, err.Error())
		return
	}
	return
}

// Disable ...
func (n Tasks) Disable(id int64) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Updates(map[string]interface{}{"enabled": false}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskUpdate, err.Error())
		return
	}
	return
}

// List ...
func (n *Tasks) List(limit, offset int, orderBy, sort string, onlyEnabled bool) (list []*Task, total int64, err error) {

	if err = n.Db.Model(Task{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskList, err.Error())
		return
	}

	list = make([]*Task, 0)
	q := n.Db.Model(&Task{})

	if onlyEnabled {
		q = q.Where("enabled = ?", true)
	}

	q = q.Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Triggers.Entity").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Actions.Entity").
		Preload("Area").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskList, err.Error())
	}
	return
}

// Search ...
func (n *Tasks) Search(query string, limit, offset int) (list []*Task, total int64, err error) {

	q := n.Db.Model(&Task{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Task, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrTaskSearch, err.Error())
	}
	return
}

// AppendTrigger ...
func (n *Tasks) AppendTrigger(id int64, trigger *Trigger) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Triggers").Append(trigger); err != nil {
		err = errors.Wrap(apperr.ErrTaskAppendTrigger, err.Error())
	}
	return
}

// DeleteTrigger ...
func (n *Tasks) DeleteTrigger(id, triggerId int64) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Triggers").Delete(&Trigger{Id: triggerId}); err != nil {
		err = errors.Wrap(apperr.ErrTaskDeleteTrigger, err.Error())
	}
	return
}

// AppendCondition ...
func (n *Tasks) AppendCondition(id int64, condition *Condition) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Conditions").Append(condition); err != nil {
		err = errors.Wrap(apperr.ErrTaskAppendCondition, err.Error())
	}
	return
}

// DeleteCondition ...
func (n *Tasks) DeleteCondition(id, conditionId int64) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Conditions").Delete(&Condition{Id: conditionId}); err != nil {
		err = errors.Wrap(apperr.ErrTaskDeleteCondition, err.Error())
	}
	return
}

// AppendAction ...
func (n *Tasks) AppendAction(id int64, action *Action) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Actions").Append(action); err != nil {
		err = errors.Wrap(apperr.ErrTaskAppendAction, err.Error())
	}
	return
}

// DeleteAction ...
func (n *Tasks) DeleteAction(id, actionId int64) (err error) {
	if err = n.Db.Model(&Task{Id: id}).Association("Actions").Delete(&Action{Id: actionId}); err != nil {
		err = errors.Wrap(apperr.ErrTaskDeleteAction, err.Error())
	}
	return
}
