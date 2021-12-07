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
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

// Storages ...
type Storages struct {
	Db *gorm.DB
}

// Storage ...
type Storage struct {
	Name      string          `gorm:"primary_key"`
	Value     json.RawMessage `gorm:"type:jsonb;not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

// TableName ...
func (s *Storage) TableName() string {
	return "storage"
}

// CreateOrUpdate ...
func (s *Storages) CreateOrUpdate(v Storage) (err error) {
	err = s.Db.Model(&Storage{}).
		Set("gorm:insert_option",
			fmt.Sprintf("ON CONFLICT (name) DO UPDATE SET value = '%s', updated_at = '%s'", v.Value, time.Now().Format(time.RFC3339))).
		Create(&v).Error
	if err != nil {
		err = errors.Wrap(err, "createOrUpdate failed")
	}
	return
}

// GetByName ...
func (s *Storages) GetByName(name string) (v Storage, err error) {
	if err = s.Db.Model(&Storage{}).Where("name = ?", name).First(&v).Error; err != nil {
		err = errors.Wrap(err, "getByName failed")
	}
	return
}

// Delete ...
func (s *Storages) Delete(name string) (err error) {
	if err = s.Db.Delete(&Storage{}, "name = ?", name).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
	}
	return
}

// Search ...
func (s *Storages) Search(query string, limit, offset int) (list []Storage, total int64, err error) {

	q := s.Db.Model(&Storage{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]Storage, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(err, "search failed")
	}
	return
}
