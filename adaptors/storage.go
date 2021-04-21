// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// Storage ...
type Storage struct {
	table *db.Storages
}

// GetStorageAdaptor ...
func GetStorageAdaptor(d *gorm.DB) *Storage {
	return &Storage{
		table: &db.Storages{Db: d},
	}
}

// CreateOrUpdate ...
func (s *Storage) CreateOrUpdate(ver m.Storage) (err error) {
	err = s.table.CreateOrUpdate(s.toDb(ver))
	return
}

// Delete ...
func (s *Storage) Delete(name string) (err error) {
	err = s.table.Delete(name)
	return
}

// Search ...
func (s *Storage) Search(query string, limit, offset int) (list []m.Storage, total int64, err error) {
	var dbList []db.Storage
	if dbList, total, err = s.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]m.Storage, 0)
	for _, dbVer := range dbList {
		list = append(list, s.fromDb(dbVer))
	}

	return
}

// GetByName ...
func (s *Storage) GetByName(name string) (ver m.Storage, err error) {

	var dbVer db.Storage
	if dbVer, err = s.table.GetByName(name); err != nil {
		return
	}

	ver = s.fromDb(dbVer)

	return
}

func (s *Storage) fromDb(dbVer db.Storage) (ver m.Storage) {
	ver = m.Storage{
		Name:      dbVer.Name,
		Value:     dbVer.Value,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}
	return
}

func (s *Storage) toDb(ver m.Storage) (dbVer db.Storage) {
	dbVer = db.Storage{
		Name:      ver.Name,
		Value:     ver.Value,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
	return
}
