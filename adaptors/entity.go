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

package adaptors

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IEntity ...
type IEntity interface {
	Add(ctx context.Context, ver *m.Entity) (err error)
	GetById(ctx context.Context, id common.EntityId, preloadMetric ...bool) (ver *m.Entity, err error)
	GetByIds(ctx context.Context, ids []common.EntityId, preloadMetric ...bool) (ver []*m.Entity, err error)
	GetByIdsSimple(ctx context.Context, ids []common.EntityId) (list []*m.Entity, err error)
	Delete(ctx context.Context, id common.EntityId) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query, plugin *string, areaId *int64) (list []*m.Entity, total int64, err error)
	GetByType(ctx context.Context, t string, limit, offset int64) (list []*m.Entity, err error)
	Update(ctx context.Context, ver *m.Entity) (err error)
	UpdateSettings(ctx context.Context, entityId common.EntityId, settings m.Attributes) (err error)
	Search(ctx context.Context, query string, limit, offset int64) (list []*m.Entity, total int64, err error)
	AppendMetric(ctx context.Context, entityId common.EntityId, metric *m.Metric) (err error)
	DeleteMetric(ctx context.Context, entityId common.EntityId, metric *m.Metric) (err error)
	UpdateAutoload(ctx context.Context, entityId common.EntityId, autoLoad bool) (err error)
	Import(ctx context.Context, entity *m.Entity) (err error)
	preloadMetric(ctx context.Context, ver *m.Entity)
	fromDb(dbVer *db.Entity) (ver *m.Entity)
	toDb(ver *m.Entity) (dbVer *db.Entity)
}

// Entity ...
type Entity struct {
	IEntity
	table *db.Entities
	db    *gorm.DB
}

// GetEntityAdaptor ...
func GetEntityAdaptor(d *gorm.DB) IEntity {
	return &Entity{
		table: &db.Entities{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Entity) Add(ctx context.Context, ver *m.Entity) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.Entities{Db: tx}
	if err = table.Add(ctx, n.toDb(ver)); err != nil {
		return
	}

	//actions
	if len(ver.Actions) > 0 {
		for i := range ver.Actions {
			ver.Actions[i].EntityId = ver.Id
		}
		entityAction := GetEntityActionAdaptor(tx)
		if err = entityAction.AddMultiple(ctx, ver.Actions); err != nil {
			return
		}
	}

	//states
	if len(ver.States) > 0 {
		for i := range ver.States {
			ver.States[i].EntityId = ver.Id
		}
		stateAdaptor := GetEntityStateAdaptor(tx)
		if err = stateAdaptor.AddMultiple(ctx, ver.States); err != nil {
			return
		}
	}

	//metrics
	metricAdaptor := GetMetricAdaptor(tx, nil)
	for _, metric := range ver.Metrics {
		if err = table.AppendMetric(ctx, ver.Id, metricAdaptor.toDb(metric)); err != nil {
			return
		}
	}

	//scripts
	scriptAdaptor := GetScriptAdaptor(tx)
	for _, script := range ver.Scripts {
		if err = table.AppendScript(ctx, ver.Id, scriptAdaptor.toDb(script)); err != nil {
			return
		}
	}

	return
}

// Import ...
func (n *Entity) Import(ctx context.Context, ver *m.Entity) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	// area
	if ver.Area != nil {
		areaAdaptor := GetAreaAdaptor(tx)
		var area *m.Area
		if area, err = areaAdaptor.GetByName(ctx, ver.Area.Name); err != nil {
			if ver.Area.Id, err = areaAdaptor.Add(ctx, ver.Area); err != nil {
				return
			}
		} else {
			ver.Area.Id = area.Id
			ver.AreaId = common.Int64(area.Id)
		}

	}

	// entity
	table := db.Entities{Db: tx}
	if err = table.Add(ctx, n.toDb(ver)); err != nil {
		return
	}

	scriptAdaptor := GetScriptAdaptor(tx)

	//actions
	if len(ver.Actions) > 0 {
		for i, action := range ver.Actions {
			if action.Script != nil {
				var foundedScript *m.Script
				if foundedScript, err = scriptAdaptor.GetByName(action.Script.Name); err == nil {
					action.Script = foundedScript
				} else {
					action.Script.Id = 0
					if action.Script.Id, err = scriptAdaptor.Add(action.Script); err != nil {
						return
					}
				}
			}
			ver.Actions[i].EntityId = ver.Id
		}
		entityAction := GetEntityActionAdaptor(tx)
		if err = entityAction.AddMultiple(ctx, ver.Actions); err != nil {
			return
		}
	}

	//states
	if len(ver.States) > 0 {
		for i := range ver.States {
			ver.States[i].EntityId = ver.Id
		}
		stateAdaptor := GetEntityStateAdaptor(tx)
		if err = stateAdaptor.AddMultiple(ctx, ver.States); err != nil {
			return
		}
	}

	//metrics
	metricAdaptor := GetMetricAdaptor(tx, nil)
	for _, metric := range ver.Metrics {
		if metric.Id, err = metricAdaptor.Add(ctx, metric); err != nil {
			return
		}
		if err = table.AppendMetric(ctx, ver.Id, metricAdaptor.toDb(metric)); err != nil {
			return
		}
	}

	// scripts
	for _, script := range ver.Scripts {
		var foundedScript *m.Script
		if foundedScript, err = scriptAdaptor.GetByName(script.Name); err == nil {
			script = foundedScript
		} else {
			script.Id = 0
			if script.Id, err = scriptAdaptor.Add(script); err != nil {
				return
			}
		}
		if err = table.AppendScript(ctx, ver.Id, scriptAdaptor.toDb(script)); err != nil {
			return
		}
	}
	return
}

// GetById ...
func (n *Entity) GetById(ctx context.Context, id common.EntityId, preloadMetric ...bool) (ver *m.Entity, err error) {

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
func (n *Entity) GetByIds(ctx context.Context, ids []common.EntityId, preloadMetric ...bool) (list []*m.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByIds(ctx, ids); err != nil {
		return
	}
	list = make([]*m.Entity, len(dbList))
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
func (n *Entity) GetByIdsSimple(ctx context.Context, ids []common.EntityId) (list []*m.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByIdsSimple(ctx, ids); err != nil {
		return
	}
	list = make([]*m.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// Delete ...
func (n *Entity) Delete(ctx context.Context, id common.EntityId) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := &db.Entities{Db: tx}
	if err = table.Delete(ctx, id); err != nil {
		return
	}

	return
}

// List ...
func (n *Entity) List(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query, plugin *string, areaId *int64) (list []*m.Entity, total int64, err error) {

	var dbList []*db.Entity
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, autoLoad, query, plugin, areaId); err != nil {
		return
	}

	list = make([]*m.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// GetByType ...
func (n *Entity) GetByType(ctx context.Context, t string, limit, offset int64) (list []*m.Entity, err error) {

	var dbList []*db.Entity
	if dbList, err = n.table.GetByType(ctx, t, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// UpdateSettings ...
func (n *Entity) UpdateSettings(ctx context.Context, entity common.EntityId, settings m.Attributes) (err error) {
	b, _ := json.Marshal(m.EntitySettings{
		Settings: settings,
	})
	err = n.table.UpdateSettings(ctx, entity, b)
	return
}

// Update ...
func (n *Entity) Update(ctx context.Context, ver *m.Entity) (err error) {

	var oldVer *m.Entity
	if oldVer, err = n.GetById(ctx, ver.Id); err != nil {
		return
	}

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.Entities{Db: tx}
	if err = table.Update(ctx, n.toDb(ver)); err != nil {
		return
	}

	entityActionAdaptor := GetEntityActionAdaptor(tx)
	if err = entityActionAdaptor.DeleteByEntityId(ctx, ver.Id); err != nil {
		return
	}

	entityStateAdaptor := GetEntityStateAdaptor(tx)
	if err = entityStateAdaptor.DeleteByEntityId(ctx, ver.Id); err != nil {
		return
	}

	//actions
	for _, action := range ver.Actions {
		action.EntityId = ver.Id
	}
	if err = entityActionAdaptor.AddMultiple(ctx, ver.Actions); err != nil {
		log.Error(err.Error())
		return
	}

	//states
	for _, state := range ver.States {
		state.EntityId = ver.Id
	}
	if err = entityStateAdaptor.AddMultiple(ctx, ver.States); err != nil {
		log.Errorf(err.Error())
		return
	}

	//metrics
	for _, oldMetric := range oldVer.Metrics {
		var exist bool
		for _, metric := range ver.Metrics {
			if metric.Id == oldMetric.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.DeleteMetric(ctx, oldVer.Id, oldMetric.Id); err != nil {
				return
			}
		}
	}

	metricAdaptor := GetMetricAdaptor(tx, nil)
	for _, metric := range ver.Metrics {
		var exist bool
		for _, oldMetric := range oldVer.Metrics {
			if metric.Id == oldMetric.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.AppendMetric(ctx, ver.Id, metricAdaptor.toDb(metric)); err != nil {
				return
			}
		} else {
			if err = metricAdaptor.Update(ctx, metric); err != nil {
				return
			}
		}
	}

	// script
	for _, oldScript := range oldVer.Scripts {
		var exist bool
		for _, script := range ver.Scripts {
			if script.Id == oldScript.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.DeleteScript(ctx, oldVer.Id, oldScript.Id); err != nil {
				return
			}
		}
	}

	scriptAdaptor := GetScriptAdaptor(tx)
	for _, script := range ver.Scripts {
		var exist bool
		for _, oldMetric := range oldVer.Scripts {
			if script.Id == oldMetric.Id {
				exist = true
			}
		}
		if !exist {
			if err = n.table.AppendScript(ctx, ver.Id, scriptAdaptor.toDb(script)); err != nil {
				return
			}
		}
	}

	return
}

// Search ...
func (n *Entity) Search(ctx context.Context, query string, limit, offset int64) (list []*m.Entity, total int64, err error) {
	var dbList []*db.Entity
	if dbList, total, err = n.table.Search(ctx, query, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.Entity, len(dbList))
	for i, dbVer := range dbList {
		list[i] = n.fromDb(dbVer)
	}

	return
}

// AppendMetric ...
func (n *Entity) AppendMetric(ctx context.Context, entityId common.EntityId, metric *m.Metric) (err error) {
	metricAdaptor := GetMetricAdaptor(n.db, nil)
	err = n.table.AppendMetric(ctx, entityId, metricAdaptor.toDb(metric))
	return
}

// DeleteMetric ...
func (n *Entity) DeleteMetric(ctx context.Context, entityId common.EntityId, metric *m.Metric) (err error) {
	err = n.table.DeleteMetric(ctx, entityId, metric.Id)
	return
}

// UpdateAutoload ...
func (n *Entity) UpdateAutoload(ctx context.Context, entityId common.EntityId, autoLoad bool) (err error) {
	err = n.table.UpdateAutoload(ctx, entityId, autoLoad)
	return
}

func (n *Entity) preloadMetric(ctx context.Context, ver *m.Entity) {

	var err error

	// load preview metrics data
	if ver.Metrics == nil || len(ver.Metrics) == 0 {
		return
	}
	bucketMetricBucketAdaptor := GetMetricBucketAdaptor(n.db, nil)
	for i, metric := range ver.Metrics {

		var optionItems = make([]string, len(metric.Options.Items))
		for i, item := range metric.Options.Items {
			optionItems[i] = item.Name
		}

		if ver.Metrics[i].Data, err = bucketMetricBucketAdaptor.SimpleListWithSoftRange(ctx, nil, nil, metric.Id, common.String(common.MetricRange24H.String()), optionItems); err != nil {
			log.Error(err.Error())
			return
		}

		ver.Metrics[i].RangesByType()
	}
}

func (n *Entity) fromDb(dbVer *db.Entity) (ver *m.Entity) {
	ver = &m.Entity{
		Id:          dbVer.Id,
		Description: dbVer.Description,
		PluginName:  dbVer.PluginName,
		Actions:     make([]*m.EntityAction, 0),
		States:      make([]*m.EntityState, 0),
		Icon:        dbVer.Icon,
		AutoLoad:    dbVer.AutoLoad,
		ParentId:    dbVer.ParentId,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
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
		ver.AreaId = common.Int64(dbVer.Area.Id)
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
		ver.Scripts = make([]*m.Script, 0)
	}

	// deserialize payload
	b, _ := dbVer.Payload.MarshalJSON()
	payload := m.EntityPayload{}
	_ = json.Unmarshal(b, &payload)
	ver.Attributes = payload.AttributeSignature

	// deserialize settings
	b, _ = dbVer.Settings.MarshalJSON()
	settings := m.EntitySettings{}
	_ = json.Unmarshal(b, &settings)
	ver.Settings = settings.Settings

	// storage
	if len(dbVer.Storage) > 0 {
		data := map[string]interface{}{}
		_ = json.Unmarshal(dbVer.Storage[0].Attributes, &data)
		_, _ = ver.Attributes.Deserialize(data)
	}

	return
}

func (n *Entity) toDb(ver *m.Entity) (dbVer *db.Entity) {

	dbVer = &db.Entity{
		Id:          ver.Id,
		Description: ver.Description,
		PluginName:  ver.PluginName,
		Icon:        ver.Icon,
		AutoLoad:    ver.AutoLoad,
		ParentId:    ver.ParentId,
		AreaId:      ver.AreaId,
		ImageId:     ver.ImageId,
	}

	// serialize payload
	b, _ := json.Marshal(m.EntityPayload{
		AttributeSignature: ver.Attributes.Signature(),
	})
	_ = dbVer.Payload.UnmarshalJSON(b)

	// serialize settings
	b, _ = json.Marshal(m.EntitySettings{
		Settings: ver.Settings,
	})
	_ = dbVer.Settings.UnmarshalJSON(b)

	// storage
	// ...

	return
}
