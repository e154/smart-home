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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type IEntityStorage interface {
	Add(ver m.EntityStorage) (id int64, err error)
	GetLastByEntityId(entityId common.EntityId) (ver m.EntityStorage, err error)
	List(limit, offset int64, orderBy, sort string) (list []m.EntityStorage, total int64, err error)
	fromDb(dbVer db.EntityStorage) (ver m.EntityStorage)
	toDb(ver m.EntityStorage) (dbVer db.EntityStorage)
}

// EntityStorage ...
type EntityStorage struct {
	IEntityStorage
	table *db.EntityStorages
	db    *gorm.DB
}

// GetEntityStorageAdaptor ...
func GetEntityStorageAdaptor(d *gorm.DB) IEntityStorage {
	return &EntityStorage{
		table: &db.EntityStorages{Db: d},
		db:    d,
	}
}

// Add ...
func (n *EntityStorage) Add(ver m.EntityStorage) (id int64, err error) {
	id, err = n.table.Add(n.toDb(ver))
	return
}

// GetLastByEntityId ...
func (n *EntityStorage) GetLastByEntityId(entityId common.EntityId) (ver m.EntityStorage, err error) {
	var dbVer db.EntityStorage
	if dbVer, err = n.table.GetLastByEntityId(entityId); err != nil {
		return
	}
	ver = n.fromDb(dbVer)
	return
}

// List ...
func (n *EntityStorage) List(limit, offset int64, orderBy, sort string) (list []m.EntityStorage, total int64, err error) {
	var dbList []db.EntityStorage
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]m.EntityStorage, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

func (n *EntityStorage) fromDb(dbVer db.EntityStorage) (ver m.EntityStorage) {
	ver = m.EntityStorage{
		Id:         dbVer.Id,
		EntityId:   dbVer.EntityId,
		State:      dbVer.State,
		Attributes: m.EntityAttributeValue{},
		CreatedAt:  dbVer.CreatedAt,
	}

	if len(dbVer.Attributes) > 0 {
		json.Unmarshal(dbVer.Attributes, &ver.Attributes)
	}

	return
}

func (n *EntityStorage) toDb(ver m.EntityStorage) (dbVer db.EntityStorage) {
	dbVer = db.EntityStorage{
		Id:        ver.Id,
		EntityId:  ver.EntityId,
		State:     ver.State,
		CreatedAt: ver.CreatedAt,
	}

	if ver.Attributes != nil {
		b, _ := json.Marshal(ver.Attributes)
		dbVer.Attributes.UnmarshalJSON(b)
	}

	return
}
