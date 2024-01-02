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
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/orm"
)

// ITrigger ...
type ITrigger interface {
	Add(ctx context.Context, ver *m.NewTrigger) (id int64, err error)
	GetById(ctx context.Context, id int64) (metric *m.Trigger, err error)
	Update(ctx context.Context, ver *m.UpdateTrigger) error
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
	orm   *orm.Orm
}

// GetTriggerAdaptor ...
func GetTriggerAdaptor(d *gorm.DB, orm *orm.Orm) ITrigger {
	return &Trigger{
		table: &db.Triggers{Db: d},
		db:    d,
		orm:   orm,
	}
}

// Add ...
func (n *Trigger) Add(ctx context.Context, ver *m.NewTrigger) (id int64, err error) {
	dbVer := &db.Trigger{
		Name:        ver.Name,
		Description: ver.Description,
		ScriptId:    ver.ScriptId,
		AreaId:      ver.AreaId,
		PluginName:  ver.PluginName,
		Enabled:     ver.Enabled,
	}

	// entities
	for _, entityId := range ver.EntityIds {
		dbVer.Entities = append(dbVer.Entities, &db.Entity{
			Id: common.EntityId(entityId),
		})
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)
	id, err = n.table.Add(ctx, dbVer)
	return
}

// Update ...
func (n *Trigger) Update(ctx context.Context, ver *m.UpdateTrigger) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	dbVer := &db.Trigger{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		ScriptId:    ver.ScriptId,
		AreaId:      ver.AreaId,
		PluginName:  ver.PluginName,
		Enabled:     ver.Enabled,
	}

	// entities
	for _, entityId := range ver.EntityIds {
		dbVer.Entities = append(dbVer.Entities, &db.Entity{
			Id: common.EntityId(entityId),
		})
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)

	table := db.Triggers{Db: tx}
	if err = table.DeleteEntity(ctx, dbVer.Id); err != nil {
		return
	}

	err = table.Update(ctx, dbVer)

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
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Entities:    make([]*m.Entity, 0, len(dbVer.Entities)),
		ScriptId:    dbVer.ScriptId,
		AreaId:      dbVer.AreaId,
		PluginName:  dbVer.PluginName,
		Enabled:     dbVer.Enabled,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}
	// entities
	if dbVer.Entities != nil {
		entityAdaptor := GetEntityAdaptor(n.db, n.orm)
		for _, entity := range dbVer.Entities {
			ver.Entities = append(ver.Entities, entityAdaptor.fromDb(entity))
		}
	}
	// aea
	if dbVer.Area != nil {
		entityAdaptor := GetAreaAdaptor(n.db)
		ver.Area = entityAdaptor.fromDb(dbVer.Area)
	}

	// deserialize payload
	payload := m.TriggerPayload{}
	_ = json.Unmarshal([]byte(dbVer.Payload), &payload)
	ver.Payload = payload.Obj

	return
}

func (n *Trigger) toDb(ver *m.Trigger) (dbVer *db.Trigger) {
	dbVer = &db.Trigger{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		ScriptId:    ver.ScriptId,
		AreaId:      ver.AreaId,
		PluginName:  ver.PluginName,
		Enabled:     ver.Enabled,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}

	if ver.Script != nil {
		dbVer.ScriptId = common.Int64(ver.Script.Id)
	}

	if ver.Area != nil {
		dbVer.AreaId = common.Int64(ver.Area.Id)
	}

	// entities
	for _, entity := range dbVer.Entities {
		dbVer.Entities = append(dbVer.Entities, &db.Entity{
			Id: entity.Id,
		})
	}

	// serialize payload
	b, _ := json.Marshal(m.TriggerPayload{
		Obj: ver.Payload,
	})
	dbVer.Payload = string(b)

	return
}
