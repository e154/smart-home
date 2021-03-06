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
	"encoding/json"
	"fmt"
	. "github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

// MapElements ...
type MapElements struct {
	Db *gorm.DB
}

// Prototype ...
type Prototype struct {
	*MapImage
	*MapText
	*Entity
}

// Entity ...
type MapElement struct {
	Id            int64 `gorm:"primary_key"`
	Name          string
	Description   string
	PrototypeId   MapElementPrototypeId
	PrototypeType MapElementPrototypeType
	Prototype     Prototype
	Map           *Map
	MapId         int64
	MapLayer      *MapLayer
	MapLayerId    int64
	GraphSettings json.RawMessage `gorm:"type:jsonb;not null"`
	Status        StatusType
	Weight        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName ...
func (d *MapElement) TableName() string {
	return "map_elements"
}

// Add ...
func (n MapElements) Add(v *MapElement) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n MapElements) GetById(mapId int64) (v *MapElement, err error) {
	v = &MapElement{Id: mapId}
	if err = n.Db.
		First(&v).Error; err != nil {
		return
	}

	err = n.gePrototype(v)
	return
}

// GetByName ...
func (n MapElements) GetByName(name string) (v *MapElement, err error) {
	v = &MapElement{}
	if err = n.Db.
		Where("name = ?", name).
		First(&v).Error; err != nil {
		return
	}

	err = n.gePrototype(v)
	return
}

func (n MapElements) gePrototype(v *MapElement) (err error) {

	if v.PrototypeId == "" {
		return
	}

	switch v.PrototypeType {
	case MapElementPrototypeText:
		text := &MapText{}
		err = n.Db.Model(&MapText{}).
			Where("id = ?", v.PrototypeId).
			First(&text).
			Error

		if err == nil {
			v.Prototype = Prototype{
				MapText: text,
			}
		}
	case MapElementPrototypeImage:
		image := &MapImage{}
		err = n.Db.Model(&MapImage{}).
			Where("id = ?", v.PrototypeId).
			First(&image).
			Error
		if err == nil {
			v.Prototype = Prototype{
				MapImage: image,
			}
		}
	case MapElementPrototypeEntity:
		device := &Entity{}
		err = n.Db.Model(&Entity{}).
			Where("id = ?", v.PrototypeId).
			Preload("Image").
			Preload("States").
			Preload("States.Image").
			Preload("States.DeviceState").
			Preload("Actions").
			Preload("Actions.Image").
			Preload("Actions.DeviceAction").
			Preload("Area").
			First(&device).
			Error
		if err == nil {
			v.Prototype = Prototype{
				Entity: device,
			}
		}
	}

	return
}

// Update ...
func (n MapElements) Update(m *MapElement) (err error) {
	err = n.Db.Model(&MapElement{Id: m.Id}).Updates(map[string]interface{}{
		"name":           m.Name,
		"description":    m.Description,
		"prototype_id":   m.PrototypeId,
		"prototype_type": m.PrototypeType,
		"map_id":         m.MapId,
		"layer_id":       m.MapLayerId,
		"graph_settings": m.GraphSettings,
		"status":         m.Status,
		"weight":         m.Weight,
	}).Error
	return
}

// Sort ...
func (n MapElements) Sort(m *MapElement) (err error) {
	err = n.Db.Model(&MapElement{Id: m.Id}).Updates(map[string]interface{}{
		"weight": m.Weight,
	}).Error
	return
}

// Delete ...
func (n MapElements) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapElement{Id: mapId}).Error
	return
}

// List ...
func (n *MapElements) List(limit, offset int64, orderBy, sort string) (list []*MapElement, total int64, err error) {

	if err = n.Db.Model(MapElement{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapElement, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

// GetActiveElements ...
func (n *MapElements) GetActiveElements(limit, offset int64, orderBy, sort string) (list []*MapElement, total int64, err error) {

	list = make([]*MapElement, 0)

	q := n.Db.Model(&MapElement{}).
		Where("prototype_type = 'entity'")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = n.Db.Model(&MapElement{}).
		Where("prototype_type = 'entity'").
		Limit(limit).
		Offset(offset)

	if orderBy != "" && sort != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		return
	}

	deviceIds := make([]MapElementPrototypeId, 0)
	for _, e := range list {
		deviceIds = append(deviceIds, e.PrototypeId)
	}

	devices := make([]*Entity, 0)
	if len(deviceIds) > 0 {
		err = n.Db.Model(&Entity{}).
			Where("id in (?)", deviceIds).
			Preload("Image").
			Preload("States").
			Preload("States.Image").
			Preload("States.DeviceState").
			Preload("Actions").
			Preload("Actions.Image").
			Preload("Actions.DeviceAction").
			Preload("Area").
			Find(&devices).
			Error

		if err != nil {
			return
		}
	}

	for _, e := range list {
		for _, device := range devices {
			if device.Id == e.PrototypeId {
				e.Prototype = Prototype{
					Entity: device,
				}
				continue
			}
		}
	}

	return
}
