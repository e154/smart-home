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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Entities ...
type Entities struct {
	*Common
}

// Entity ...
type Entity struct {
	Id           pkgCommon.EntityId `gorm:"primary_key"`
	Description  string
	PluginName   string
	Image        *Image
	ImageId      *int64
	States       []*EntityState
	Actions      []*EntityAction
	AreaId       *int64
	Area         *Area
	Metrics      []*Metric `gorm:"many2many:entity_metrics;"`
	Scripts      []*Script `gorm:"many2many:entity_scripts;"`
	Tags         []*Tag    `gorm:"many2many:entity_tags;"`
	Icon         *string
	Payload      json.RawMessage `gorm:"type:jsonb;not null"`
	Settings     json.RawMessage `gorm:"type:jsonb;not null"`
	Storage      []*EntityStorage
	AutoLoad     bool
	RestoreState bool
	ParentId     *pkgCommon.EntityId `gorm:"column:parent_id"`
	CreatedAt    time.Time           `gorm:"<-:create"`
	UpdatedAt    time.Time
}

// TableName ...
func (d *Entity) TableName() string {
	return "entities"
}

type EntitiesStatistic struct {
	Total  int32
	Used   int32
	Unused int32
}

// Add ...
func (n Entities) Add(ctx context.Context, v *Entity) (err error) {

	err = n.DB(ctx).
		Omit("Metrics.*").
		Omit("Tags.*").
		Omit("Scripts.*").
		Create(&v).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "entities_pkey") {
					err = fmt.Errorf("%s: %w", fmt.Sprintf("entity name \"%s\" not unique", v.Id), apperr.ErrEntityAdd)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityAdd)
	}

	return
}

// Update ...
func (n Entities) Update(ctx context.Context, v *Entity) (err error) {

	err = n.DB(ctx).
		Omit("Metrics.*").
		Omit("Tags.*").
		Omit("Scripts.*").
		Save(v).Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityUpdate)
	}

	return
}

// GetById ...
func (n Entities) GetById(ctx context.Context, id pkgCommon.EntityId) (v *Entity, err error) {
	v = &Entity{}
	err = n.DB(ctx).Model(v).
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
		Preload("Tags").
		Preload("Storage", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1).Order("entity_storage.created_at DESC")
		}).
		First(&v).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("id \"%s\"", id), apperr.ErrEntityNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	return
}

// GetByIds ...
func (n Entities) GetByIds(ctx context.Context, ids []pkgCommon.EntityId) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.DB(ctx).Model(Entity{}).
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
		Preload("Tags").
		//Preload("Storage", func(db *gorm.DB) *gorm.DB {
		//	return db.Limit(1).Order("entity_storage.created_at DESC")
		//}).
		Find(&list).Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	if err = n.PreloadStorage(ctx, list); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	return
}

// GetByIdsSimple ...
func (n Entities) GetByIdsSimple(ctx context.Context, ids []pkgCommon.EntityId) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.DB(ctx).Model(Entity{}).
		Preload("States").
		Where("id IN (?)", ids).
		Find(&list).Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	return
}

// Delete ...
func (n Entities) Delete(ctx context.Context, id pkgCommon.EntityId) (err error) {

	if err = n.DB(ctx).Delete(&Entity{Id: id}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityDelete)
		return
	}

	return
}

// List ...
func (n *Entities) List(ctx context.Context, limit, offset int, orderBy, sort string, autoLoad bool,
	query, plugin *string, areaId *int64) (list []*Entity, total int64, err error) {

	list = make([]*Entity, 0)
	q := n.DB(ctx).Model(Entity{})
	if autoLoad {
		q = q.Where("auto_load = ?", true)
	}
	if query != nil {
		q = q.Where("id ILIKE ?", "%"+*query+"%")
	}
	if plugin != nil {
		q = q.Where("plugin_name = ?", *plugin)
	}
	if areaId != nil {
		q = q.Where("area_id = ?", *areaId)
	}
	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityList)
		return
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
		Preload("Tags").
		//Preload("Storage", func(db *gorm.DB) *gorm.DB {
		//	return db.Limit(1).Order("entity_storage.created_at DESC")
		//}).
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityList)
		return
	}

	if err = n.PreloadStorage(ctx, list); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	return
}

// ListPlain ...
func (n *Entities) ListPlain(ctx context.Context, limit, offset int, orderBy, sort string, autoLoad bool,
	query, plugin *string, areaId *int64, tags *[]string) (list []*Entity, total int64, err error) {

	list = make([]*Entity, 0)
	q := n.DB(ctx).Model(Entity{})
	if autoLoad {
		q = q.Where("auto_load = ?", true)
	}
	if query != nil {
		q = q.Where("id ILIKE ?", "%"+*query+"%")
	}
	if plugin != nil {
		q = q.Where("plugin_name = ?", *plugin)
	}
	if areaId != nil {
		q = q.Where("area_id = ?", *areaId)
	}
	if tags != nil {
		q = q.Joins(`left join entity_tags on entity_tags.entity_id = entities.id`)
		q = q.Joins(`left join tags on entity_tags.tag_id = tags.id`)
		q = q.Where("tags.name in (?)", *tags)
	}
	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityList)
		return
	}
	q = q.
		Preload("Tags").
		Preload("Area").
		Group("entities.id").
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityList)
		return
	}

	return
}

// GetByType ...
func (n *Entities) GetByType(ctx context.Context, t string, limit, offset int) (list []*Entity, err error) {

	list = make([]*Entity, 0)
	err = n.DB(ctx).
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
		Preload("Tags").
		//Preload("Storage", func(db *gorm.DB) *gorm.DB {
		//	return db.Order("entity_storage.created_at DESC").Limit(1)
		//}).
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	// todo: remove
	if err = n.PreloadStorage(ctx, list); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityGet)
		return
	}

	return
}

// Search ...
func (n *Entities) Search(ctx context.Context, query string, limit, offset int) (list []*Entity, total int64, err error) {

	q := n.DB(ctx).Model(&Entity{}).
		Where("id ILIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntitySearch)
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("id ASC")

	list = make([]*Entity, 0)
	if err = q.Find(&list).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntitySearch)
	}

	return
}

// UpdateAutoload ...
func (n Entities) UpdateAutoload(ctx context.Context, entityId pkgCommon.EntityId, autoLoad bool) (err error) {
	q := map[string]interface{}{
		"auto_load": autoLoad,
	}

	if err = n.DB(ctx).Model(&Entity{Id: entityId}).Updates(q).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityUpdate)
	}
	return
}

// DeleteScripts ...
func (n Entities) DeleteScripts(ctx context.Context, id pkgCommon.EntityId) (err error) {
	if err = n.DB(ctx).Model(&Entity{Id: id}).Association("Scripts").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityDeleteScript)
	}
	return
}

// DeleteTags ...
func (n Entities) DeleteTags(ctx context.Context, id pkgCommon.EntityId) (err error) {
	if err = n.DB(ctx).Model(&Entity{Id: id}).Association("Tags").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityDeleteTag)
	}
	return
}

// PreloadStorage ...
func (n Entities) PreloadStorage(ctx context.Context, list []*Entity) (err error) {

	//todo: fix
	// temporary solution because Preload("Storage", func(db *gorm.DB) *gorm.DB { - does not work ...
	for _, item := range list {
		err = n.DB(ctx).Model(&EntityStorage{}).
			Order("created_at desc").
			Limit(2).
			Find(&item.Storage, "entity_id = ?", item.Id).
			Error
		if err != nil {
			err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityStorageGet)
			return
		}
	}

	return
}

// Statistic ...
func (n *Entities) Statistic(ctx context.Context) (statistic *EntitiesStatistic, err error) {

	statistic = &EntitiesStatistic{}
	//
	var usedList []struct {
		Count    int32
		AutoLoad bool
	}
	err = n.DB(ctx).Raw(`
select count(e.id), e.auto_load
from entities as e
group by e.auto_load`).
		Scan(&usedList).
		Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrEntityStat)
		return
	}

	for _, item := range usedList {
		statistic.Total += item.Count
		if item.AutoLoad {
			statistic.Used = item.Count

			continue
		}
		statistic.Unused = item.Count
	}

	return
}
