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
	"encoding/json"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

// Trigger ...
type Trigger struct {
	table *db.Triggers
	db    *gorm.DB
}

// GetTriggerAdaptor ...
func GetTriggerAdaptor(d *gorm.DB) *Trigger {
	return &Trigger{
		table: &db.Triggers{Db: d},
		db:    d,
	}
}

// DeleteByTaskId ...
func (n *Trigger) DeleteByTaskId(id int64) (err error) {
	err = n.table.DeleteByTaskId(id)
	return
}

// AddMultiple ...
func (n *Trigger) AddMultiple(items []*m.Trigger) (err error) {

	insertRecords := make([]interface{}, 0, len(items))

	for _, ver := range items {
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))
	return
}

func (n *Trigger) fromDb(dbVer *db.Trigger) (ver *m.Trigger) {
	ver = &m.Trigger{
		Id:         dbVer.Id,
		Name:       dbVer.Name,
		EntityId:   dbVer.EntityId,
		TaskId:     dbVer.TaskId,
		ScriptId:   dbVer.ScriptId,
		PluginName: dbVer.PluginName,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}

	// deserialize payload
	payload := m.TriggerPayload{}
	json.Unmarshal([]byte(dbVer.Payload), &payload)
	ver.Payload = payload.Obj

	return
}

func (n *Trigger) toDb(ver *m.Trigger) (dbVer *db.Trigger) {
	dbVer = &db.Trigger{
		Id:         ver.Id,
		Name:       ver.Name,
		EntityId:   ver.EntityId,
		TaskId:     ver.TaskId,
		ScriptId:   ver.ScriptId,
		PluginName: ver.PluginName,
	}

	if ver.Script != nil {
		dbVer.ScriptId = ver.Script.Id
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)

	return
}
