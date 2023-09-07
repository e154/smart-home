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
	"fmt"
	"github.com/jinzhu/gorm"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ITask ...
type ITask interface {
	Add(ver *m.NewTask) (id int64, err error)
	Import(ver *m.Task) (err error)
	Update(ver *m.UpdateTask) (err error)
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

// Import ...
func (n *Task) Import(ver *m.Task) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			fmt.Println(err.Error())
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
	conditionAdaptor := GetConditionAdaptor(tx)
	for _, condition := range ver.Conditions {
		if condition.Id, err = conditionAdaptor.Add(condition); err != nil {
			return
		}
		if err = table.AppendCondition(ver.Id, conditionAdaptor.toDb(condition)); err != nil {
			return
		}
	}

	//triggers
	triggerAdaptor := GetTriggerAdaptor(tx)
	for _, trigger := range ver.Triggers {
		if trigger.Id, err = triggerAdaptor.Add(trigger); err != nil {
			return
		}
		if err = table.AppendTrigger(ver.Id, triggerAdaptor.toDb(trigger)); err != nil {
			return
		}
	}

	//actions
	actionAdaptor := GetActionAdaptor(tx)
	for _, action := range ver.Actions {
		if action.Id, err = actionAdaptor.Add(action); err != nil {
			return
		}
		if err = table.AppendAction(ver.Id, actionAdaptor.toDb(action)); err != nil {
			return
		}
	}

	return
}

// Add ...
func (n *Task) Add(ver *m.NewTask) (taskId int64, err error) {

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
	taskId, err = table.Add(&db.Task{
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		Condition:   ver.Condition,
		AreaId:      ver.AreaId,
	})
	if err != nil {
		return
	}

	//conditions
	conditionAdaptor := GetConditionAdaptor(tx)
	for _, id := range ver.ConditionIds {
		if err = table.AppendCondition(taskId, conditionAdaptor.toDb(&m.Condition{Id: id})); err != nil {
			return
		}
	}

	//triggers
	triggerAdaptor := GetTriggerAdaptor(tx)
	for _, id := range ver.TriggerIds {
		if err = table.AppendTrigger(taskId, triggerAdaptor.toDb(&m.Trigger{Id: id})); err != nil {
			return
		}
	}

	//actions
	actionAdaptor := GetActionAdaptor(tx)
	for _, id := range ver.ActionIds {
		if err = table.AppendAction(taskId, actionAdaptor.toDb(&m.Action{Id: id})); err != nil {
			return
		}
	}

	return
}

// Update ...
func (n *Task) Update(ver *m.UpdateTask) (err error) {

	var oldVer *m.Task
	if oldVer, err = n.GetById(ver.Id); err != nil {
		return
	}

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
	if err = table.Update(&db.Task{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		Condition:   ver.Condition,
		AreaId:      ver.AreaId,
	}); err != nil {
		return
	}

	//conditions
	for _, oldCondition := range oldVer.Conditions {
		var exist bool
		for _, id := range ver.ConditionIds {
			if id == oldCondition.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.DeleteCondition(oldVer.Id, oldCondition.Id); err != nil {
				return
			}
		}
	}

	conditionAdaptor := GetConditionAdaptor(tx)
	for _, id := range ver.ConditionIds {
		var exist bool
		for _, oldCondition := range oldVer.Conditions {
			if id == oldCondition.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.AppendCondition(ver.Id, conditionAdaptor.toDb(&m.Condition{Id: id})); err != nil {
				return
			}
		}
	}

	//triggers
	for _, oldTrigger := range oldVer.Triggers {
		var exist bool
		for _, id := range ver.TriggerIds {
			if id == oldTrigger.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.DeleteTrigger(oldVer.Id, oldTrigger.Id); err != nil {
				return
			}
		}
	}

	triggerAdaptor := GetTriggerAdaptor(tx)
	for _, id := range ver.TriggerIds {
		var exist bool
		for _, oldTrigger := range oldVer.Triggers {
			if id == oldTrigger.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.AppendTrigger(ver.Id, triggerAdaptor.toDb(&m.Trigger{Id: id})); err != nil {
				return
			}
		}
	}

	//actions
	for _, oldAction := range oldVer.Actions {
		var exist bool
		for _, id := range ver.ActionIds {
			if id == oldAction.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.DeleteAction(oldVer.Id, oldAction.Id); err != nil {
				return
			}
		}
	}

	actionAdaptor := GetActionAdaptor(tx)
	for _, id := range ver.ActionIds {
		var exist bool
		for _, oldAction := range oldVer.Actions {
			if id == oldAction.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.AppendAction(ver.Id, actionAdaptor.toDb(&m.Action{Id: id})); err != nil {
				return
			}
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
		AreaId:      dbVer.AreaId,
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
		AreaId:      ver.AreaId,
	}

	return
}
