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
	"time"
)

// MapLayers ...
type MapLayers struct {
	Db *gorm.DB
}

// MapLayer ...
type MapLayer struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Map         *Map
	MapId       int64
	Status      string
	Weight      int64
	Elements    []*MapElement
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *MapLayer) TableName() string {
	return "map_layers"
}

// Add ...
func (n MapLayers) Add(v *MapLayer) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n MapLayers) GetById(mapId int64) (v *MapLayer, err error) {
	v = &MapLayer{Id: mapId}
	if err = n.Db.First(&v).Error; err != nil {
		err = errors.Wrap(err, "getById failed")
	}
	return
}

// Update ...
func (n MapLayers) Update(m *MapLayer) (err error) {
	err = n.Db.Model(&MapLayer{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"weight":      m.Weight,
		"map_id":      m.MapId,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Sort ...
func (n MapLayers) Sort(m *MapLayer) (err error) {
	err = n.Db.Model(&MapLayer{Id: m.Id}).Updates(map[string]interface{}{
		"weight": m.Weight,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "sort failed")
	}
	return
}

// Delete ...
func (n MapLayers) Delete(mapId int64) (err error) {
	if err = n.Db.Delete(&MapLayer{Id: mapId}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
	}
	return
}

// List ...
func (n *MapLayers) List(limit, offset int64, orderBy, sort string) (list []*MapLayer, total int64, err error) {

	if err = n.Db.Model(MapLayer{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*MapLayer, 0)
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
