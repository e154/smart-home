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
	"time"
)

// ICondition ...
type ICondition interface {
	Add(ver *m.Condition) (id int64, err error)
	GetById(id int64) (metric *m.Condition, err error)
	Update(ver *m.Condition) error
	Delete(deviceId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.Condition, total int64, err error)
	Search(query string, limit, offset int) (list []*m.Condition, total int64, err error)
	fromDb(dbVer *db.Condition) (ver *m.Condition)
	toDb(ver *m.Condition) (dbVer *db.Condition)
}

// Condition ...
type Condition struct {
	ICondition
	table *db.Conditions
	db    *gorm.DB
}

// GetConditionAdaptor ...
func GetConditionAdaptor(d *gorm.DB) ICondition {
	return &Condition{
		table: &db.Conditions{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Condition) Add(ver *m.Condition) (id int64, err error) {
	id, err = n.table.Add(n.toDb(ver))
	return
}

// GetById ...
func (n *Condition) GetById(id int64) (metric *m.Condition, err error) {
	var dbVer *db.Condition
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Condition) GetByIdWithData(id int64, from, to *time.Time, metricRange *string) (metric *m.Condition, err error) {
	var dbVer *db.Condition
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Condition) Update(ver *m.Condition) error {
	return n.table.Update(n.toDb(ver))
}

// Delete ...
func (n *Condition) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

// List ...
func (n *Condition) List(limit, offset int64, orderBy, sort string) (list []*m.Condition, total int64, err error) {
	var dbList []*db.Condition
	if dbList, total, err = n.table.List(int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Condition, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Condition) Search(query string, limit, offset int) (list []*m.Condition, total int64, err error) {
	var dbList []*db.Condition
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Condition, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Condition) fromDb(dbVer *db.Condition) (ver *m.Condition) {
	ver = &m.Condition{
		Id:        dbVer.Id,
		Name:      dbVer.Name,
		ScriptId:  dbVer.ScriptId,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
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
		Id:        ver.Id,
		Name:      ver.Name,
		ScriptId:  ver.ScriptId,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}

	if ver.Script != nil {
		dbVer.ScriptId = ver.Script.Id
	}

	return
}
