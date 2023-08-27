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
	"github.com/jinzhu/gorm"
	"time"
)

// IAction ...
type IAction interface {
	Add(ver *m.Action) (id int64, err error)
	GetById(id int64) (metric *m.Action, err error)
	Update(ver *m.Action) error
	Delete(deviceId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Action, total int64, err error)
	Search(query string, limit, offset int) (list []*m.Action, total int64, err error)
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

// Add ...
func (n *Action) Add(ver *m.Action) (id int64, err error) {
	id, err = n.table.Add(n.toDb(ver))
	return
}

// GetById ...
func (n *Action) GetById(id int64) (metric *m.Action, err error) {
	var dbVer *db.Action
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Action) GetByIdWithData(id int64, from, to *time.Time, metricRange *string) (metric *m.Action, err error) {
	var dbVer *db.Action
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Action) Update(ver *m.Action) error {
	return n.table.Update(n.toDb(ver))
}

// Delete ...
func (n *Action) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// AddMultiple ...
func (n *Action) AddMultiple(items []*m.Action) (err error) {

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
func (n *Action) List(limit, offset int64, orderBy, sort string) (list []*m.Action, total int64, err error) {
	var dbList []*db.Action
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Action, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Action) Search(query string, limit, offset int) (list []*m.Action, total int64, err error) {
	var dbList []*db.Action
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Action, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Action) fromDb(dbVer *db.Action) (ver *m.Action) {
	ver = &m.Action{
		Id:               dbVer.Id,
		Name:             dbVer.Name,
		ScriptId:         dbVer.ScriptId,
		EntityId:         dbVer.EntityId,
		EntityActionName: dbVer.EntityActionName,
		CreatedAt:        dbVer.CreatedAt,
		UpdatedAt:        dbVer.UpdatedAt,
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
	return
}

func (n *Action) toDb(ver *m.Action) (dbVer *db.Action) {
	dbVer = &db.Action{
		Id:               ver.Id,
		Name:             ver.Name,
		ScriptId:         ver.ScriptId,
		EntityId:         ver.EntityId,
		EntityActionName: ver.EntityActionName,
		CreatedAt:        ver.CreatedAt,
		UpdatedAt:        ver.UpdatedAt,
	}

	if ver.Entity != nil {
		dbVer.EntityId = common.NewEntityId(ver.Entity.Id.String())
	}

	if ver.Script != nil {
		dbVer.ScriptId = common.Int64(ver.Script.Id)
	}

	return
}
