// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Zigbee2mqtts ...
type Zigbee2mqtts struct {
	Db *gorm.DB
}

// Zigbee2mqtt ...
type Zigbee2mqtt struct {
	Id                int64 `gorm:"primary_key"`
	Name              string
	Login             string
	Devices           []*Zigbee2mqttDevice
	EncryptedPassword string
	PermitJoin        bool
	BaseTopic         string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// TableName ...
func (m *Zigbee2mqtt) TableName() string {
	return "zigbee2mqtt"
}

// Add ...
func (z Zigbee2mqtts) Add(ctx context.Context, v *Zigbee2mqtt) (id int64, err error) {
	if err = z.Db.WithContext(ctx).Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrZigbee2mqttAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (z Zigbee2mqtts) GetById(ctx context.Context, id int64) (v *Zigbee2mqtt, err error) {
	v = &Zigbee2mqtt{Id: id}
	err = z.Db.WithContext(ctx).First(&v).
		Preload("Devices").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrZigbee2mqttNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrZigbee2mqttGet, err.Error())
	}
	return
}

// Update ...
func (z Zigbee2mqtts) Update(ctx context.Context, m *Zigbee2mqtt) (err error) {
	q := map[string]interface{}{
		"Name":               m.Name,
		"Login":              m.Login,
		"PermitJoin":         m.PermitJoin,
		"BaseTopic":          m.BaseTopic,
		"encrypted_password": m.EncryptedPassword,
	}

	if err = z.Db.WithContext(ctx).Model(&Zigbee2mqtt{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrZigbee2mqttUpdate, err.Error())
	}
	return
}

// Delete ...
func (z Zigbee2mqtts) Delete(ctx context.Context, id int64) (err error) {
	if err = z.Db.WithContext(ctx).Delete(&Zigbee2mqtt{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrZigbee2mqttDelete, err.Error())
	}
	return
}

// List ...
func (z *Zigbee2mqtts) List(ctx context.Context, limit, offset int) (list []*Zigbee2mqtt, total int64, err error) {

	if err = z.Db.WithContext(ctx).Model(Zigbee2mqtt{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Zigbee2mqtt, 0)
	err = z.Db.WithContext(ctx).
		Limit(limit).
		Preload("Devices").
		Offset(offset).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrZigbee2mqttList, err.Error())
	}
	return
}

// GetByLogin ...
func (z *Zigbee2mqtts) GetByLogin(ctx context.Context, login string) (bridge *Zigbee2mqtt, err error) {

	bridge = &Zigbee2mqtt{}
	err = z.Db.WithContext(ctx).Model(bridge).
		Where("login = ?", login).
		First(&bridge).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrZigbee2mqttGet, err.Error())
	}
	return
}
