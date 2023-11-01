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
	"github.com/e154/smart-home/system/orm"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IAction ...
type IAction interface {
	Add(ctx context.Context, ver *m.Action) (id int64, err error)
	GetById(ctx context.Context, id int64) (metric *m.Action, err error)
	Update(ctx context.Context, ver *m.Action) error
	Delete(ctx context.Context, deviceId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Action, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Action, total int64, err error)
	fromDb(dbVer *db.Action) (ver *m.Action)
	toDb(ver *m.Action) (dbVer *db.Action)
}

// Action ...
type Action struct {
	IAction
	table *db.Actions
	db    *gorm.DB
	orm   *orm.Orm
}

// GetActionAdaptor ...
func GetActionAdaptor(d *gorm.DB, orm *orm.Orm) IAction {
	return &Action{
		table: &db.Actions{Db: d},
		db:    d,
		orm:   orm,
	}
}

// Add ...
func (n *Action) Add(ctx context.Context, ver *m.Action) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetById ...
func (n *Action) GetById(ctx context.Context, id int64) (metric *m.Action, err error) {
	var dbVer *db.Action
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Action) GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *string) (metric *m.Action, err error) {
	var dbVer *db.Action
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Action) Update(ctx context.Context, ver *m.Action) error {
	return n.table.Update(ctx, n.toDb(ver))
}

// Delete ...
func (n *Action) Delete(ctx context.Context, deviceId int64) (err error) {
	err = n.table.Delete(ctx, deviceId)
	return
}

// List ...
func (n *Action) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Action, total int64, err error) {
	var dbList []*db.Action
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Action, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Action) Search(ctx context.Context, query string, limit, offset int) (list []*m.Action, total int64, err error) {
	var dbList []*db.Action
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
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
		Description:      dbVer.Description,
		ScriptId:         dbVer.ScriptId,
		EntityId:         dbVer.EntityId,
		AreaId:           dbVer.AreaId,
		EntityActionName: dbVer.EntityActionName,
		CreatedAt:        dbVer.CreatedAt,
		UpdatedAt:        dbVer.UpdatedAt,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}
	// area
	if dbVer.Area != nil {
		scriptAdaptor := GetAreaAdaptor(n.db)
		ver.Area = scriptAdaptor.fromDb(dbVer.Area)
	}
	// entity
	if dbVer.Entity != nil {
		entityAdaptor := GetEntityAdaptor(n.db, n.orm)
		ver.Entity = entityAdaptor.fromDb(dbVer.Entity)
	}
	return
}

func (n *Action) toDb(ver *m.Action) (dbVer *db.Action) {
	dbVer = &db.Action{
		Id:               ver.Id,
		Name:             ver.Name,
		Description:      ver.Description,
		ScriptId:         ver.ScriptId,
		EntityId:         ver.EntityId,
		AreaId:           ver.AreaId,
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

	if ver.Area != nil {
		dbVer.AreaId = common.Int64(ver.Area.Id)
	}

	return
}
