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

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// EntityStorages ...
type EntityStorages struct {
	Db *gorm.DB
}

// EntityStorage ...
type EntityStorage struct {
	Id         int64 `gorm:"primary_key"`
	EntityId   common.EntityId
	State      string
	Attributes json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time
}

// TableName ...
func (d *EntityStorage) TableName() string {
	return "entity_storage"
}

// Add ...
func (n *EntityStorages) Add(ctx context.Context, v EntityStorage) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetLastByEntityId ...
func (n *EntityStorages) GetLastByEntityId(ctx context.Context, entityId common.EntityId) (v EntityStorage, err error) {
	v = EntityStorage{}
	err = n.Db.WithContext(ctx).Model(&EntityStorage{}).
		Order("created_at desc").
		First(&v, "entity_id = ?", entityId).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrEntityNotFound, fmt.Sprintf("id \"%s\"", entityId))
			return
		}
		err = errors.Wrap(apperr.ErrEntityStorageGet, err.Error())
		return
	}
	return
}

// List ...
func (n *EntityStorages) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []EntityStorage, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(EntityStorage{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageList, err.Error())
		return
	}

	list = make([]EntityStorage, 0)
	q := n.Db.WithContext(ctx).Model(&EntityStorage{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageList, err.Error())
		return
	}
	return
}

// ListByEntityId ...
func (n *EntityStorages) ListByEntityId(ctx context.Context, limit, offset int, orderBy, sort string, entityIds []*common.EntityId, startDate, endDate *time.Time) (list []EntityStorage, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&EntityStorage{})


	if entityIds != nil && len(entityIds) > 0 {
		var ids = make([]string, 0, len(entityIds))
		for _, id := range entityIds {
			ids = append(ids, id.String())
		}
		q = q.Where("entity_id in (?)", ids)
	}

	if startDate != nil {
		q = q.Where("created_at >= ?", &startDate)
	}
	if endDate != nil {
		q = q.Where("created_at <= ?", &endDate)
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageList, err.Error())
		return
	}

	list = make([]EntityStorage, 0)
	q = q.
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageList, err.Error())
		return
	}
	return
}

// DeleteOldest ...
func (n *EntityStorages) DeleteOldest(ctx context.Context, days int) (err error) {
	storage := &EntityStorage{}
	if err = n.Db.WithContext(ctx).Last(&storage).Error; err != nil {
		err = errors.Wrap(apperr.ErrLogDelete, err.Error())
		return
	}
	err = n.Db.WithContext(ctx).Delete(&EntityStorage{},
		fmt.Sprintf(`created_at < CAST('%s' AS DATE) - interval '%d days'`,
			storage.CreatedAt.UTC().Format("2006-01-02 15:04:05"), days)).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageDelete, err.Error())
	}
	return
}
