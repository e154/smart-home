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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// IUserMeta ...
type IUserMeta interface {
	UpdateOrCreate(meta *m.UserMeta) (id int64, err error)
	fromDb(dbMeta *db.UserMeta) (meta *m.UserMeta)
	toDb(meta *m.UserMeta) (dbMeta *db.UserMeta)
}

// UserMeta ...
type UserMeta struct {
	IUserMeta
	table *db.UserMetas
	db    *gorm.DB
}

// GetUserMetaAdaptor ...
func GetUserMetaAdaptor(d *gorm.DB) *UserMeta {
	return &UserMeta{
		table: &db.UserMetas{Db: d},
		db:    d,
	}
}

// UpdateOrCreate ...
func (n *UserMeta) UpdateOrCreate(meta *m.UserMeta) (id int64, err error) {

	dbMeta := n.toDb(meta)
	if id, err = n.table.UpdateOrCreate(dbMeta); err != nil {
		return
	}

	return
}

func (n *UserMeta) fromDb(dbMeta *db.UserMeta) (meta *m.UserMeta) {
	meta = &m.UserMeta{
		Id:     dbMeta.Id,
		Key:    dbMeta.Key,
		UserId: dbMeta.UserId,
		Value:  dbMeta.Value,
	}
	return
}

func (n *UserMeta) toDb(meta *m.UserMeta) (dbMeta *db.UserMeta) {
	dbMeta = &db.UserMeta{
		Id:     meta.Id,
		Key:    meta.Key,
		UserId: meta.UserId,
		Value:  meta.Value,
	}
	return
}
