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
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AlexaIntents ...
type AlexaIntents struct {
	Db *gorm.DB
}

// AlexaIntent ...
type AlexaIntent struct {
	Name         string `gorm:"primary_key"`
	AlexaSkill   *AlexaSkill
	AlexaSkillId int64
	Script       *Script
	ScriptId     int64
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName ...
func (d *AlexaIntent) TableName() string {
	return "alexa_intents"
}

// Add ...
func (n AlexaIntents) Add(v *AlexaIntent) (err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaIntentAdd, err.Error())
	}
	return
}

// GetByName ...
func (n AlexaIntents) GetByName(name string) (intent *AlexaIntent, err error) {
	intent = &AlexaIntent{}
	if err = n.Db.Model(intent).Where("name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrAlexaIntentNotFound, fmt.Sprintf("name \"w%s\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrAlexaIntentGet, err.Error())
	}
	return
}

// Update ...
func (n AlexaIntents) Update(v *AlexaIntent) (err error) {
	err = n.Db.Model(v).Where("name = ? and alexa_skill_id = ?", v.Name, v.AlexaSkillId).Updates(map[string]interface{}{
		"name":        v.Name,
		"description": v.Description,
		"script_id":   v.ScriptId,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrAlexaIntentUpdate, err.Error())
	}
	return
}

// Delete ...
func (n AlexaIntents) Delete(v *AlexaIntent) (err error) {
	if err = n.Db.Delete(&AlexaIntent{}, "name = ? and alexa_skill_id = ?", v.Name, v.AlexaSkillId).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaIntentDelete, err.Error())
	}
	return
}
