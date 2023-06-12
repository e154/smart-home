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

	. "github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Scripts ...
type Scripts struct {
	Db *gorm.DB
}

// Script ...
type Script struct {
	Id          int64 `gorm:"primary_key"`
	Lang        ScriptLang
	Name        string
	Source      string
	Description string
	Compiled    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Script) TableName() string {
	return "scripts"
}

// Add ...
func (n Scripts) Add(script *Script) (id int64, err error) {
	if err = n.Db.Create(&script).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptAdd, err.Error())
		return
	}
	id = script.Id
	return
}

// GetById ...
func (n Scripts) GetById(scriptId int64) (script *Script, err error) {
	script = &Script{Id: scriptId}
	if err = n.Db.First(&script).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("id \"%d\"", scriptId))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}
	return
}

// GetByName ...
func (n Scripts) GetByName(name string) (script *Script, err error) {
	script = &Script{Name: name}
	if err = n.Db.First(&script).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("name \"%d\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}
	return
}

// Update ...
func (n Scripts) Update(m *Script) (err error) {
	err = n.Db.Model(&Script{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"lang":        m.Lang,
		"source":      m.Source,
		"compiled":    m.Compiled,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrScriptUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Scripts) Delete(scriptId int64) (err error) {
	if err = n.Db.Delete(&Script{Id: scriptId}).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptDelete, err.Error())
	}
	return
}

// List ...
func (n *Scripts) List(limit, offset int, orderBy, sort string) (list []*Script, total int64, err error) {

	if err = n.Db.Model(Script{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptList, err.Error())
		return
	}

	list = make([]*Script, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrScriptList, err.Error())
	}
	return
}

// Search ...
func (n *Scripts) Search(query string, limit, offset int) (list []*Script, total int64, err error) {

	q := n.Db.Model(&Script{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Script, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptSearch, err.Error())
	}
	return
}
