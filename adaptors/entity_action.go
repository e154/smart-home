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
	"context"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// IEntityAction ...
type IEntityAction interface {
	Add(ctx context.Context, ver *m.EntityAction) (id int64, err error)
	DeleteByEntityId(ctx context.Context, id common.EntityId) (err error)
	AddMultiple(ctx context.Context, items []*m.EntityAction) (err error)
	fromDb(dbVer *db.EntityAction) (ver *m.EntityAction)
	toDb(ver *m.EntityAction) (dbVer *db.EntityAction)
}

// EntityAction ...
type EntityAction struct {
	IEntityAction
	table *db.EntityActions
	db    *gorm.DB
}

// GetEntityActionAdaptor ...
func GetEntityActionAdaptor(d *gorm.DB) IEntityAction {
	return &EntityAction{
		table: &db.EntityActions{Db: d},
		db:    d,
	}
}

// Add ...
func (n *EntityAction) Add(ctx context.Context, ver *m.EntityAction) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(ctx, dbVer); err != nil {
		return
	}

	return
}

// DeleteByEntityId ...
func (n *EntityAction) DeleteByEntityId(ctx context.Context, id common.EntityId) (err error) {
	err = n.table.DeleteByEntityId(ctx, id)
	return
}

// AddMultiple ...
func (n *EntityAction) AddMultiple(ctx context.Context, items []*m.EntityAction) (err error) {

	if len(items) == 0 {
		return
	}

	insertRecords := make([]*db.EntityAction, 0, len(items))

	for _, ver := range items {
		//if ver.ImageId == 0 {
		//	continue
		//}
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	if err = n.table.AddMultiple(ctx, insertRecords); err != nil {
		err = errors.Wrap(apperr.ErrEntityActionAdd, err.Error())
	}

	return
}

func (n *EntityAction) fromDb(dbVer *db.EntityAction) (ver *m.EntityAction) {
	ver = &m.EntityAction{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Icon:        dbVer.Icon,
		EntityId:    dbVer.EntityId,
		ImageId:     dbVer.ImageId,
		ScriptId:    dbVer.ScriptId,
		Type:        dbVer.Type,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// script
	if dbVer.Script != nil {
		scriptAdaptor := GetScriptAdaptor(n.db)
		ver.Script, _ = scriptAdaptor.fromDb(dbVer.Script)
	}

	return
}

func (n *EntityAction) toDb(ver *m.EntityAction) (dbVer *db.EntityAction) {
	dbVer = &db.EntityAction{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Icon:        ver.Icon,
		EntityId:    ver.EntityId,
		ImageId:     ver.ImageId,
		ScriptId:    ver.ScriptId,
		Type:        ver.Type,
	}
	if ver.Image != nil && ver.Image.Id != 0 {
		dbVer.ImageId = common.Int64(ver.Image.Id)
	}
	if ver.Script != nil && ver.Script.Id != 0 {
		dbVer.ScriptId = common.Int64(ver.Script.Id)
	}
	return
}
