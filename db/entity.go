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

// Entities ...
type Entities struct {
	Db *gorm.DB
}

// Entity ...
type Entity struct {
	Id          common.EntityId `gorm:"primary_key"`
	Description string
	PluginName  string
	Image       *Image
	ImageId     *int64
	States      []*EntityState
	Actions     []*EntityAction
	AreaId      *int64
	Area        *Area
	Metrics     []Metric `gorm:"many2many:entity_metrics;"`
	Scripts     []Script `gorm:"many2many:entity_scripts;"`
	Icon        *string
	Payload     json.RawMessage `gorm:"type:jsonb;not null"`
	Settings    json.RawMessage `gorm:"type:jsonb;not null"`
	Storage     []*EntityStorage
	AutoLoad    bool
	ParentId    *common.EntityId `gorm:"column:parent"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Entity) TableName() string {
	return "entities"
}

// Add ...
func (n Entities) Add(v *Entity) (err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
	}
	return
}

// Update ...
func (n Entities) Update(v *Entity) (err error) {
	q := map[string]interface{}{
		"image_id":    v.ImageId,
		"area_id":     v.AreaId,
		"description": v.Description,
		"plugin_name": v.PluginName,
		"icon":        v.Icon,
		"payload":     v.Payload,
		"settings":    v.Settings,
		"auto_load":   v.AutoLoad,
	}

	if err = n.Db.Model(&Entity{Id: v.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// UpdateSettings ...
func (n Entities) UpdateSettings(entityId common.EntityId, settings []byte) (err error) {
	q := map[string]interface{}{
		"settings": settings,
	}

	if err = n.Db.Model(&Entity{Id: entityId}).Updates(q).Error; err != nil {
		err = errors.Wrap(err, "update settings failed")
	}
	return
}

// GetById ...
func (n Entities) GetById(id common.EntityId) (v *Entity, err error) {
	v = &Entity{}
	err = n.Db.Model(v).
		Where("id = ?", id).
		Preload("Image").
		Preload("States").
		Preload("States.Image").
		Preload("Actions").
		Preload("Actions.Image").
		Preload("Actions.Script").
		Preload("Area").
		Preload("Metrics").
		Preload("Scripts").
		First(&v).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("id \"%s\"", id))
			return
		}
		err = errors.Wrap(err, "getById failed")
		return
	}

	if err = n.preloadStorage(v); err != nil {
		err = errors.Wrap(err, "preload failed")
		return
	}

	return
}

// Delete ...
func (n Entities) Delete(id common.EntityId) (err error) {

	if err = n.Db.Delete(&Entity{Id: id}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
		return
	}

	err = n.Db.Model(&MapElement{}).
		Where("prototype_id = ? and prototype_type = 'entity'", id).
		Update("prototype_id", "").
		Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
	}

	return
}

// List ...
func (n *Entities) List(limit, offset int64, orderBy, sort string, autoLoad bool) (list []*Entity, total int64, err error) {

	if err = n.Db.Model(Entity{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*Entity, 0)
	q := n.Db
	if autoLoad {
		q = q.Where("auto_load = ?", true)
	}
	q = q.
		Preload("Image").
		Preload("States").
		Preload("States.Image").
		Preload("Actions").
		Preload("Actions.Image").
		Preload("Actions.Script").
		Preload("Area").
		Preload("Metrics").
		Preload("Scripts").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(err, "find failed")
		return
	}

	for _, entity := range list {
		if err = n.preloadStorage(entity); err != nil {
			err = errors.Wrap(err, "preload failed")
			return
		}
	}

	return
}

// GetByType ...
func (n *Entities) GetByType(t string, limit, offset int64) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.Db.Model(&Entity{}).
		Where("plugin_name = ? and auto_load = true", t).
		Preload("Image").
		Preload("States").
		Preload("States.Image").
		Preload("Actions").
		Preload("Actions.Image").
		Preload("Actions.Script").
		Preload("Area").
		Preload("Metrics").
		Preload("Scripts").
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("type \"%s\"", t))
			return
		}
		err = errors.Wrap(err, "getByType failed")
		return
	}

	for _, entity := range list {
		if err = n.preloadStorage(entity); err != nil {
			err = errors.Wrap(err, "preload failed")
			return
		}
	}

	return
}

// Search ...
func (n *Entities) Search(query string, limit, offset int64) (list []*Entity, total int64, err error) {

	q := n.Db.Model(&Entity{}).
		Where("id LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("id ASC")

	list = make([]*Entity, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(err, "search failed")
	}

	return
}

// AppendMetric ...
func (n Entities) AppendMetric(id common.EntityId, metric Metric) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Metrics").Append(&metric).Error; err != nil {
		err = errors.Wrap(err, "append metric failed")
	}
	return
}

// DeleteMetric ...
func (n Entities) DeleteMetric(id common.EntityId, metricId int64) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Metrics").Delete(&Metric{Id: metricId}).Error; err != nil {
		err = errors.Wrap(err, "delete metric failed")
	}
	return
}

// ReplaceMetric ...
func (n Entities) ReplaceMetric(id common.EntityId, metric Metric) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Metrics").Replace(&metric).Error; err != nil {
		err = errors.Wrap(err, "replace metric failed")
	}
	return
}

// AppendScript ...
func (n Entities) AppendScript(id common.EntityId, script *Script) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Scripts").Append(script).Error; err != nil {
		err = errors.Wrap(err, "append script failed")
	}
	return
}

// DeleteScript ...
func (n Entities) DeleteScript(id common.EntityId, scriptId int64) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Scripts").Delete(&Script{Id: scriptId}).Error; err != nil {
		err = errors.Wrap(err, "delete metric failed")
	}
	return
}

// ReplaceScript ...
func (n Entities) ReplaceScript(id common.EntityId, script *Script) (err error) {
	if err = n.Db.Model(&Entity{Id: id}).Association("Scripts").Replace(script).Error; err != nil {
		err = errors.Wrap(err, "replace metric failed")
	}
	return
}

func (n *Entities) preloadStorage(entity *Entity) (err error) {
	err = n.Db.Model(&EntityStorage{}).
		Where("entity_id = ?", entity.Id).
		Order("created_at DESC").
		Limit(1).
		Scan(&entity.Storage).
		Error
	return
}
