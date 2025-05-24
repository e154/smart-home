// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

package adaptors

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.EntityRepo = (*Entity)(nil)

// Entity ...
type Entity struct {
	table *db.Entities
	db    *gorm.DB
	orm   *orm.Orm
}

// GetEntityAdaptor ...
func GetEntityAdaptor(d *gorm.DB, orm *orm.Orm) *Entity {
	return &Entity{
		table: &db.Entities{&db.Common{Db: d}},
		db:    d,
		orm:   orm,
	}
}

// Add ...
func (n *Entity) Add(ctx context.Context, ver *models.Entity) (err error) {

	if strings.Contains(ver.Id.String(), " ") {
		err = fmt.Errorf("entity name \"%s\" contains spaces: %w", ver.Id, apperr.ErrEntityAdd)
		return
	}

	err = n.table.Add(ctx, n.toDb(ver))

	return
}

// GetById ...
func (n *Entity) GetById(ctx context.Context, id pkgCommon.EntityId, preloadMetric ...bool) (ver *models.Entity, err error) {

	var dbVer *db.Entity
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	if len(preloadMetric) > 0 && preloadMetric[0] {
		n.preloadMetric(ctx, ver)
	}

	return
}

// GetByIds ...
func (n *Entity) GetByIds(ctx context.Context, ids []pkgCommon.EntityId, preloadMetric ...bool) (list []*models.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByIds(ctx, ids); err != nil {
		return
	}
	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		if len(preloadMetric) > 0 && preloadMetric[0] {
			n.preloadMetric(ctx, ver)
		}
		list[i] = ver
	}

	return
}

// GetByIdsSimple ...
func (n *Entity) GetByIdsSimple(ctx context.Context, ids []pkgCommon.EntityId) (list []*models.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByIdsSimple(ctx, ids); err != nil {
		return
	}
	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Delete ...
func (n *Entity) Delete(ctx context.Context, id pkgCommon.EntityId) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// List ...
func (n *Entity) List(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query, plugin *string, areaId *int64) (list []*models.Entity, total int64, err error) {

	var dbList []*db.Entity
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, autoLoad, query, plugin, areaId); err != nil {
		return
	}

	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// ListPlain ...
func (n *Entity) ListPlain(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query,
	plugin *string, areaId *int64, tags *[]string) (list []*models.Entity, total int64, err error) {

	var dbList []*db.Entity
	if dbList, total, err = n.table.ListPlain(ctx, int(limit), int(offset), orderBy, sort, autoLoad, query, plugin,
		areaId, tags); err != nil {
		return
	}

	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// GetByType ...
func (n *Entity) GetByType(ctx context.Context, t string, limit, offset int64) (list []*models.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByType(ctx, t, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

func (n *Entity) DeleteScripts(ctx context.Context, entityID pkgCommon.EntityId) (err error) {
	err = n.table.DeleteScripts(ctx, entityID)
	return
}

func (n *Entity) DeleteTags(ctx context.Context, entityID pkgCommon.EntityId) (err error) {
	err = n.table.DeleteTags(ctx, entityID)
	return
}

// Update ...
func (n *Entity) Update(ctx context.Context, ver *models.Entity) (err error) {
	err = n.table.Update(ctx, n.toDb(ver))
	return
}

// Search ...
func (n *Entity) Search(ctx context.Context, query string, limit, offset int64) (list []*models.Entity, total int64, err error) {
	var dbList []*db.Entity
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*models.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// UpdateAutoload ...
func (n *Entity) UpdateAutoload(ctx context.Context, entityId pkgCommon.EntityId, autoLoad bool) (err error) {
	err = n.table.UpdateAutoload(ctx, entityId, autoLoad)
	return
}

func (n *Entity) preloadMetric(ctx context.Context, ver *models.Entity) {

	var err error

	// load preview metrics data
	if ver.Metrics == nil || len(ver.Metrics) == 0 {
		return
	}
	bucketMetricBucketAdaptor := GetMetricBucketAdaptor(n.db, n.orm)
	for i, metric := range ver.Metrics {

		var optionItems = make([]string, len(metric.Options.Items))
		for i, item := range metric.Options.Items {
			optionItems[i] = item.Name
		}

		if ver.Metrics[i].Data, err = bucketMetricBucketAdaptor.List(ctx, nil, nil, metric.Id, optionItems, pkgCommon.MetricRange24H.Ptr()); err != nil {
			log.Error(err.Error())
			return
		}

		ver.Metrics[i].RangesByType()
	}
}

func (n *Entity) Statistic(ctx context.Context) (statistic *models.EntitiesStatistic, err error) {
	var dbVer *db.EntitiesStatistic
	if dbVer, err = n.table.Statistic(ctx); err != nil {
		return
	}
	statistic = &models.EntitiesStatistic{
		Total:  dbVer.Total,
		Used:   dbVer.Used,
		Unused: dbVer.Unused,
	}
	return
}

func (n *Entity) fromDb(dbVer *db.Entity) (ver *models.Entity) {
	ver = &models.Entity{
		Id:           dbVer.Id,
		Description:  dbVer.Description,
		PluginName:   dbVer.PluginName,
		Actions:      make([]*models.EntityAction, 0),
		States:       make([]*models.EntityState, 0),
		Icon:         dbVer.Icon,
		AutoLoad:     dbVer.AutoLoad,
		RestoreState: dbVer.RestoreState,
		ParentId:     dbVer.ParentId,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
	}

	// actions
	entityActionAdaptor := GetEntityActionAdaptor(n.db)
	for _, dbAction := range dbVer.Actions {
		action := entityActionAdaptor.fromDb(dbAction)
		ver.Actions = append(ver.Actions, action)
	}

	// states
	entityStateAdaptor := GetEntityStateAdaptor(n.db)
	for _, dbState := range dbVer.States {
		state := entityStateAdaptor.fromDb(dbState)
		ver.States = append(ver.States, state)
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// Area
	if dbVer.Area != nil {
		areaAdaptor := GetAreaAdaptor(n.db)
		ver.Area = areaAdaptor.fromDb(dbVer.Area)
		ver.AreaId = pkgCommon.Int64(dbVer.Area.Id)
	}

	// metrics
	if dbVer.Metrics != nil && len(dbVer.Metrics) > 0 {
		metricAdaptor := GetMetricAdaptor(n.db, nil)
		for _, metric := range dbVer.Metrics {
			ver.Metrics = append(ver.Metrics, metricAdaptor.fromDb(metric))
		}
	}

	// scripts
	if dbVer.Scripts != nil && len(dbVer.Scripts) > 0 {
		scriptAdaptor := GetScriptAdaptor(n.db)
		for _, script := range dbVer.Scripts {
			s, _ := scriptAdaptor.fromDb(script)
			ver.Scripts = append(ver.Scripts, s)
		}
	} else {
		ver.Scripts = make([]*models.Script, 0)
	}

	// deserialize payload
	b, _ := dbVer.Payload.MarshalJSON()
	payload := models.EntityPayload{}
	_ = json.Unmarshal(b, &payload)
	ver.Attributes = payload.AttributeSignature

	// deserialize settings
	b, _ = dbVer.Settings.MarshalJSON()
	settings := models.EntitySettings{}
	_ = json.Unmarshal(b, &settings)
	ver.Settings = settings.Settings

	// storage
	if len(dbVer.Storage) > 0 {
		data := map[string]interface{}{}
		_ = json.Unmarshal(dbVer.Storage[0].Attributes, &data)
		_, _ = ver.Attributes.Deserialize(data)
		storageAdaptor := GetEntityStorageAdaptor(n.db)
		for _, store := range dbVer.Storage {
			ver.Storage = append(ver.Storage, storageAdaptor.fromDb(store))
		}
	}

	// tags
	for _, tag := range dbVer.Tags {
		ver.Tags = append(ver.Tags, &models.Tag{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}

	return
}

func (n *Entity) toDb(ver *models.Entity) (dbVer *db.Entity) {

	dbVer = &db.Entity{
		Id:           ver.Id,
		Description:  ver.Description,
		PluginName:   ver.PluginName,
		Icon:         ver.Icon,
		AutoLoad:     ver.AutoLoad,
		RestoreState: ver.RestoreState,
		ParentId:     ver.ParentId,
		AreaId:       ver.AreaId,
		ImageId:      ver.ImageId,
		CreatedAt:    ver.CreatedAt,
		UpdatedAt:    ver.UpdatedAt,
	}

	// serialize payload
	b, _ := json.Marshal(models.EntityPayload{
		AttributeSignature: ver.Attributes.Signature(),
	})
	_ = dbVer.Payload.UnmarshalJSON(b)

	// serialize settings
	b, _ = json.Marshal(models.EntitySettings{
		Settings: ver.Settings,
	})
	_ = dbVer.Settings.UnmarshalJSON(b)

	// states
	entityState := GetEntityStateAdaptor(nil)
	if len(ver.States) > 0 {
		for i := range ver.States {
			ver.States[i].EntityId = ver.Id
		}
		dbVer.States = make([]*db.EntityState, 0, len(ver.States))
		for _, state := range ver.States {
			dbVer.States = append(dbVer.States, entityState.toDb(state))
		}
	} else {
		dbVer.States = make([]*db.EntityState, 0)
	}

	// actions
	entityAction := GetEntityActionAdaptor(nil)
	if len(ver.Actions) > 0 {
		for i := range ver.Actions {
			ver.Actions[i].EntityId = ver.Id
		}
		dbVer.Actions = make([]*db.EntityAction, 0, len(ver.Actions))
		for _, action := range ver.Actions {
			dbVer.Actions = append(dbVer.Actions, entityAction.toDb(action))
		}
	} else {
		dbVer.Actions = make([]*db.EntityAction, 0)
	}

	// metrics
	if len(ver.Metrics) > 0 {
		metricAdaptor := GetMetricAdaptor(nil, nil)
		dbVer.Metrics = make([]*db.Metric, 0, len(ver.Metrics))
		for _, metric := range ver.Metrics {
			dbVer.Metrics = append(dbVer.Metrics, metricAdaptor.toDb(metric))
		}
	} else {
		dbVer.Metrics = make([]*db.Metric, 0)
	}

	// scripts
	if len(ver.Scripts) > 0 {
		dbVer.Scripts = make([]*db.Script, 0, len(ver.Scripts))
		for _, script := range ver.Scripts {
			dbVer.Scripts = append(dbVer.Scripts, &db.Script{
				Id: script.Id,
			})
		}
	} else {
		dbVer.Scripts = make([]*db.Script, 0)
	}

	// tags
	if len(ver.Tags) > 0 {
		dbVer.Tags = make([]*db.Tag, 0, len(ver.Tags))
		for _, tag := range ver.Tags {
			dbVer.Tags = append(dbVer.Tags, &db.Tag{
				Id:   tag.Id,
				Name: tag.Name,
			})
		}
	} else {
		dbVer.Tags = make([]*db.Tag, 0)
	}

	return
}
