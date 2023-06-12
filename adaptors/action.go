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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IAction ...
type IAction interface {
	DeleteByTaskId(id int64) (err error)
	AddMultiple(items []*m.Action) (err error)
	fromDb(dbVer *db.Action) (ver *m.Action)
	toDb(ver *m.Action) (dbVer *db.Action)
}

// Action ...
type Action struct {
	IAction
	table *db.Actions
	db    *gorm.DB
}

// GetActionAdaptor ...
func GetActionAdaptor(d *gorm.DB) IAction {
	return &Action{
		table: &db.Actions{Db: d},
		db:    d,
	}
}

// DeleteByTaskId ...
func (n *Action) DeleteByTaskId(id int64) (err error) {
	err = n.table.DeleteByTaskId(id)
	return
}

// AddMultiple ...
func (n *Action) AddMultiple(items []*m.Action) (err error) {

	insertRecords := make([]*db.Action, 0, len(items))

	for _, ver := range items {
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = n.table.AddMultiple(insertRecords)
	return
}

func (n *Action) fromDb(dbVer *db.Action) (ver *m.Action) {
	ver = &m.Action{
		Id:       dbVer.Id,
		Name:     dbVer.Name,
		TaskId:   dbVer.TaskId,
		ScriptId: dbVer.ScriptId,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}
	return
}

func (n *Action) toDb(ver *m.Action) (dbVer *db.Action) {
	dbVer = &db.Action{
		Id:       ver.Id,
		Name:     ver.Name,
		TaskId:   ver.TaskId,
		ScriptId: ver.ScriptId,
	}

	if ver.Script != nil {
		dbVer.ScriptId = ver.Script.Id
	}

	return
}
