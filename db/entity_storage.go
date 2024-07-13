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
	"strings"
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
	CreatedAt  time.Time       `gorm:"<-:create"`
}

// TableName ...
func (d *EntityStorage) TableName() string {
	return "entity_storage"
}

// Add ...
func (n *EntityStorages) Add(ctx context.Context, v *EntityStorage) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetLastByEntityId ...
func (n *EntityStorages) GetLastByEntityId(ctx context.Context, entityId common.EntityId) (v *EntityStorage, err error) {
	v = &EntityStorage{}
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

// GetLastThreeById ...
func (n *EntityStorages) GetLastThreeById(ctx context.Context, entityId common.EntityId, id int64) (list []*EntityStorage, err error) {
	list = make([]*EntityStorage, 2)
	err = n.Db.WithContext(ctx).Model(&EntityStorage{}).
		Order("id desc").
		Where("entity_id = ? and id <= ?", entityId, id).
		Limit(3).
		Find(&list).
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
func (n *EntityStorages) List(ctx context.Context, limit, offset int, orderBy, sort string, entityIds []common.EntityId, startDate, endDate *time.Time) (list []*EntityStorage, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&EntityStorage{})

	if len(entityIds) > 0 {
		q = q.Where("entity_id in (?)", entityIds)
	}

	if startDate != nil {
		q = q.Where("created_at > ?", startDate.UTC().Format(time.RFC3339))
	}
	if endDate != nil {
		q = q.Where("created_at <= ?", endDate.UTC().Format(time.RFC3339))
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityStorageList, err.Error())
		return
	}

	list = make([]*EntityStorage, 0)

	if sort != "" && orderBy != "" {
		switch sort {
		case "id", "entity_id", "state", "created_at":
			q = q.
				Order(fmt.Sprintf("%s %s", sort, orderBy))
		default:
			sort = strings.ReplaceAll(sort, "attributes_", "")
			q = q.
				Order(fmt.Sprintf("attributes->>'%s' %s", sort, orderBy))
		}
	}

	q = q.
		Limit(limit).
		Offset(offset)

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
