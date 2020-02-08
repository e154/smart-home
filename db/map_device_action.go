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

type MapDeviceActions struct {
	Db *gorm.DB
}

type MapDeviceAction struct {
	Id             int64 `gorm:"primary_key"`
	DeviceAction   *DeviceAction
	DeviceActionId int64
	MapDevice      *MapDevice
	MapDeviceId    int64
	Image          *Image
	ImageId        int64
	Type           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d *MapDeviceAction) TableName() string {
	return "map_device_actions"
}

func (n MapDeviceActions) Add(v *MapDeviceAction) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapDeviceActions) GetById(mapId int64) (v *MapDeviceAction, err error) {
	v = &MapDeviceAction{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapDeviceActions) Update(m *MapDeviceAction) (err error) {
	err = n.Db.Model(&MapDeviceAction{Id: m.Id}).Updates(map[string]interface{}{
		"device_action_id": m.DeviceActionId,
		"map_device_id":    m.MapDeviceId,
		"image_id":         m.ImageId,
		"type":             m.Type,
	}).Error
	return
}

func (n MapDeviceActions) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapDeviceAction{Id: mapId}).Error
	return
}

func (n *MapDeviceActions) List(limit, offset int64, orderBy, sort string) (list []*MapDeviceAction, total int64, err error) {

	if err = n.Db.Model(MapDeviceAction{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDeviceAction, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
