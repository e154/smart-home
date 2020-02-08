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
	"fmt"
)

type MapTexts struct {
	Db *gorm.DB
}

type MapText struct {
	Id        int64 `gorm:"primary_key"`
	Text      string
	Style     string
}

func (d *MapText) TableName() string {
	return "map_texts"
}

func (n MapTexts) Add(v *MapText) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapTexts) GetById(mapId int64) (v *MapText, err error) {
	v = &MapText{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapTexts) Update(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	return
}

func (n MapTexts) Sort(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	return
}

func (n MapTexts) Delete(id int64) (err error) {

	if err = n.Db.Delete(&MapText{Id: id}).Error; err != nil {
		return
	}

	if id != 0 {
		err = n.Db.Model(&MapElement{}).
			Where("prototype_id = ? and prototype_type = 'text'", id).
			Update("prototype_id", "").
			Error
	}

	return
}

func (n *MapTexts) List(limit, offset int64, orderBy, sort string) (list []*MapText, total int64, err error) {

	if err = n.Db.Model(MapText{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapText, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
