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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

type Devices struct {
	Db *gorm.DB
}

type Device struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Device      *Device `gorm:"foreignkey:DeviceId"`
	DeviceId    sql.NullInt64
	Node        *Node
	NodeId      sql.NullInt64
	Status      string
	Type        common.DeviceType
	Properties  json.RawMessage `gorm:"type:jsonb;not null"`
	States      []*DeviceState
	Actions     []*DeviceAction
	Devices     []*Device
	IsGroup     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Device) TableName() string {
	return "devices"
}

func (n Devices) Add(device *Device) (id int64, err error) {
	if err = n.Db.Create(&device).Error; err != nil {
		return
	}
	id = device.Id
	return
}

func (n Devices) GetAllEnabled() (list []*Device, err error) {
	list = make([]*Device, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	if err != nil {
		return
	}

	for _, device := range list {
		n.DependencyLoading(device)
	}

	return
}

func (n Devices) GetById(deviceId int64) (device *Device, err error) {
	device = &Device{Id: deviceId}
	if err = n.Db.First(&device).Error; err != nil {
		return
	}
	err = n.DependencyLoading(device)
	return
}

func (n Devices) GetByDeviceActionId(deviceActionId int64) (device *Device, err error) {
	device = &Device{}
	err = n.Db.Raw(`select d.*
from devices d
left join device_actions da on d.id = da.device_id
where da.id = ? and da notnull`, deviceActionId).Scan(device).Error
	if err != nil {
		return
	}
	if device.Id == 0 {
		err = errors.New("record not found")
		return
	}
	err = n.DependencyLoading(device)
	return
}

func (n Devices) Update(m *Device) (err error) {
	err = n.Db.Model(&Device{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"properties":  m.Properties,
		"device_id":   m.DeviceId,
		"node":        m.Node,
		"type":        m.Type,
	}).Error
	return
}

func (n Devices) Delete(deviceId int64) (err error) {
	err = n.Db.Delete(&Device{Id: deviceId}).Error
	return
}

func (n *Devices) List(limit, offset int64, orderBy, sort string) (list []*Device, total int64, err error) {

	if err = n.Db.Model(Device{}).Count(&total).Error; err != nil {
		return
	}

	q := n.Db.Model(&Device{}).
		Preload("Device").
		Preload("Node").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	list = make([]*Device, 0)
	err = q.
		Find(&list).
		Error

	if err != nil {
		return
	}

	for _, device := range list {
		n.DependencyLoading(device)
	}

	return
}

func (n *Devices) DependencyLoading(device *Device) (err error) {

	// actions
	device.Actions = make([]*DeviceAction, 0)
	n.Db.Model(&DeviceAction{}).
		Where("device_id = ?", device.Id).
		Preload("Script").
		Find(&device.Actions)

	// states
	device.States = make([]*DeviceState, 0)
	n.Db.Model(&DeviceState{}).
		Where("device_id = ?", device.Id).
		Find(&device.States)

	// node
	if device.NodeId.Valid {
		device.Node = &Node{Id: device.NodeId.Int64}
		n.Db.Model(device.Node).
			Find(&device.Node)
	}

	// parent device
	if device.DeviceId.Valid {
		device.Device = &Device{}

		n.Db.Model(device).
			Related(&device.Device)

		// node
		if device.Device.NodeId.Valid {
			device.Device.Node = &Node{Id: device.Device.NodeId.Int64}
			n.Db.Model(device.Device.Node).
				Find(&device.Device.Node)
		}

		// actions
		device.Device.Actions = make([]*DeviceAction, 0)
		n.Db.Model(&DeviceAction{}).
			Where("device_id = ?", device.Device.Id).
			Preload("Script").
			Find(&device.Device.Actions)

		// states
		device.Device.States = make([]*DeviceState, 0)
		n.Db.Model(&DeviceState{}).
			Where("device_id = ?", device.Device.Id).
			Find(&device.Device.States)
	}

	// slave devices
	device.Devices = make([]*Device, 0)
	err = n.Db.Model(device).
		Where("device_id = ?", device.Id).
		Find(&device.Devices).
		Error

	return
}

func (n *Devices) RemoveState(deviceId, stateId int64) (err error) {
	if deviceId == 0 || stateId == 0 {
		err = errors.New("bad deviceId or stateId")
		return
	}
	err = n.Db.Delete(&DeviceState{DeviceId: deviceId, Id: stateId}).Error
	return
}

func (n *Devices) RemoveAction(deviceId, actionId int64) (err error) {
	if deviceId == 0 || actionId == 0 {
		err = errors.New("bad deviceId or actionId")
		return
	}
	err = n.Db.Delete(&DeviceAction{DeviceId: deviceId, Id: actionId}).Error
	return
}

func (n *Devices) Search(query string, limit, offset int) (list []*Device, total int64, err error) {

	q := n.Db.Model(&Device{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Device, 0)
	err = q.Find(&list).Error

	return
}
