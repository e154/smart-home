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
	"encoding/json"
	"time"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ITrigger ...
type ITrigger interface {
	Add(ctx context.Context, ver *m.Trigger) (id int64, err error)
	GetById(ctx context.Context, id int64) (metric *m.Trigger, err error)
	Update(ctx context.Context, ver *m.Trigger) error
	Delete(ctx context.Context, deviceId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Trigger, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Trigger, total int64, err error)
	Enable(ctx context.Context, id int64) (err error)
	Disable(ctx context.Context, id int64) (err error)
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
func (n *Trigger) Add(ctx context.Context, ver *m.Trigger) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetById ...
func (n *Trigger) GetById(ctx context.Context, id int64) (metric *m.Trigger, err error) {
	var dbVer *db.Trigger
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// GetByIdWithData ...
func (n *Trigger) GetByIdWithData(ctx context.Context, id int64, from, to *time.Time, metricRange *string) (metric *m.Trigger, err error) {
	var dbVer *db.Trigger
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}
	metric = n.fromDb(dbVer)
	return
}

// Update ...
func (n *Trigger) Update(ctx context.Context, ver *m.Trigger) error {
	return n.table.Update(ctx, n.toDb(ver))
}

// Delete ...
func (n *Trigger) Delete(ctx context.Context, deviceId int64) (err error) {
	err = n.table.Delete(ctx, deviceId)
	return
}

// List ...
func (n *Trigger) List(ctx context.Context, limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Trigger, total int64, err error) {
	var dbList []*db.Trigger
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, onlyEnabled); err != nil {
		return
	}

	list = make([]*m.Trigger, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Search ...
func (n *Trigger) Search(ctx context.Context, query string, limit, offset int) (list []*m.Trigger, total int64, err error) {
	var dbList []*db.Trigger
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Trigger, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Enable ...
func (n *Trigger) Enable(ctx context.Context, id int64) (err error) {
	err = n.table.Enable(ctx, id)
	return
}

// Disable ...
func (n *Trigger) Disable(ctx context.Context, id int64) (err error) {
	err = n.table.Disable(ctx, id)
	return
}

func (n *Trigger) fromDb(dbVer *db.Trigger) (ver *m.Trigger) {
	ver = &m.Trigger{
		Id:         dbVer.Id,
		Name:       dbVer.Name,
		EntityId:   dbVer.EntityId,
		ScriptId:   dbVer.ScriptId,
		PluginName: dbVer.PluginName,
		Enabled:    dbVer.Enabled,
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
		Enabled:    ver.Enabled,
		CreatedAt:  ver.CreatedAt,
		UpdatedAt:  ver.UpdatedAt,
	}

	if ver.Script != nil {
		dbVer.ScriptId = common.Int64(ver.Script.Id)
	}

	if ver.Entity != nil {
		dbVer.EntityId = common.NewEntityId(ver.Entity.Id.String())
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)

	return
}
