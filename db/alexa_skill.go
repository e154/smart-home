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
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

// AlexaSkills ...
type AlexaSkills struct {
	Db *gorm.DB
}

// AlexaSkill ...
type AlexaSkill struct {
	Id                   int64 `gorm:"primary_key"`
	SkillId              string
	Description          string
	Intents              []*AlexaIntent `gorm:"foreignkey:AlexaSkillId"`
	Status               common.StatusType
	OnLaunchScript       *Script `gorm:"foreignkey:OnLaunchScriptId"`
	OnLaunchScriptId     *int64  `gorm:"column:on_launch"`
	OnSessionEndScript   *Script
	OnSessionEndScriptId *int64 `gorm:"column:on_session_end"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// TableName ...
func (d *AlexaSkill) TableName() string {
	return "alexa_skills"
}

// Add ...
func (n AlexaSkills) Add(v *AlexaSkill) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n AlexaSkills) GetById(id int64) (v *AlexaSkill, err error) {
	v = &AlexaSkill{Id: id}
	err = n.Db.Model(v).
		Preload("OnLaunchScript").
		Preload("OnSessionEndScript").
		Find(v).
		Error
	if err != nil {
		return
	}
	err = n.preload(v)

	return
}

// List ...
func (n *AlexaSkills) List(limit, offset int64, orderBy, sort string) (list []*AlexaSkill, total int64, err error) {

	if err = n.Db.Model(AlexaSkill{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*AlexaSkill, 0)
	q := n.Db.Model(&AlexaSkill{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.Find(&list).Error

	return
}

// ListEnabled ...
func (n *AlexaSkills) ListEnabled(limit, offset int64) (list []*AlexaSkill, err error) {

	list = make([]*AlexaSkill, 0)
	err = n.Db.Model(&AlexaSkill{}).
		Where("status = 'enabled'").
		Limit(limit).
		Offset(offset).
		Preloads("Intents").        //TODO fix Not work?!
		Preloads("Intents.Script"). //TODO fix Not work?!
		Preload("OnLaunchScript").
		Preload("OnSessionEndScript").
		Find(&list).Error

	if err != nil {
		return
	}

	//????
	for _, skill := range list {
		_ = n.preload(skill)
	}

	return
}

func (n AlexaSkills) preload(v *AlexaSkill) (err error) {
	err = n.Db.Model(v).
		Related(&v.Intents).Error

	for _, intent := range v.Intents {
		intent.Script = &Script{Id: intent.ScriptId}
		err = n.Db.Model(intent).
			Related(intent.Script).Error
	}
	return
}

// Update ...
func (n AlexaSkills) Update(v *AlexaSkill) (err error) {
	q := map[string]interface{}{
		"skill_id":    v.SkillId,
		"status":      v.Status,
		"description": v.Description,
	}
	if v.OnLaunchScriptId != nil {
		q["on_launch"] = common.Int64Value(v.OnLaunchScriptId)
	}
	if v.OnSessionEndScriptId != nil {
		q["on_session_end"] = common.Int64Value(v.OnSessionEndScriptId)
	}
	err = n.Db.Model(&AlexaSkill{}).Where("id = ?", v.Id).Updates(q).Error
	return
}

// Delete ...
func (n AlexaSkills) Delete(id int64) (err error) {
	err = n.Db.Delete(&AlexaSkill{}, "id = ?", id).Error
	return
}
