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
	"fmt"
	"strings"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.EntityStateRepo = (*EntityState)(nil)

// EntityState ...
type EntityState struct {
	table *db.EntityStates
	db    *gorm.DB
}

// GetEntityStateAdaptor ...
func GetEntityStateAdaptor(d *gorm.DB) *EntityState {
	return &EntityState{
		table: &db.EntityStates{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *EntityState) Add(ctx context.Context, ver *m.EntityState) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(ctx, dbVer); err != nil {
		return
	}

	return
}

// DeleteByEntityId ...
func (n *EntityState) DeleteByEntityId(ctx context.Context, entityId pkgCommon.EntityId) (err error) {
	err = n.table.DeleteByEntityId(ctx, entityId)
	return
}

// AddMultiple ...
func (n *EntityState) AddMultiple(ctx context.Context, items []*m.EntityState) (err error) {

	if len(items) == 0 {
		return
	}

	insertRecords := make([]*db.EntityState, 0, len(items))

	for _, ver := range items {
		//if ver.ImageId == 0 {
		//	continue
		//}
		insertRecords = append(insertRecords, n.toDb(ver))
	}

	if err = n.table.AddMultiple(ctx, insertRecords); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityStateAdd)
	}

	return
}

func (n *EntityState) fromDb(dbVer *db.EntityState) (ver *m.EntityState) {
	ver = &m.EntityState{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Icon:        dbVer.Icon,
		//DeviceStateId: dbVer.DeviceStateId,
		EntityId:  dbVer.EntityId,
		ImageId:   dbVer.ImageId,
		Style:     dbVer.Style,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// state
	//if dbVer.DeviceState != nil {
	//	stateAdaptor := GetDeviceStateAdaptor(n.db)
	//	ver.DeviceState = stateAdaptor.fromDb(dbVer.DeviceState)
	//}

	return
}

func (n *EntityState) toDb(ver *m.EntityState) (dbVer *db.EntityState) {
	dbVer = &db.EntityState{
		Id:          ver.Id,
		Name:        strings.TrimSpace(ver.Name),
		Description: ver.Description,
		Icon:        ver.Icon,
		//DeviceStateId: ver.DeviceStateId,
		EntityId: ver.EntityId,
		ImageId:  ver.ImageId,
		Style:    ver.Style,
	}
	//if ver.DeviceState != nil && ver.DeviceState.Id != 0 {
	//	dbVer.DeviceStateId = ver.DeviceState.Id
	//}
	if ver.Image != nil && ver.Image.Id != 0 {
		dbVer.ImageId = pkgCommon.Int64(ver.Image.Id)
	}
	return
}
