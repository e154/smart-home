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
	"github.com/e154/smart-home/pkg/common/encryptor"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.Zigbee2mqttRepo = (*Zigbee2mqtt)(nil)

// Zigbee2mqtt ...
type Zigbee2mqtt struct {
	table *db.Zigbee2mqtts
	db    *gorm.DB
}

// GetZigbee2mqttAdaptor ...
func GetZigbee2mqttAdaptor(d *gorm.DB) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		table: &db.Zigbee2mqtts{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *Zigbee2mqtt) Add(ctx context.Context, ver *models.Zigbee2mqtt) (id int64, err error) {

	id, err = n.table.Add(ctx, n.toDb(ver))

	return
}

// GetById ...
func (n *Zigbee2mqtt) GetById(ctx context.Context, id int64) (ver *models.Zigbee2mqtt, err error) {

	var dbVer *db.Zigbee2mqtt
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Zigbee2mqtt) Update(ctx context.Context, ver *models.Zigbee2mqtt) (err error) {
	err = n.table.Update(ctx, n.toDb(ver))
	return
}

// Delete ...
func (n *Zigbee2mqtt) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// List ...
func (n *Zigbee2mqtt) List(ctx context.Context, limit, offset int64) (list []*models.Zigbee2mqtt, total int64, err error) {
	var dbList []*db.Zigbee2mqtt
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*models.Zigbee2mqtt, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// GetByLogin ...
func (a *Zigbee2mqtt) GetByLogin(ctx context.Context, login string) (ver *models.Zigbee2mqtt, err error) {

	var dbVer *db.Zigbee2mqtt
	if dbVer, err = a.table.GetByLogin(ctx, login); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

func (n *Zigbee2mqtt) fromDb(dbVer *db.Zigbee2mqtt) (ver *models.Zigbee2mqtt) {
	ver = &models.Zigbee2mqtt{
		Id:                dbVer.Id,
		Login:             dbVer.Login,
		Name:              dbVer.Name,
		PermitJoin:        dbVer.PermitJoin,
		BaseTopic:         dbVer.BaseTopic,
		CreatedAt:         dbVer.CreatedAt,
		UpdatedAt:         dbVer.UpdatedAt,
		EncryptedPassword: dbVer.EncryptedPassword,
	}

	if len(dbVer.Devices) > 0 {
		zigbee2mqttDeviceAdaptor := GetZigbee2mqttDeviceAdaptor(n.db)
		for _, dbDev := range dbVer.Devices {
			dev := zigbee2mqttDeviceAdaptor.fromDb(dbDev)
			ver.Devices = append(ver.Devices, dev)
		}
	} else {
		ver.Devices = make([]*models.Zigbee2mqttDevice, 0)
	}

	return
}

func (n *Zigbee2mqtt) toDb(ver *models.Zigbee2mqtt) (dbVer *db.Zigbee2mqtt) {
	dbVer = &db.Zigbee2mqtt{
		Id:                ver.Id,
		Login:             ver.Login,
		Name:              ver.Name,
		PermitJoin:        ver.PermitJoin,
		BaseTopic:         ver.BaseTopic,
		EncryptedPassword: ver.EncryptedPassword,
	}
	if ver.Password != nil {
		if *ver.Password == "" {
			dbVer.EncryptedPassword = ""
		} else {
			dbVer.EncryptedPassword, _ = encryptor.HashPassword(*ver.Password)
		}
	}
	return
}
