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
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.VariableRepo = (*Variable)(nil)

// Variable ...
type Variable struct {
	table *db.Variables
	db    *gorm.DB
}

// GetVariableAdaptor ...
func GetVariableAdaptor(d *gorm.DB) *Variable {
	return &Variable{
		table: &db.Variables{&db.Common{Db: d}},
		db:    d,
	}
}

// CreateOrUpdate ...
func (n *Variable) CreateOrUpdate(ctx context.Context, ver models.Variable) (err error) {
	err = n.table.CreateOrUpdate(ctx, n.toDb(ver))
	return
}

// GetByName ...
func (n *Variable) GetByName(ctx context.Context, name string) (ver models.Variable, err error) {

	var dbVer db.Variable
	if dbVer, err = n.table.GetByName(ctx, name); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Delete ...
func (n *Variable) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// List ...
func (n *Variable) List(ctx context.Context, limit, offset int64, orderBy, sort string, system bool, name string) (list []models.Variable, total int64, err error) {
	var dbList []db.Variable
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, system, name); err != nil {
		return
	}

	list = make([]models.Variable, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// Search ...
func (n *Variable) Search(ctx context.Context, query string, limit, offset int) (list []models.Variable, total int64, err error) {
	var dbList []db.Variable
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]models.Variable, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// DeleteTags ...
func (n *Variable) DeleteTags(ctx context.Context, name string) (err error) {
	return n.table.DeleteTags(ctx, name)
}

func (n *Variable) fromDb(dbVer db.Variable) (ver models.Variable) {
	ver = models.Variable{
		Name:      dbVer.Name,
		Value:     dbVer.Value,
		System:    dbVer.System,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
		EntityId:  dbVer.EntityId,
	}
	// tags
	for _, tag := range dbVer.Tags {
		ver.Tags = append(ver.Tags, &models.Tag{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}
	return
}

func (n *Variable) toDb(ver models.Variable) (dbVer db.Variable) {
	dbVer = db.Variable{
		Name:     ver.Name,
		Value:    ver.Value,
		System:   ver.System,
		EntityId: ver.EntityId,
	}
	// tags
	if len(ver.Tags) > 0 {
		dbVer.Tags = make([]*db.Tag, 0, len(ver.Tags))
		for _, tag := range ver.Tags {
			dbVer.Tags = append(dbVer.Tags, &db.Tag{
				Id:   tag.Id,
				Name: tag.Name,
			})
		}
	} else {
		dbVer.Tags = make([]*db.Tag, 0)
	}
	return
}
