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
	"fmt"
	"time"
)

type MapDeviceStates struct {
	Db *gorm.DB
}

type MapDeviceState struct {
	Id            int64 `gorm:"primary_key"`
	DeviceState   *DeviceState
	DeviceStateId int64
	MapDevice     *MapDevice
	MapDeviceId   int64
	Image         *Image
	ImageId       int64
	Style         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *MapDeviceState) TableName() string {
	return "map_device_states"
}

func (n MapDeviceStates) Add(v *MapDeviceState) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapDeviceStates) GetById(mapId int64) (v *MapDeviceState, err error) {
	v = &MapDeviceState{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapDeviceStates) Update(m *MapDeviceState) (err error) {
	err = n.Db.Model(&MapDeviceState{Id: m.Id}).Updates(map[string]interface{}{
		"device_state_id": m.DeviceStateId,
		"map_device_id":   m.MapDeviceId,
		"image_id":        m.ImageId,
		"style":           m.Style,
	}).Error
	return
}

func (n MapDeviceStates) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapDeviceState{Id: mapId}).Error
	return
}

func (n *MapDeviceStates) List(limit, offset int64, orderBy, sort string) (list []*MapDeviceState, total int64, err error) {

	if err = n.Db.Model(MapDeviceState{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDeviceState, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
