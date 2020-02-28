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

package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FlowZigbee2mqttDevices struct {
	db *gorm.DB
}

func NewFlowZigbee2mqttDevices(db *gorm.DB) *FlowZigbee2mqttDevices {
	return &FlowZigbee2mqttDevices{db: db}
}

type FlowZigbee2mqttDevice struct {
	Id                  int64 `gorm:"primary_key"`
	Flow                *Flow
	FlowId              int64
	Zigbee2mqttDevice   *Zigbee2mqttDevice
	Zigbee2mqttDeviceId string
	CreatedAt           time.Time
}

func (d *FlowZigbee2mqttDevice) TableName() string {
	return "flow_zigbee2mqtt_devices"
}

func (f *FlowZigbee2mqttDevices) Add(sub *FlowZigbee2mqttDevice) (err error) {
	err = f.db.Create(sub).Error
	return
}

func (f *FlowZigbee2mqttDevices) Delete(flowId int64, ids []string) (err error) {
	err = f.db.Delete(&FlowZigbee2mqttDevice{}, "zigbee2mqtt_device_id in (?) and flow_id = ?", ids, flowId).Error
	return
}
