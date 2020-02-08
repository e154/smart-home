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
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/e154/smart-home/system/uuid"
	"encoding/json"
)

type Connections struct {
	Db *gorm.DB
}

type Connection struct {
	Uuid          uuid.UUID `gorm:"primary_key"`
	Name          string
	ElementFrom   uuid.UUID
	ElementTo     uuid.UUID
	PointFrom     int64
	PointTo       int64
	Flow          *Flow
	FlowId        int64
	Direction     string
	GraphSettings json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *Connection) TableName() string {
	return "connections"
}

func (n Connections) Add(connection *Connection) (id uuid.UUID, err error) {
	if err = n.Db.Create(&connection).Error; err != nil {
		return
	}
	id = connection.Uuid
	return
}

func (n Connections) GetById(id uuid.UUID) (connection *Connection, err error) {
	connection = &Connection{Uuid: id}
	err = n.Db.First(&connection).Error
	return
}

func (n Connections) Update(m *Connection) (err error) {
	err = n.Db.Model(&Connection{Uuid: m.Uuid}).Updates(map[string]interface{}{
		"name":           m.Name,
		"element_from":   m.ElementFrom,
		"element_to":     m.ElementTo,
		"point_from":     m.PointFrom,
		"point_to":       m.PointTo,
		"flow_id":        m.FlowId,
		"direction":      m.Direction,
		"graph_settings": m.GraphSettings,
	}).Error
	return
}

func (n Connections) Delete(ids []uuid.UUID) (err error) {
	err = n.Db.Delete(&Connection{}, "uuid in (?)", ids).Error
	return
}

func (n *Connections) List(limit, offset int64, orderBy, sort string) (list []*Connection, total int64, err error) {

	if err = n.Db.Model(Connection{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Connection, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
