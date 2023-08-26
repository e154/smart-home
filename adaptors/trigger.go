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
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// ITrigger ...
type ITrigger interface {
	Add(ver *m.Trigger) (id int64, err error)
	GetById(id int64) (metric *m.Trigger, err error)
	Update(ver *m.Trigger) error
	Delete(deviceId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Trigger, total int64, err error)
	Search(query string, limit, offset int) (list []*m.Trigger, total int64, err error)
	AddMultiple(items []*m.Trigger) (err error)
	fromDb(dbVer *db.Trigger) (ver *m.Trigger)
	toDb(ver *m.Trigger) (dbVer *db.Trigger)
}

// Trigger ...
type Trigger struct {
	ITrigger
	table *db.Triggers
	db    *gorm.DB
}

// GetTriggerAdaptor ...
func GetTriggerAdaptor(d *gorm.DB) ITrigger {
	return &Trigger{
		table: &db.Triggers{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Trigger) Add(ver *m.Trigger) (id int64, err error) {
	id, err = n.table.Add(n.toDb(ver))
	return
}

// GetById ...
func (n *Trigger) GetById(id int64) (metric *m.Trigger, err error) {
	var dbVer *db.Trigger
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Trigger) GetByIdWithData(id int64, from, to *time.Time, metricRange *string) (metric *m.Trigger, err error) {
	var dbVer *db.Trigger
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Trigger) Update(ver *m.Trigger) error {
	return n.table.Update(n.toDb(ver))
}

// Delete ...
func (n *Trigger) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// AddMultiple ...
func (n *Trigger) AddMultiple(items []*m.Trigger) (err error) {

	//TODO not work
	//insertRecords := make([]interface{}, 0, len(items))
	//for _, ver := range items {
	//	insertRecords = append(insertRecords, n.toDb(ver))
	//}
	//
	//err = gormbulk.BulkInsert(n.db, insertRecords, len(insertRecords))

	for _, ver := range items {
		if _, err = n.table.Add(n.toDb(ver)); err != nil {
			return
		}
	}

	return
}

// List ...
func (n *Trigger) List(limit, offset int64, orderBy, sort string) (list []*m.Trigger, total int64, err error) {
	var dbList []*db.Trigger
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Trigger, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Trigger) Search(query string, limit, offset int) (list []*m.Trigger, total int64, err error) {
	var dbList []*db.Trigger
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Trigger, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Trigger) fromDb(dbVer *db.Trigger) (ver *m.Trigger) {
	ver = &m.Trigger{
		Id:         dbVer.Id,
		Name:       dbVer.Name,
		EntityId:   dbVer.EntityId,
		ScriptId:   dbVer.ScriptId,
		PluginName: dbVer.PluginName,
		CreatedAt:  dbVer.CreatedAt,
		UpdatedAt:  dbVer.UpdatedAt,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}
	// entity
	if dbVer.Entity != nil {
		entityAdaptor := GetEntityAdaptor(n.db)
		ver.Entity = entityAdaptor.fromDb(dbVer.Entity)
	}

	// deserialize payload
	payload := m.TriggerPayload{}
	_ = json.Unmarshal([]byte(dbVer.Payload), &payload)
	ver.Payload = payload.Obj

	return
}

func (n *Trigger) toDb(ver *m.Trigger) (dbVer *db.Trigger) {
	dbVer = &db.Trigger{
		Id:         ver.Id,
		Name:       ver.Name,
		EntityId:   ver.EntityId,
		ScriptId:   ver.ScriptId,
		PluginName: ver.PluginName,
		CreatedAt:  ver.CreatedAt,
		UpdatedAt:  ver.UpdatedAt,
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)

	return
}
