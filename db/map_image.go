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

type MapImages struct {
	Db *gorm.DB
}

type MapImage struct {
	Id        int64 `gorm:"primary_key"`
	Image     *Image
	ImageId   int64
	Style     string
}

func (d *MapImage) TableName() string {
	return "map_images"
}

func (n MapImages) Add(v *MapImage) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapImages) GetById(mapId int64) (v *MapImage, err error) {
	v = &MapImage{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapImages) Update(m *MapImage) (err error) {
	err = n.Db.Model(&MapImage{Id: m.Id}).Updates(map[string]interface{}{
		"image_id": m.ImageId,
		"style":    m.Style,
	}).Error
	return
}

func (n MapImages) Sort(m *MapImage) (err error) {
	err = n.Db.Model(&MapImage{Id: m.Id}).Updates(map[string]interface{}{
		"image_id": m.ImageId,
		"style":    m.Style,
	}).Error
	return
}

func (n MapImages) Delete(id int64) (err error) {

	if err = n.Db.Delete(&MapImage{Id: id}).Error; err != nil {
		return
	}

	if id != 0 {
		err = n.Db.Model(&MapElement{}).
			Where("prototype_id = ? and prototype_type = 'image'", id).
			Update("prototype_id", "").
			Error
	}

	return
}

func (n *MapImages) List(limit, offset int64, orderBy, sort string) (list []*MapImage, total int64, err error) {

	if err = n.Db.Model(MapImage{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapImage, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
