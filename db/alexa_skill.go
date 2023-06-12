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

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AlexaSkills ...
type AlexaSkills struct {
	Db *gorm.DB
}

// AlexaSkill ...
type AlexaSkill struct {
	Id          int64 `gorm:"primary_key"`
	SkillId     string
	Description string
	Intents     []*AlexaIntent `gorm:"foreignkey:AlexaSkillId"`
	Status      common.StatusType
	Script      *Script
	ScriptId    *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *AlexaSkill) TableName() string {
	return "alexa_skills"
}

// Add ...
func (n AlexaSkills) Add(v *AlexaSkill) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n AlexaSkills) GetById(id int64) (v *AlexaSkill, err error) {
	v = &AlexaSkill{Id: id}
	err = n.Db.Model(v).
		Preload("Script").
		Preload("Intents").
		Preload("Intents.Script").
		Find(v).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrAlexaSkillNotFound, fmt.Sprintf("id \"%s\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrAlexaSkillGet, err.Error())
		return
	}
	if err = n.preload(v); err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillGet, err.Error())
	}

	return
}

// List ...
func (n *AlexaSkills) List(limit, offset int, orderBy, sort string) (list []*AlexaSkill, total int64, err error) {

	if err = n.Db.Model(AlexaSkill{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillList, err.Error())
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

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillList, err.Error())
	}

	return
}

// ListEnabled ...
func (n *AlexaSkills) ListEnabled(limit, offset int) (list []*AlexaSkill, err error) {

	list = make([]*AlexaSkill, 0)
	err = n.Db.Model(&AlexaSkill{}).
		Where("status = 'enabled'").
		Limit(limit).
		Offset(offset).
		Preload("Intents").
		Preload("Intents.Script").
		Preload("Script").
		Find(&list).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillList, err.Error())
		return
	}

	//????
	for _, skill := range list {
		_ = n.preload(skill)
	}

	return
}

func (n AlexaSkills) preload(v *AlexaSkill) (err error) {
	//todo fix
	//err = n.Db.Model(v).
	//	Related(&v.Intents).Error
	//
	//if err != nil {
	//	err = errors.Wrap(err, "get related intents failed")
	//	return
	//}
	//
	//for _, intent := range v.Intents {
	//	intent.Script = &Script{Id: intent.ScriptId}
	//	err = n.Db.Model(intent).
	//		Related(intent.Script).Error
	//}
	return
}

// Update ...
func (n AlexaSkills) Update(v *AlexaSkill) (err error) {
	q := map[string]interface{}{
		"skill_id":    v.SkillId,
		"status":      v.Status,
		"description": v.Description,
	}
	if v.ScriptId != nil {
		q["script_id"] = common.Int64Value(v.ScriptId)
	}
	if err = n.Db.Model(&AlexaSkill{}).Where("id = ?", v.Id).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillUpdate, err.Error())
	}
	return
}

// Delete ...
func (n AlexaSkills) Delete(id int64) (err error) {
	if err = n.Db.Delete(&AlexaSkill{}, "id = ?", id).Error; err != nil {
		err = errors.Wrap(apperr.ErrAlexaSkillDelete, err.Error())
	}
	return
}
