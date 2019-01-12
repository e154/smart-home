package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

type DeviceStates struct {
	Db *gorm.DB
}

type DeviceState struct {
	Id          int64 `gorm:"primary_key"`
	Device      *Device
	DeviceId    int64 `gorm:"column:device_id"`
	Description string
	SystemName  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *DeviceState) TableName() string {
	return "device_states"
}

func (n DeviceStates) Add(state *DeviceState) (id int64, err error) {
	if err = n.Db.Create(&state).Error; err != nil {
		return
	}
	id = state.Id
	return
}

func (n DeviceStates) GetById(stateId int64) (state *DeviceState, err error) {
	state = &DeviceState{Id: stateId}
	err = n.Db.First(&state).Error
	return
}

func (n DeviceStates) Update(m *DeviceState) (err error) {
	err = n.Db.Model(&DeviceState{Id: m.Id}).Updates(map[string]interface{}{
		"system_name": m.SystemName,
		"description": m.Description,
	}).Error
	return
}

func (n DeviceStates) Delete(stateId int64) (err error) {
	err = n.Db.Delete(&DeviceState{Id: stateId}).Error
	return
}

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
