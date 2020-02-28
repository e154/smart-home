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

type FlowZigbee2mqttDevice struct {
	db    *gorm.DB
	table *db.FlowZigbee2mqttDevices
}

func GetFlowZigbee2mqttDeviceAdaptor(Db *gorm.DB) *FlowZigbee2mqttDevice {
	return &FlowZigbee2mqttDevice{
		db:    Db,
		table: db.NewFlowZigbee2mqttDevices(Db),
	}
}

func (f *FlowZigbee2mqttDevice) Add(sub *m.FlowZigbee2mqttDevice) (err error) {
	err = f.table.Add(f.toDb(sub))
	return
}

func (f *FlowZigbee2mqttDevice) Remove(flowId int64, ids []string) (err error) {
	err = f.table.Delete(flowId, ids)
	return
}

func (f *FlowZigbee2mqttDevice) fromDb(dbVer *db.FlowZigbee2mqttDevice) (ver *m.FlowZigbee2mqttDevice) {

	ver = &m.FlowZigbee2mqttDevice{
		Id:                  dbVer.Id,
		FlowId:              dbVer.FlowId,
		Zigbee2mqttDeviceId: dbVer.Zigbee2mqttDeviceId,
		CreatedAt:           dbVer.CreatedAt,
	}

	return
}

func (f *FlowZigbee2mqttDevice) toDb(ver *m.FlowZigbee2mqttDevice) (dbVer *db.FlowZigbee2mqttDevice) {

	dbVer = &db.FlowZigbee2mqttDevice{
		Id:                  ver.Id,
		FlowId:              ver.FlowId,
		Zigbee2mqttDeviceId: ver.Zigbee2mqttDeviceId,
		CreatedAt:           ver.CreatedAt,
	}

	return
}
