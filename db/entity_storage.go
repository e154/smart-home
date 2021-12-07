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

package db

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
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
func (n *EntityStorages) Add(v EntityStorage) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetLastByEntityId ...
func (n *EntityStorages) GetLastByEntityId(entityId common.EntityId) (v EntityStorage, err error) {
	v = EntityStorage{}
	if err = n.Db.Model(&EntityStorage{}).
		Order("created_at desc").
		First(&v, "entity_id = ?", entityId).
		Error; err != nil {
		err = errors.Wrap(err, "getLastByEntityId failed")
		return
	}
	return
}

// List ...
func (n *EntityStorages) List(limit, offset int64, orderBy, sort string) (list []EntityStorage, total int64, err error) {

	if err = n.Db.Model(EntityStorage{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]EntityStorage, 0)
	q := n.Db.Model(&EntityStorage{}).
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
		err = errors.Wrap(err, "list failed")
		return
	}
	return
}
