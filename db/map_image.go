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

// MapImages ...
type MapImages struct {
	Db *gorm.DB
}

// MapImage ...
type MapImage struct {
	Id      int64 `gorm:"primary_key"`
	Image   *Image
	ImageId int64
	Style   string
}

// TableName ...
func (d *MapImage) TableName() string {
	return "map_images"
}

// Add ...
func (n MapImages) Add(v *MapImage) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n MapImages) GetById(mapId int64) (v *MapImage, err error) {
	v = &MapImage{Id: mapId}
	if err = n.Db.First(&v).Error; err != nil {
		err = errors.Wrap(err, "getById failed")
	}
	return
}

// Update ...
func (n MapImages) Update(m *MapImage) (err error) {
	err = n.Db.Model(&MapImage{Id: m.Id}).Updates(map[string]interface{}{
		"image_id": m.ImageId,
		"style":    m.Style,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Sort ...
func (n MapImages) Sort(m *MapImage) (err error) {
	err = n.Db.Model(&MapImage{Id: m.Id}).Updates(map[string]interface{}{
		"image_id": m.ImageId,
		"style":    m.Style,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "sort failed")
	}
	return
}

// Delete ...
func (n MapImages) Delete(id int64) (err error) {

	if err = n.Db.Delete(&MapImage{Id: id}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
		return
	}

	if id != 0 {
		err = n.Db.Model(&MapElement{}).
			Where("prototype_id = ? and prototype_type = 'image'", id).
			Update("prototype_id", "").
			Error
	}
	if err != nil {
		err = errors.Wrap(err, "update mapElement failed")
	}
	return
}

// List ...
func (n *MapImages) List(limit, offset int64, orderBy, sort string) (list []*MapImage, total int64, err error) {

	if err = n.Db.Model(MapImage{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count mapElement failed")
		return
	}

	list = make([]*MapImage, 0)
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
