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
	"context"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IAction ...
type IAction interface {
	Add(ctx context.Context, ver *m.Action) (id int64, err error)
	GetById(ctx context.Context,id int64) (metric *m.Action, err error)
	Update(ctx context.Context,ver *m.Action) error
	Delete(ctx context.Context,deviceId int64) (err error)
	List(ctx context.Context,limit, offset int64, orderBy, sort string) (list []*m.Action, total int64, err error)
	Search(ctx context.Context,query string, limit, offset int) (list []*m.Action, total int64, err error)
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
