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
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

// Messages ...
type Messages struct {
	Db *gorm.DB
}

// Message ...
type Message struct {
	Id        int64 `gorm:"primary_key"`
	Type      string
	Payload   json.RawMessage `gorm:"type:jsonb;not null"`
	Statuses  []*MessageDelivery
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName ...
func (d *Message) TableName() string {
	return "messages"
}

// Add ...
func (n Messages) Add(msg Message) (id int64, err error) {
	if err = n.Db.Create(&msg).Error; err != nil {
		return
	}
	id = msg.Id
	return
}
