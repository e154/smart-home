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
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// DeviceStates ...
type DeviceStates struct {
	Db *gorm.DB
}

// DeviceState ...
type DeviceState struct {
	Id          int64 `gorm:"primary_key"`
	Device      *Device
	DeviceId    int64 `gorm:"column:device_id"`
	Description string
	SystemName  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (m *DeviceState) TableName() string {
	return "device_states"
}

// Add ...
func (n DeviceStates) Add(state *DeviceState) (id int64, err error) {
	if err = n.Db.Create(&state).Error; err != nil {
		return
	}
	id = state.Id
	return
}

// GetById ...
func (n DeviceStates) GetById(stateId int64) (state *DeviceState, err error) {
	state = &DeviceState{Id: stateId}
	err = n.Db.First(&state).Error
	return
}

// Update ...
func (n DeviceStates) Update(m *DeviceState) (err error) {
	err = n.Db.Model(&DeviceState{Id: m.Id}).Updates(map[string]interface{}{
		"system_name": m.SystemName,
		"description": m.Description,
	}).Error
	return
}

// Delete ...
func (n DeviceStates) Delete(stateId int64) (err error) {
	err = n.Db.Delete(&DeviceState{Id: stateId}).Error
	return
}

// List ...
func (n *DeviceStates) List(limit, offset int64, orderBy, sort string) (list []*DeviceState, total int64, err error) {

	if err = n.Db.Model(DeviceState{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*DeviceState, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

// GetByDeviceId ...
func (n DeviceStates) GetByDeviceId(deviceId int64) (actions []*DeviceState, err error) {
	actions = make([]*DeviceState, 0)
	err = n.Db.Model(&DeviceState{}).
		Where("device_id = ?", deviceId).
		Find(&actions).
		Error
	return
}
