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

type Zigbee2mqtts struct {
	Db *gorm.DB
}

type Zigbee2mqtt struct {
	Id                int64 `gorm:"primary_key"`
	Name              string
	Login             string
	Devices           []*Zigbee2mqttDevice
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (m *Zigbee2mqtt) TableName() string {
	return "zigbee2mqtt"
}

func (z Zigbee2mqtts) Add(v *Zigbee2mqtt) (id int64, err error) {
	if err = z.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (z Zigbee2mqtts) GetById(id int64) (v *Zigbee2mqtt, err error) {
	v = &Zigbee2mqtt{Id: id}
	err = z.Db.First(&v).
		Preload("Devices").Error
	return
}

func (z Zigbee2mqtts) Update(m *Zigbee2mqtt) (err error) {
	q := map[string]interface{}{
		"Name":              m.Name,
		"Login":             m.Login,
		"EncryptedPassword": m.EncryptedPassword,
	}
	err = z.Db.Model(&Zigbee2mqtt{Id: m.Id}).Updates(q).Error
	return
}

func (z Zigbee2mqtts) Delete(id int64) (err error) {
	err = z.Db.Delete(&Zigbee2mqtt{Id: id}).Error
	return
}

func (z *Zigbee2mqtts) List(limit, offset int64) (list []*Zigbee2mqtt, total int64, err error) {

	if err = z.Db.Model(Zigbee2mqtt{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Zigbee2mqtt, 0)
	err = z.Db.
		Limit(limit).
		Preload("Devices").
		Offset(offset).
		Find(&list).
		Error

	return
}
