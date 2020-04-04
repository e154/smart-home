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

type AlexaIntents struct {
	Db *gorm.DB
}

type AlexaIntent struct {
	Name               string `gorm:"primary_key"`
	AlexaApplication   *AlexaApplication
	AlexaApplicationId string
	Flow               *Flow
	FlowId             int64
	Description        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (d *AlexaIntent) TableName() string {
	return "alexa_intents"
}

func (n AlexaIntents) Add(v *AlexaIntent) (id int64, err error) {
	err = n.Db.Create(&v).Error
	return
}

func (n AlexaIntents) Update(v *AlexaIntent) (err error) {
	err = n.Db.Model(v).Where("name = ? and alexa_application_id = ?", v.Name, v.AlexaApplicationId).Updates(&map[string]interface{}{
		"name":        v.Name,
		"description": v.Description,
		"flow_id":     v.FlowId,
	}).Error
	return
}

func (n AlexaIntents) Delete(v *AlexaIntent) (err error) {
	err = n.Db.Model(v).Delete("name = ? and alexa_application_id = ?", v.Name, v.AlexaApplicationId).Error
	return
}
