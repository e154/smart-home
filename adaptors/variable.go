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

package adaptors

import (
	"context"

	"gorm.io/gorm"
	
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IVariable ...
type IVariable interface {
	Add(ctx context.Context, ver m.Variable) (err error)
	CreateOrUpdate(ctx context.Context, ver m.Variable) (err error)
	GetByName(ctx context.Context, name string) (ver m.Variable, err error)
	Update(ctx context.Context, variable m.Variable) (err error)
	Delete(ctx context.Context, name string) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, system bool, name string) (list []m.Variable, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []m.Variable, total int64, err error)
	fromDb(dbVer db.Variable) (ver m.Variable)
	toDb(ver m.Variable) (dbVer db.Variable)
}

// Variable ...
type Variable struct {
	IVariable
	table *db.Variables
	db    *gorm.DB
}

// GetVariableAdaptor ...
func GetVariableAdaptor(d *gorm.DB) IVariable {
	return &Variable{
		table: &db.Variables{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Variable) Add(ctx context.Context, ver m.Variable) (err error) {
	// tags
	tagAdaptor := GetTagAdaptor(n.db)
	for _, tag := range ver.Tags {
		var foundedTag *m.Tag
		if foundedTag, err = tagAdaptor.GetByName(ctx, tag.Name); err == nil {
			tag.Id = foundedTag.Id
		} else {
			tag.Id = 0
			if tag.Id, err = tagAdaptor.Add(ctx, tag); err != nil {
				return
			}
		}
	}

	err = n.table.Add(ctx, n.toDb(ver))
	return
}

// CreateOrUpdate ...
func (n *Variable) CreateOrUpdate(ctx context.Context, ver m.Variable) (err error) {

	if err = n.table.DeleteTags(ctx, ver.Name); err != nil {
		return
	}

	// tags
	tagAdaptor := GetTagAdaptor(n.db)
	for _, tag := range ver.Tags {
		var foundedTag *m.Tag
		if foundedTag, err = tagAdaptor.GetByName(ctx, tag.Name); err == nil {
			tag.Id = foundedTag.Id
		} else {
			tag.Id = 0
			if tag.Id, err = tagAdaptor.Add(ctx, tag); err != nil {
				return
			}
		}
	}

	err = n.table.CreateOrUpdate(ctx, n.toDb(ver))
	return
}

// GetByName ...
func (n *Variable) GetByName(ctx context.Context, name string) (ver m.Variable, err error) {

	var dbVer db.Variable
	if dbVer, err = n.table.GetByName(ctx, name); err != nil {

		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Variable) Update(ctx context.Context, ver m.Variable) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.Variables{Db: tx}

	var oldVer db.Variable
	if oldVer, err = table.GetByName(ctx, ver.Name); err != nil {
		err = n.Add(ctx, ver)
		return
	}

	if err = table.DeleteTags(ctx, oldVer.Name); err != nil {
		return
	}

	// tags
	tagAdaptor := GetTagAdaptor(tx)
	for _, tag := range ver.Tags {
		var foundedTag *m.Tag
		if foundedTag, err = tagAdaptor.GetByName(ctx, tag.Name); err == nil {
			tag.Id = foundedTag.Id
		} else {
			tag.Id = 0
			if tag.Id, err = tagAdaptor.Add(ctx, tag); err != nil {
				return
			}
		}
	}

	err = table.Update(ctx, n.toDb(ver))
	return
}

// Delete ...
func (n *Variable) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// List ...
func (n *Variable) List(ctx context.Context, limit, offset int64, orderBy, sort string, system bool, name string) (list []m.Variable, total int64, err error) {
	var dbList []db.Variable
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, system, name); err != nil {
		return
	}

	list = make([]m.Variable, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// Search ...
func (n *Variable) Search(ctx context.Context, query string, limit, offset int) (list []m.Variable, total int64, err error) {
	var dbList []db.Variable
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]m.Variable, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

func (n *Variable) fromDb(dbVer db.Variable) (ver m.Variable) {
	ver = m.Variable{
		Name:      dbVer.Name,
		Value:     dbVer.Value,
		System:    dbVer.System,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
		EntityId:  dbVer.EntityId,
	}
	// tags
	for _, tag := range dbVer.Tags {
		ver.Tags = append(ver.Tags, &m.Tag{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}
	return
}

func (n *Variable) toDb(ver m.Variable) (dbVer db.Variable) {
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
