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
	"github.com/lib/pq"
	"time"
)

type Zigbee2mqttDevices struct {
	Db *gorm.DB
}

type Zigbee2mqttDevice struct {
	Id            string `gorm:"primary_key"`
	Zigbee2mqtt   *Zigbee2mqtt
	Zigbee2mqttId int64
	Name          string
	Type          string
	Model         string
	Description   string
	Manufacturer  string
	Status        string
	Functions     pq.StringArray `gorm:"type:varchar(100)[]"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (m *Zigbee2mqttDevice) TableName() string {
	return "zigbee2mqtt_devices"
}

func (z Zigbee2mqttDevices) Add(v *Zigbee2mqttDevice) (err error) {
	if err = z.Db.Create(&v).Error; err != nil {
		return
	}
	return
}

func (z Zigbee2mqttDevices) GetById(id string) (v *Zigbee2mqttDevice, err error) {
	v = &Zigbee2mqttDevice{Id: id}
	err = z.Db.First(&v).Error
	return
}

func (z Zigbee2mqttDevices) Update(m *Zigbee2mqttDevice) (err error) {
	err = z.Db.Model(&Zigbee2mqttDevice{Id: m.Id}).Updates(map[string]interface{}{
		"Name":         m.Name,
		"Type":         m.Type,
		"Model":        m.Model,
		"Description":  m.Description,
		"Manufacturer": m.Manufacturer,
		"Functions":    m.Functions,
		"Status":       m.Status,
	}).Error
	return
}

func (z Zigbee2mqttDevices) Delete(id string) (err error) {
	err = z.Db.Delete(&Zigbee2mqttDevice{Id: id}).Error
	return
}

func (z *Zigbee2mqttDevices) List(limit, offset int64) (list []*Zigbee2mqttDevice, total int64, err error) {

	if err = z.Db.Model(Zigbee2mqttDevice{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Zigbee2mqttDevice, 0)
	err = z.Db.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	return
}

func (z *Zigbee2mqttDevices) Search(query string, limit, offset int) (list []*Zigbee2mqttDevice, total int64, err error) {

	q := z.Db.Model(&Zigbee2mqttDevice{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Zigbee2mqttDevice, 0)
	err = q.Find(&list).Error

	return
}
