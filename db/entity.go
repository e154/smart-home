// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
	Metrics     []*Metric `gorm:"many2many:entity_metrics;"`
	Scripts     []*Script `gorm:"many2many:entity_scripts;"`
	Icon        *string
	Payload     json.RawMessage `gorm:"type:jsonb;not null"`
	Settings    json.RawMessage `gorm:"type:jsonb;not null"`
	Storage     []*EntityStorage
	AutoLoad    bool
	ParentId    *common.EntityId `gorm:"column:parent_id"`
	CreatedAt   time.Time        `gorm:"<-:create"`
	UpdatedAt   time.Time
}

// TableName ...
func (d *Entity) TableName() string {
	return "entities"
}

// Add ...
func (n Entities) Add(ctx context.Context, v *Entity) (err error) {
	err = n.Db.WithContext(ctx).Omit("Metrics.*").Omit("Scripts.*").Create(&v).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "entities_pkey") {
					err = errors.Wrap(apperr.ErrEntityAdd, fmt.Sprintf("entity name \"%s\" not unique", v.Id))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrEntityAdd, err.Error())
	}
	return
}

// Update ...
func (n Entities) Update(ctx context.Context, v *Entity) (err error) {

	err = n.Db.WithContext(ctx).
		Omit("Metrics.*").
		Omit("Scripts.*").
		Save(v).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityUpdate, err.Error())
	}
	return
}

// GetById ...
func (n Entities) GetById(ctx context.Context, id common.EntityId) (v *Entity, err error) {
	v = &Entity{}
	err = n.Db.WithContext(ctx).Model(v).
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
		Preload("Storage", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1).Order("created_at DESC")
		}).
		First(&v).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrEntityNotFound, fmt.Sprintf("id \"%s\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrEntityGet, err.Error())
		return
	}

	return
}

// GetByIds ...
func (n Entities) GetByIds(ctx context.Context, ids []common.EntityId) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.Db.WithContext(ctx).Model(Entity{}).
		Where("id IN (?)", ids).
		Preload("Image").
		Preload("States").
		Preload("States.Image").
		Preload("Actions").
		Preload("Actions.Image").
		Preload("Actions.Script").
		Preload("Area").
		Preload("Metrics").
		Preload("Scripts").
		Preload("Storage", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1).Order("created_at DESC")
		}).
		Find(&list).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityGet, err.Error())
		return
	}

	return
}

// GetByIdsSimple ...
func (n Entities) GetByIdsSimple(ctx context.Context, ids []common.EntityId) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.Db.WithContext(ctx).Model(Entity{}).
		Preload("States").
		Where("id IN (?)", ids).
		Find(&list).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityGet, err.Error())
		return
	}

	return
}

// Delete ...
func (n Entities) Delete(ctx context.Context, id common.EntityId) (err error) {

	if err = n.Db.WithContext(ctx).Delete(&Entity{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityDelete, err.Error())
		return
	}

	return
}

// List ...
func (n *Entities) List(ctx context.Context, limit, offset int, orderBy, sort string, autoLoad bool,
	query, plugin *string, areaId *int64) (list []*Entity, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(Entity{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityList, err.Error())
		return
	}

	list = make([]*Entity, 0)
	q := n.Db
	if autoLoad {
		q = q.Where("auto_load = ?", true)
	}
	if query != nil {
		q = q.Where("id LIKE ?", "%"+*query+"%")
	}
	if plugin != nil {
		q = q.Where("plugin_name = ?", *plugin)
	}
	if areaId != nil {
		q = q.Where("area_id = ?", *areaId)
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
		Preload("Storage", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1).Order("created_at DESC")
		}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		WithContext(ctx).
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityList, err.Error())
		return
	}

	return
}

// GetByType ...
func (n *Entities) GetByType(ctx context.Context, t string, limit, offset int) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.Db.WithContext(ctx).
		Model(&Entity{}).
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
		Preload("Storage", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1).Order("created_at DESC")
		}).
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrEntityGet, fmt.Sprintf("type \"%s\"", t))
			return
		}
		err = errors.Wrap(apperr.ErrEntityGet, err.Error())
		return
	}

	return
}

// Search ...
func (n *Entities) Search(ctx context.Context, query string, limit, offset int) (list []*Entity, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Entity{}).
		Where("id LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntitySerch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("id ASC")

	list = make([]*Entity, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntitySerch, err.Error())
	}

	return
}

// UpdateAutoload ...
func (n Entities) UpdateAutoload(ctx context.Context, entityId common.EntityId, autoLoad bool) (err error) {
	q := map[string]interface{}{
		"auto_load": autoLoad,
	}

	if err = n.Db.WithContext(ctx).Model(&Entity{Id: entityId}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityUpdate, err.Error())
	}
	return
}

// DeleteScripts ...
func (n Entities) DeleteScripts(ctx context.Context, id common.EntityId) (err error) {
	if err = n.Db.WithContext(ctx).Model(&Entity{Id: id}).Association("Scripts").Clear(); err != nil {
		err = errors.Wrap(apperr.ErrEntityDeleteScript, err.Error())
	}
	return
}
