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

// MapDevices ...
type MapDevices struct {
	Db *gorm.DB
}

// MapDevice ...
type MapDevice struct {
	Id        int64 `gorm:"primary_key"`
	Image     *Image
	ImageId   int64
	Device    *Device
	DeviceId  int64
	States    []*MapDeviceState
	Actions   []*MapDeviceAction
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName ...
func (d *MapDevice) TableName() string {
	return "map_devices"
}

// Add ...
func (n MapDevices) Add(v *MapDevice) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n MapDevices) GetById(mapId int64) (v *MapDevice, err error) {
	v = &MapDevice{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

// Delete ...
func (n MapDevices) Delete(id int64) (err error) {

	if err = n.Db.Delete(&MapDevice{Id: id}).Error; err != nil {
		return
	}

	if id != 0 {
		err = n.Db.Model(&MapElement{}).
			Where("prototype_id = ? and prototype_type = 'device'", id).
			Update("prototype_id", "").
			Error
	}

	return
}

// List ...
func (n *MapDevices) List(limit, offset int64, orderBy, sort string) (list []*MapDevice, total int64, err error) {

	if err = n.Db.Model(MapDevice{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapDevice, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
