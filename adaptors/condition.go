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

package adaptors

import (
	"context"
	"time"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ICondition ...
type ICondition interface {
	Add(ctx context.Context, ver *m.Condition) (id int64, err error)
	GetById(ctx context.Context, id int64) (metric *m.Condition, err error)
	Update(ctx context.Context, ver *m.Condition) error
	Delete(ctx context.Context, deviceId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, ids *[]uint64) (list []*m.Condition, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Condition, total int64, err error)
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
func (n *Condition) Add(ctx context.Context, ver *m.Condition) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetById ...
func (n *Condition) GetById(ctx context.Context, id int64) (metric *m.Condition, err error) {
	var dbVer *db.Condition
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Condition) GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *string) (metric *m.Condition, err error) {
	var dbVer *db.Condition
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Condition) Update(ctx context.Context, ver *m.Condition) error {
	return n.table.Update(ctx, n.toDb(ver))
}

// Delete ...
func (n *Condition) Delete(ctx context.Context, deviceId int64) (err error) {
	err = n.table.Delete(ctx, deviceId)
	return
}

// List ...
func (n *Condition) List(ctx context.Context, limit, offset int64, orderBy, sort string, ids *[]uint64) (list []*m.Condition, total int64, err error) {
	var dbList []*db.Condition
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, ids); err != nil {
		return
	}

	list = make([]*m.Condition, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Condition) Search(ctx context.Context, query string, limit, offset int) (list []*m.Condition, total int64, err error) {
	var dbList []*db.Condition
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
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
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		ScriptId:    dbVer.ScriptId,
		AreaId:      dbVer.AreaId,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}
	// area
	if dbVer.Area != nil {
		areaAdaptor := GetAreaAdaptor(n.db)
		ver.Area = areaAdaptor.fromDb(dbVer.Area)
	}
	return
}

func (n *Condition) toDb(ver *m.Condition) (dbVer *db.Condition) {
	dbVer = &db.Condition{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		ScriptId:    ver.ScriptId,
		AreaId:      ver.AreaId,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}

	if ver.Script != nil {
		dbVer.ScriptId = common.Int64(ver.Script.Id)
	}

	if ver.Area != nil {
		dbVer.AreaId = common.Int64(ver.Area.Id)
	}

	return
}
