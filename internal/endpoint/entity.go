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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"
	"github.com/pkg/errors"
)

// EntityEndpoint ...
type EntityEndpoint struct {
	*CommonEndpoint
}

// NewEntityEndpoint ...
func NewEntityEndpoint(common *CommonEndpoint) *EntityEndpoint {
	return &EntityEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *EntityEndpoint) Add(ctx context.Context, entity *models.Entity) (result *models.Entity, err error) {

	if ok, errs := n.validation.Valid(entity); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		for _, tag := range entity.Tags {
			var foundedTag *models.Tag
			if foundedTag, err = n.adaptors.Tag.GetByName(ctx, tag.Name); err == nil {
				tag.Id = foundedTag.Id
			} else {
				tag.Id = 0
				if tag.Id, err = n.adaptors.Tag.Add(ctx, tag); err != nil {
					return err
				}
			}
		}

		return n.adaptors.Entity.Add(ctx, entity)
	})

	result, err = n.adaptors.Entity.GetById(ctx, entity.Id)
	if err != nil {
		return
	}

	log.Infof("added new entity id:(%s)", result.Id)

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventCreatedEntityModel{
		EntityId: result.Id,
	})

	return
}

// Import ...
func (n *EntityEndpoint) Import(ctx context.Context, entity *models.Entity) (err error) {

	for _, action := range entity.Actions {
		if action.Script != nil {
			var engine scripts.Engine
			if engine, err = n.scriptService.NewEngine(action.Script); err != nil {
				err = errors.Wrap(apperr.ErrInvalidRequest, err.Error())
				return
			}

			if err = engine.Compile(); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}
		}
	}
	for _, script := range entity.Scripts {
		if script != nil {
			var engine scripts.Engine
			if engine, err = n.scriptService.NewEngine(script); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}

			if err = engine.Compile(); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}
		}
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		// area
		if entity.Area != nil {
			var area *models.Area
			if area, err = n.adaptors.Area.GetByName(ctx, entity.Area.Name); err != nil {
				if entity.Area.Id, err = n.adaptors.Area.Add(ctx, entity.Area); err != nil {
					return err
				}
			} else {
				entity.Area.Id = area.Id
				entity.AreaId = pkgCommon.Int64(area.Id)
			}

		}

		// scripts
		for _, script := range entity.Scripts {
			var foundedScript *models.Script
			if foundedScript, err = n.adaptors.Script.GetByName(ctx, script.Name); err == nil {
				script.Id = foundedScript.Id
			} else {
				script.Id = 0
				if script.Id, err = n.adaptors.Script.Add(ctx, script); err != nil {
					return err
				}
			}
		}

		// tags
		for _, tag := range entity.Tags {
			var foundedTag *models.Tag
			if foundedTag, err = n.adaptors.Tag.GetByName(ctx, tag.Name); err == nil {
				tag.Id = foundedTag.Id
			} else {
				tag.Id = 0
				if tag.Id, err = n.adaptors.Tag.Add(ctx, tag); err != nil {
					return err
				}
			}
		}

		//actions
		if len(entity.Actions) > 0 {
			for i, action := range entity.Actions {
				action.Id = 0
				action.EntityId = entity.Id
				if action.Script != nil {
					var foundedScript *models.Script
					if foundedScript, err = n.adaptors.Script.GetByName(ctx, action.Script.Name); err == nil {
						action.Script.Id = foundedScript.Id
					} else {
						action.Script.Id = 0
						if action.Script.Id, err = n.adaptors.Script.Add(ctx, action.Script); err != nil {
							return err
						}
					}
				}
				if entity.Actions[i].Icon != nil && *entity.Actions[i].Icon == "" {
					entity.Actions[i].Icon = nil
				}
				entity.Actions[i].EntityId = entity.Id
			}
		}

		//states
		if len(entity.States) > 0 {
			for _, state := range entity.States {
				state.Id = 0
				state.EntityId = entity.Id
			}
		}

		//metrics
		for _, metric := range entity.Metrics {
			metric.Id = 0
			if metric.Id, err = n.adaptors.Metric.Add(ctx, metric); err != nil {
				return err
			}
		}

		// entity
		if err = n.adaptors.Entity.Add(ctx, entity); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventCreatedEntityModel{
		EntityId: entity.Id,
	})

	log.Infof("entity id:(%s) was imported", entity.Id)

	return
}

// GetById ...
func (n *EntityEndpoint) GetById(ctx context.Context, id pkgCommon.EntityId) (result *models.Entity, err error) {
	if result, err = n.adaptors.Entity.GetById(ctx, id); err != nil {
		return
	}
	result.IsLoaded = n.supervisor.EntityIsLoaded(id)
	return
}

// Update ...
func (n *EntityEndpoint) Update(ctx context.Context, params *models.Entity) (result *models.Entity, err error) {

	var entity *models.Entity
	if entity, err = n.adaptors.Entity.GetById(ctx, params.Id); err != nil {
		return
	}

	entity.Description = params.Description
	entity.PluginName = params.PluginName
	entity.Icon = params.Icon
	entity.ImageId = params.ImageId
	entity.Actions = params.Actions
	entity.States = params.States
	entity.AreaId = params.AreaId
	entity.Metrics = params.Metrics
	entity.Scripts = params.Scripts
	entity.Tags = params.Tags
	entity.Hidden = params.Hidden
	entity.Attributes = params.Attributes
	entity.Settings = params.Settings
	entity.ParentId = params.ParentId
	entity.RestoreState = params.RestoreState
	entity.AutoLoad = params.AutoLoad

	if ok, errs := n.validation.Valid(entity); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var oldVer *models.Entity
	if oldVer, err = n.adaptors.Entity.GetById(ctx, params.Id); err != nil {
		return
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		// entity action
		if err = n.adaptors.EntityAction.DeleteByEntityId(ctx, entity.Id); err != nil {
			return err
		}

		// entity action
		if err = n.adaptors.EntityState.DeleteByEntityId(ctx, entity.Id); err != nil {
			return err
		}

		// scripts
		if err = n.adaptors.Entity.DeleteScripts(ctx, oldVer.Id); err != nil {
			return err
		}

		// tags
		if err = n.adaptors.Entity.DeleteTags(ctx, oldVer.Id); err != nil {
			return err
		}

		// tags
		for _, tag := range entity.Tags {
			var foundedTag *models.Tag
			if foundedTag, err = n.adaptors.Tag.GetByName(ctx, tag.Name); err == nil {
				tag.Id = foundedTag.Id
			} else {
				tag.Id = 0
				if tag.Id, err = n.adaptors.Tag.Add(ctx, tag); err != nil {
					return err
				}
			}
		}

		if err = n.adaptors.Entity.Update(ctx, entity); err != nil {
			return err
		}

		//metrics
		for _, oldMetric := range oldVer.Metrics {
			var exist bool
			for _, metric := range entity.Metrics {
				if metric.Id == oldMetric.Id {
					exist = true
				}
			}
			if !exist {
				if err = n.adaptors.Metric.Delete(ctx, oldMetric.Id); err != nil {
					return err
				}
			}
		}

		for _, metric := range entity.Metrics {
			var exist bool
			for _, oldMetric := range oldVer.Metrics {
				if metric.Id == oldMetric.Id {
					exist = true
				}
			}
			if exist {
				if err = n.adaptors.Metric.Update(ctx, metric); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	result, err = n.adaptors.Entity.GetById(ctx, entity.Id)
	if err != nil {
		return
	}

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventUpdatedEntityModel{
		EntityId: result.Id,
	})

	log.Infof("updated entity id:(%s)", result.Id)

	return
}

// List ...
func (n *EntityEndpoint) List(ctx context.Context, pagination common.PageParams, query, plugin *string, areaId *int64, tags *[]string) (entities []*models.Entity, total int64, err error) {
	entities, total, err = n.adaptors.Entity.ListPlain(ctx, pagination.Limit, pagination.Offset, pagination.Order,
		pagination.SortBy, false, query, plugin, areaId, tags)
	if err != nil {
		return
	}
	for _, entity := range entities {
		entity.IsLoaded = n.supervisor.EntityIsLoaded(entity.Id)
	}
	return
}

// Delete ...
func (n *EntityEndpoint) Delete(ctx context.Context, id pkgCommon.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {
		var entity *models.Entity
		if entity, err = n.adaptors.Entity.GetById(ctx, id); err != nil {
			return err
		}

		for _, metric := range entity.Metrics {
			if err = n.adaptors.Metric.Delete(ctx, metric.Id); err != nil {
				return err
			}
		}
		return n.adaptors.Entity.Delete(ctx, entity.Id)
	})

	n.eventBus.Publish("system/models/entities/"+id.String(), events.CommandUnloadEntity{
		EntityId: id,
	})

	log.Infof("entity id:(%s) was deleted", id)

	return
}

// Search ...
func (n *EntityEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*models.Entity, total int64, err error) {

	result, total, err = n.adaptors.Entity.Search(ctx, query, limit, offset)
	return
}

// Enable ...
func (n *EntityEndpoint) Enable(ctx context.Context, id pkgCommon.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *models.Entity
	entity, err = n.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Entity.UpdateAutoload(ctx, entity.Id, true); err != nil {
		return
	}

	n.eventBus.Publish("system/entities/"+id.String(), events.CommandLoadEntity{
		EntityId: id,
	})

	log.Infof("entity id:(%s) was enabled", id)

	return
}

// Disable ...
func (n *EntityEndpoint) Disable(ctx context.Context, id pkgCommon.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *models.Entity
	entity, err = n.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Entity.UpdateAutoload(ctx, entity.Id, false); err != nil {
		return
	}

	n.eventBus.Publish("system/entities/"+id.String(), events.CommandUnloadEntity{
		EntityId: id,
	})

	log.Infof("entity id:(%s) was disabled", id)

	return
}

// Statistic ...
func (n *EntityEndpoint) Statistic(ctx context.Context) (statistic []*models.Statistic, err error) {
	var stat *models.EntitiesStatistic
	if stat, err = n.adaptors.Entity.Statistic(ctx); err != nil {
		return
	}
	statistic = []*models.Statistic{
		{
			Name:        "entities.stat_total_name",
			Description: "entities.stat_total_descr",
			Value:       stat.Total,
			Diff:        0,
		},
		{
			Name:        "entities.stat_used_name",
			Description: "entities.stat_used_descr",
			Value:       stat.Used,
			Diff:        0,
		},
		{
			Name:        "entities.stat_unused_name",
			Description: "entities.stat_unused_descr",
			Value:       stat.Unused,
			Diff:        0,
		},
	}
	return
}
