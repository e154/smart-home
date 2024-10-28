// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.EntityStorageRepo = (*EntityStorage)(nil)

// EntityStorage ...
type EntityStorage struct {
	table *db.EntityStorages
	db    *gorm.DB
}

// GetEntityStorageAdaptor ...
func GetEntityStorageAdaptor(d *gorm.DB) *EntityStorage {
	return &EntityStorage{
		table: &db.EntityStorages{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *EntityStorage) Add(ctx context.Context, ver *models.EntityStorage) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(ver))
	return
}

// GetLastByEntityId ...
func (n *EntityStorage) GetLastByEntityId(ctx context.Context, entityId common.EntityId) (ver *models.EntityStorage, err error) {
	var dbVer *db.EntityStorage
	if dbVer, err = n.table.GetLastByEntityId(ctx, entityId); err != nil {
		return
	}
	ver = n.fromDb(dbVer)
	return
}

// List ...
func (n *EntityStorage) List(ctx context.Context, limit, offset int64, orderBy, sort string, entityIds []common.EntityId, startDate, endDate *time.Time) (list []*models.EntityStorage, total int64, err error) {
	var dbList []*db.EntityStorage
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, entityIds, startDate, endDate); err != nil {
		return
	}

	list = make([]*models.EntityStorage, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// GetLastThreeById ...
func (n *EntityStorage) GetLastThreeById(ctx context.Context, entityId common.EntityId, id int64) (list []*models.EntityStorage, err error) {
	var dbList []*db.EntityStorage
	if dbList, err = n.table.GetLastThreeById(ctx, entityId, id); err != nil {
		return
	}

	list = make([]*models.EntityStorage, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}
	return
}

// DeleteOldest ...
func (n *EntityStorage) DeleteOldest(ctx context.Context, days int) (err error) {
	err = n.table.DeleteOldest(ctx, days)
	return
}

func (n *EntityStorage) fromDb(dbVer *db.EntityStorage) (ver *models.EntityStorage) {
	ver = &models.EntityStorage{
		Id:         dbVer.Id,
		EntityId:   dbVer.EntityId,
		State:      dbVer.State,
		Attributes: models.AttributeValue{},
		CreatedAt:  dbVer.CreatedAt,
	}

	if len(dbVer.Attributes) > 0 {
		_ = json.Unmarshal(dbVer.Attributes, &ver.Attributes)
	}

	return
}

func (n *EntityStorage) toDb(ver *models.EntityStorage) (dbVer *db.EntityStorage) {
	dbVer = &db.EntityStorage{
		Id:        ver.Id,
		EntityId:  ver.EntityId,
		State:     ver.State,
		CreatedAt: ver.CreatedAt,
	}

	if ver.Attributes != nil {
		b, _ := json.Marshal(ver.Attributes)
		_ = dbVer.Attributes.UnmarshalJSON(b)
	}

	return
}
