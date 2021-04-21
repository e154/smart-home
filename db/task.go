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
	Conditions  []*Condition
	Triggers    []*Trigger
	Actions     []*Action
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
		return
	}
	id = task.Id
	return
}

// GetAllEnabled ...
func (n Tasks) GetAllEnabled() (list []*Task, err error) {
	list = make([]*Task, 0)
	err = n.Db.Where("enabled = ?", true).
		Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Area").
		Find(&list).Error
	return
}

// GetById ...
func (n Tasks) GetById(taskId int64) (task *Task, err error) {
	task = &Task{}
	err = n.Db.Model(task).
		Where("id = ?", taskId).
		Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Area").
		First(task).Error

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

	err = n.Db.Model(&Task{Id: m.Id}).Updates(q).Error
	return
}

// Delete ...
func (n Tasks) Delete(id int64) (err error) {
	err = n.Db.Delete(&Task{Id: id}).Error
	return
}

// List ...
func (n *Tasks) List(limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*Task, total int64, err error) {

	if err = n.Db.Model(Task{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Task, 0)
	q := n.Db.Model(&Task{})

	if onlyEnabled {
		q = q.Where("enabled = ?", true)
	}

	q = q.Preload("Triggers").
		Preload("Triggers.Script").
		Preload("Conditions").
		Preload("Conditions.Script").
		Preload("Actions").
		Preload("Actions.Script").
		Preload("Area").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

// Search ...
func (n *Tasks) Search(query string, limit, offset int) (list []*Task, total int64, err error) {

	q := n.Db.Model(&Task{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Task, 0)
	err = q.Find(&list).Error

	return
}
