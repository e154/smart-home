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

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"
	"gorm.io/gorm"
)

var _ adaptors.TagRepo = (*Tag)(nil)

// Tag ...
type Tag struct {
	table *db.Tags
	db    *gorm.DB
}

// GetTagAdaptor ...
func GetTagAdaptor(d *gorm.DB) *Tag {
	return &Tag{
		table: &db.Tags{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *Tag) Add(ctx context.Context, tag *m.Tag) (id int64, err error) {
	dbTag := n.toDb(tag)
	id, err = n.table.Add(ctx, dbTag)
	return
}

// GetByName ...
func (n *Tag) GetByName(ctx context.Context, name string) (tag *m.Tag, err error) {

	var dbTag *db.Tag
	if dbTag, err = n.table.GetByName(ctx, name); err != nil {
		return
	}

	tag, _ = n.fromDb(dbTag)

	return
}

// GetById ...
func (n *Tag) GetById(ctx context.Context, id int64) (tag *m.Tag, err error) {

	var dbTag *db.Tag
	if dbTag, err = n.table.GetById(ctx, id); err != nil {
		return
	}

	tag, _ = n.fromDb(dbTag)

	return
}

// Update ...
func (n *Tag) Update(ctx context.Context, tag *m.Tag) (err error) {
	err = n.table.Update(ctx, n.toDb(tag))
	return
}

// List ...
func (n *Tag) List(ctx context.Context, limit, offset int64, orderBy, sort string, query *string, names *[]string) (list []*m.Tag, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Tag
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, query, names); err != nil {
		return
	}

	list = make([]*m.Tag, 0)
	for _, dbTag := range dbList {
		tag, _ := n.fromDb(dbTag)
		list = append(list, tag)
	}

	return
}

// Delete ...
func (n *Tag) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// Search ...
func (n *Tag) Search(ctx context.Context, query string, limit, offset int64) (list []*m.Tag, total int64, err error) {
	var dbList []*db.Tag
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Tag, 0)
	for _, dbTag := range dbList {
		dev, _ := n.fromDb(dbTag)
		list = append(list, dev)
	}

	return
}

func (n *Tag) fromDb(dbVer *db.Tag) (ver *m.Tag, err error) {
	ver = &m.Tag{
		Id:   dbVer.Id,
		Name: dbVer.Name,
	}
	return
}

func (n *Tag) toDb(ver *m.Tag) (dbVer *db.Tag) {
	dbVer = &db.Tag{
		Id:   ver.Id,
		Name: ver.Name,
	}
	return
}
