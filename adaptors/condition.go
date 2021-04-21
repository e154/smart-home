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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

// Condition ...
type Condition struct {
	table *db.Conditions
	db    *gorm.DB
}

// GetConditionAdaptor ...
func GetConditionAdaptor(d *gorm.DB) *Condition {
	return &Condition{
		table: &db.Conditions{Db: d},
		db:    d,
	}
}

// DeleteByTaskId ...
func (n *Condition) DeleteByTaskId(id int64) (err error) {
	err = n.table.DeleteByTaskId(id)
	return
}

// AddMultiple ...
func (n *Condition) AddMultiple(items []*m.Condition) (err error) {

	insertRecords := make([]interface{}, 0, len(items))

	for _, ver := range items {
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))
	return
}

func (n *Condition) fromDb(dbVer *db.Condition) (ver *m.Condition) {
	ver = &m.Condition{
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

func (n *Condition) toDb(ver *m.Condition) (dbVer *db.Condition) {
	dbVer = &db.Condition{
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
