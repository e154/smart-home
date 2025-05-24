// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"gorm.io/gorm"
)

// Tasks ...
type Tasks struct {
	*Common
}

// Task ...
type Task struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Enabled     bool
	Condition   pkgCommon.ConditionType
	Conditions  []*Condition `gorm:"many2many:task_conditions;"`
	Actions     []*Action    `gorm:"many2many:task_actions;"`
	Triggers    []*Trigger   `gorm:"many2many:task_triggers;"`
	AreaId      *int64
	Area        *Area
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

// TableName ...
func (d *Task) TableName() string {
	return "tasks"
}

// Add ...
func (n Tasks) Add(ctx context.Context, task *Task) (id int64, err error) {
	err = n.DB(ctx).
		Omit("Triggers.*", "Conditions.*", "Actions.*").
		Create(&task).Error
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskAdd)
		return
	}
	id = task.Id
	return
}

// GetById ...
func (n Tasks) GetById(ctx context.Context, taskId int64) (task *Task, err error) {
	task = &Task{}
	err = n.DB(ctx).Model(task).
		Where("id = ?", taskId).
		Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Triggers.Entities").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Actions.Entity").
		Preload("Area").
		First(task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("id \"%d\"", taskId), apperr.ErrTaskNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskGet)
		return
	}

	return
}

// Update ...
func (n Tasks) Update(ctx context.Context, task *Task) (err error) {
	err = n.DB(ctx).
		Omit("Triggers.*", "Conditions.*", "Actions.*").
		Save(task).Error
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskUpdate)
	}
	return
}

// Delete ...
func (n Tasks) Delete(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Delete(&Task{Id: id}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskDelete)
		return
	}
	return
}

// Enable ...
func (n Tasks) Enable(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Model(&Task{Id: id}).Updates(map[string]interface{}{"enabled": true}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskUpdate)
		return
	}
	return
}

// Disable ...
func (n Tasks) Disable(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Model(&Task{Id: id}).Updates(map[string]interface{}{"enabled": false}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskUpdate)
		return
	}
	return
}

// List ...
func (n Tasks) List(ctx context.Context, limit, offset int, orderBy, sort string, onlyEnabled bool) (list []*Task, total int64, err error) {

	list = make([]*Task, 0)
	q := n.DB(ctx).Model(&Task{})

	if onlyEnabled {
		q = q.Where("enabled = ?", true)
	}

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskList)
		return
	}

	q = q.Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Triggers.Entities").
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskList)
	}
	return
}

// Search ...
func (n Tasks) Search(ctx context.Context, query string, limit, offset int) (list []*Task, total int64, err error) {

	q := n.DB(ctx).Model(&Task{}).
		Where("name ILIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskSearch)
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Task, 0)
	if err = q.Find(&list).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskSearch)
	}
	return
}

// DeleteTrigger ...
func (n Tasks) DeleteTrigger(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Model(&Task{Id: id}).Association("Triggers").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskDeleteTrigger)
	}
	return
}

// DeleteCondition ...
func (n Tasks) DeleteCondition(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Model(&Task{Id: id}).Association("Conditions").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskDeleteCondition)
	}
	return
}

// DeleteAction ...
func (n Tasks) DeleteAction(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Model(&Task{Id: id}).Association("Actions").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrTaskDeleteAction)
	}
	return
}
