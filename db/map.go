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
	"time"

	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Maps ...
type Maps struct {
	Db *gorm.DB
}

// Map ...
type Map struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Options     json.RawMessage `gorm:"type:jsonb;not null"`
	Layers      []*MapLayer
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Map) TableName() string {
	return "maps"
}

// Add ...
func (n Maps) Add(v *Map) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Maps) GetById(mapId int64) (v *Map, err error) {
	v = &Map{Id: mapId}
	if err = n.Db.First(&v).Error; err != nil {
		err = errors.Wrap(err, "getById failed")
	}
	return
}

// GetFullById ...
func (n Maps) GetFullById(mapId int64) (v *Map, err error) {

	v = &Map{}
	err = n.Db.Model(v).
		Where("id = ?", mapId).
		Preload("Layers").
		Preload("Layers.Elements").
		Preload("Layers.Elements.Zone").
		Preload("Layers.Map").
		Find(&v).
		Error

	if err != nil {
		err = errors.Wrap(err, "get map failed")
		return
	}

	imageIds := make([]int64, 0)
	textIds := make([]int64, 0)
	deviceIds := make([]int64, 0)
	for _, l := range v.Layers {
		for _, e := range l.Elements {
			id, ok := e.PrototypeId.(int64)
			if !ok {
				continue
			}
			switch e.PrototypeType {
			case common.MapElementPrototypeText:
				textIds = append(textIds, id)
			case common.MapElementPrototypeImage:
				imageIds = append(imageIds, id)
			case common.MapElementPrototypeEntity:
				deviceIds = append(deviceIds, id)
			}
		}
	}

	images := make([]*MapImage, 0)
	texts := make([]*MapText, 0)
	entities := make([]*Entity, 0)

	if len(imageIds) > 0 {
		err = n.Db.Model(&MapImage{}).
			Where("id in (?)", imageIds).
			Preload("Image").
			Find(&images).
			Error
		if err != nil {
			err = errors.Wrap(err, "get mapImage failed")
		}
	}

	if len(textIds) > 0 {
		err = n.Db.Model(&MapText{}).
			Where("id in (?)", textIds).
			Find(&texts).
			Error
		if err != nil {
			err = errors.Wrap(err, "get mapText failed")
		}
	}

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
			//Preload("Device").
			//Preload("Device.States").
			//Preload("Device.Actions").
			Preload("Zone").
			Find(&entities).
			Error
		if err != nil {
			err = errors.Wrap(err, "get mapEntity failed")
		}
	}

	for _, l := range v.Layers {
		for _, e := range l.Elements {

			if err = n.Db.Preload("Metrics").First(e).Error; err != nil {
				err = errors.Wrap(err, "preload metric failed")
				return
			}

			switch e.PrototypeType {
			case common.MapElementPrototypeText:
				for _, text := range texts {
					if text.Id == e.PrototypeId {
						e.Prototype = Prototype{
							MapText: text,
						}
						continue
					}
				}
			case common.MapElementPrototypeImage:
				for _, image := range images {
					if image.Id == e.PrototypeId {
						e.Prototype = Prototype{
							MapImage: image,
						}
						continue
					}
				}
			case common.MapElementPrototypeEntity:
				for _, entity := range entities {
					if entity.Id == e.PrototypeId {
						e.Prototype = Prototype{
							Entity: entity,
						}
						continue
					}
				}
			}
		}
	}

	//debug.Println(v)
	//fmt.Println("---")

	return
}

// Update ...
func (n Maps) Update(m *Map) (err error) {
	err = n.Db.Model(&Map{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"options":     m.Options,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Delete ...
func (n Maps) Delete(mapId int64) (err error) {
	if err = n.Db.Delete(&Map{Id: mapId}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
	}
	return
}

// List ...
func (n *Maps) List(limit, offset int64, orderBy, sort string) (list []*Map, total int64, err error) {

	if err = n.Db.Model(Map{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*Map, 0)
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

// Search ...
func (n *Maps) Search(query string, limit, offset int) (list []*Map, total int64, err error) {

	q := n.Db.Model(&Map{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Map, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(err, "search failed")
	}
	return
}
