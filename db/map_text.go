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

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// MapTexts ...
type MapTexts struct {
	Db *gorm.DB
}

// MapText ...
type MapText struct {
	Id    int64 `gorm:"primary_key"`
	Text  string
	Style string
}

// TableName ...
func (d *MapText) TableName() string {
	return "map_texts"
}

// Add ...
func (n MapTexts) Add(v *MapText) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n MapTexts) GetById(mapId int64) (v *MapText, err error) {
	v = &MapText{Id: mapId}
	if err = n.Db.First(&v).Error; err != nil {
		err = errors.Wrap(err, "geById failed")
	}
	return
}

// Update ...
func (n MapTexts) Update(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Sort ...
func (n MapTexts) Sort(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "sort failed")
	}
	return
}

// Delete ...
func (n MapTexts) Delete(id int64) (err error) {

	if err = n.Db.Delete(&MapText{Id: id}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
		return
	}

	if id != 0 {
		err = n.Db.Model(&MapElement{}).
			Where("prototype_id = ? and prototype_type = 'text'", id).
			Update("prototype_id", "").
			Error
	}
	if err != nil {
		err = errors.Wrap(err, "update mapElement failed")
	}
	return
}

// List ...
func (n *MapTexts) List(limit, offset int64, orderBy, sort string) (list []*MapText, total int64, err error) {

	if err = n.Db.Model(MapText{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*MapText, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(err, "list failed")
	}
	return
}
