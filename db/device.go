package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"errors"
	"encoding/json"
	"database/sql"
)

type Devices struct {
	Db *gorm.DB
}

type Device struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Device      *Device
	DeviceId    sql.NullInt64
	Node        *Node
	NodeId      sql.NullInt64
	Status      string
	Type        string
	Properties  json.RawMessage `gorm:"type:jsonb;not null"`
	States      []*DeviceState
	Actions     []*DeviceAction
	Devices     []*Device
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

	list = make([]*Device, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
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

	device.States = make([]*DeviceState, 0)
	device.Actions = make([]*DeviceAction, 0)
	device.Devices = make([]*Device, 0)

	n.Db.Model(device).
		Related(&device.States).
		Related(&device.Actions)

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
