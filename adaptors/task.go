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

package adaptors

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ITask ...
type ITask interface {
	Add(ver *m.Task) (err error)
	Update(ver *m.Task) (err error)
	Delete(id int64) (err error)
	GetById(id int64) (task *m.Task, err error)
	List(limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Task, total int64, err error)
	Enable(id int64) (err error)
	Disable(id int64) (err error)
	fromDb(dbVer *db.Task) (ver *m.Task)
	toDb(ver *m.Task) (dbVer *db.Task)
}

// Task ...
type Task struct {
	ITask
	table *db.Tasks
	db    *gorm.DB
}

// GetTaskAdaptor ...
func GetTaskAdaptor(d *gorm.DB) ITask {
	return &Task{
		table: &db.Tasks{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Task) Add(ver *m.Task) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.Tasks{Db: tx}
	if ver.Id, err = table.Add(n.toDb(ver)); err != nil {
		return
	}

	//conditions
	if len(ver.Conditions) > 0 {
		for i := range ver.Conditions {
			ver.Conditions[i].TaskId = ver.Id
		}
		conditionAction := GetConditionAdaptor(tx)
		if err = conditionAction.AddMultiple(ver.Conditions); err != nil {
			return
		}
	}

	//triggers
	if len(ver.Triggers) > 0 {
		for i := range ver.Triggers {
			ver.Triggers[i].TaskId = ver.Id
		}
		triggerAction := GetTriggerAdaptor(tx)
		if err = triggerAction.AddMultiple(ver.Triggers); err != nil {
			return
		}
	}

	//actions
	if len(ver.Actions) > 0 {
		for i := range ver.Actions {
			ver.Actions[i].TaskId = ver.Id
		}
		actionAction := GetActionAdaptor(tx)
		if err = actionAction.AddMultiple(ver.Actions); err != nil {
			return
		}
	}

	return
}

// Update ...
func (n *Task) Update(ver *m.Task) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.Tasks{Db: tx}
	if err = table.Update(n.toDb(ver)); err != nil {
		return
	}

	//conditions
	conditionAction := GetConditionAdaptor(tx)
	_ = conditionAction.DeleteByTaskId(ver.Id)
	if len(ver.Conditions) > 0 {
		for i := range ver.Conditions {
			ver.Conditions[i].TaskId = ver.Id
		}
		if err = conditionAction.AddMultiple(ver.Conditions); err != nil {
			return
		}
	}

	//triggers
	triggerAction := GetTriggerAdaptor(tx)
	_ = triggerAction.DeleteByTaskId(ver.Id)
	if len(ver.Triggers) > 0 {
		for i := range ver.Triggers {
			ver.Triggers[i].TaskId = ver.Id
		}
		if err = triggerAction.AddMultiple(ver.Triggers); err != nil {
			return
		}
	}

	//actions
	actionAction := GetActionAdaptor(tx)
	_ = actionAction.DeleteByTaskId(ver.Id)
	if len(ver.Actions) > 0 {
		for i := range ver.Actions {
			ver.Actions[i].TaskId = ver.Id
		}
		if err = actionAction.AddMultiple(ver.Actions); err != nil {
			return
		}
	}

	return
}

// Enable ...
func (n *Task) Enable(id int64) (err error) {
	err = n.table.Enable(id)
	return
}

// Disable ...
func (n *Task) Disable(id int64) (err error) {
	err = n.table.Disable(id)
	return
}

// GetById ...
func (n *Task) GetById(id int64) (task *m.Task, err error) {

	var dbVer *db.Task
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	task = n.fromDb(dbVer)

	return
}

// Delete ...
func (n *Task) Delete(id int64) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := &db.Tasks{Db: tx}
	if err = table.Delete(id); err != nil {
		return
	}

	return
}

// List ...
func (n *Task) List(limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Task, total int64, err error) {

	var dbList []*db.Task
	if dbList, total, err = n.table.List(int(limit), int(offset), orderBy, sort, onlyEnabled); err != nil {
		return
	}

	list = make([]*m.Task, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Task) fromDb(dbVer *db.Task) (ver *m.Task) {
	ver = &m.Task{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Enabled:     dbVer.Enabled,
		Condition:   dbVer.Condition,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// triggers
	triggerAdaptor := GetTriggerAdaptor(n.db)
	for _, dbVer := range dbVer.Triggers {
		tr := triggerAdaptor.fromDb(dbVer)
		ver.Triggers = append(ver.Triggers, tr)
	}

	// conditions
	conditionAdaptor := GetConditionAdaptor(n.db)
	for _, dbVer := range dbVer.Conditions {
		cond := conditionAdaptor.fromDb(dbVer)
		ver.Conditions = append(ver.Conditions, cond)
	}

	// actions
	actionAdaptor := GetActionAdaptor(n.db)
	for _, dbVer := range dbVer.Actions {
		act := actionAdaptor.fromDb(dbVer)
		ver.Actions = append(ver.Actions, act)
	}

	// area
	if dbVer.Area != nil {
		areaAdaptor := GetAreaAdaptor(n.db)
		ver.Area = areaAdaptor.fromDb(dbVer.Area)
	}

	return
}

func (n *Task) toDb(ver *m.Task) (dbVer *db.Task) {
	dbVer = &db.Task{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		Condition:   ver.Condition,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}

	// area
	if ver.Area != nil {
		dbVer.AreaId = common.Int64(ver.Area.Id)
	}

	return
}
