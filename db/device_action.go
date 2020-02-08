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
	"fmt"
)

type DeviceActions struct {
	Db *gorm.DB
}

type DeviceAction struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Device      *Device
	DeviceId    int64 `gorm:"column:device_id"`
	Script      *Script
	ScriptId    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *DeviceAction) TableName() string {
	return "device_actions"
}

func (n DeviceActions) Add(action *DeviceAction) (id int64, err error) {
	if err = n.Db.Create(&action).Error; err != nil {
		return
	}
	id = action.Id
	return
}

func (n DeviceActions) GetById(actionId int64) (action *DeviceAction, err error) {
	action = &DeviceAction{Id: actionId}
	err = n.Db.Model(action).
		Preload("Script").
		First(&action).
		Error
	return
}

func (n DeviceActions) GetByDeviceId(deviceId int64) (actions []*DeviceAction, err error) {
	actions = make([]*DeviceAction, 0)
	err = n.Db.Model(&DeviceAction{}).
		Where("device_id = ?", deviceId).
		Preload("Script").
		Find(&actions).
		Error
	return
}

func (n DeviceActions) Update(m *DeviceAction) (err error) {
	err = n.Db.Model(&DeviceAction{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"script_id":   m.ScriptId,
		"device_id":   m.DeviceId,
	}).Error
	return
}

func (n DeviceActions) Delete(actionId int64) (err error) {
	err = n.Db.Delete(&DeviceAction{Id: actionId}).Error
	return
}

func (n *DeviceActions) List(limit, offset int64, orderBy, sort string) (list []*DeviceAction, total int64, err error) {

	if err = n.Db.Model(DeviceAction{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*DeviceAction, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *DeviceActions) Search(query string, limit, offset int) (list []*DeviceAction, total int64, err error) {

	q := n.Db.Model(&DeviceAction{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*DeviceAction, 0)
	err = q.Find(&list).Error

	return
}
