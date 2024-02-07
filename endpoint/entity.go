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
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
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
func (n *EntityEndpoint) Add(ctx context.Context, entity *m.Entity) (result *m.Entity, err error) {

	if ok, errs := n.validation.Valid(entity); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Entity.Add(ctx, entity); err != nil {
		return
	}

	result, err = n.adaptors.Entity.GetById(ctx, entity.Id)
	if err != nil {
		return
	}

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventCreatedEntityModel{
		EntityId: result.Id,
	})

	return
}

// Import ...
func (n *EntityEndpoint) Import(ctx context.Context, entity *m.Entity) (err error) {

	for _, action := range entity.Actions {
		if action.Script != nil {
			var engine *scripts.Engine
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
			var engine *scripts.Engine
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

	if err = n.adaptors.Entity.Import(ctx, entity); err != nil {
		return
	}

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventCreatedEntityModel{
		EntityId: entity.Id,
	})

	return
}

// GetById ...
func (n *EntityEndpoint) GetById(ctx context.Context, id common.EntityId) (result *m.Entity, err error) {
	if result, err = n.adaptors.Entity.GetById(ctx, id); err != nil {
		return
	}
	result.IsLoaded = n.supervisor.EntityIsLoaded(id)
	return
}

// Update ...
func (n *EntityEndpoint) Update(ctx context.Context, params *m.Entity) (result *m.Entity, err error) {

	var entity *m.Entity
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
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Entity.Update(ctx, entity); err != nil {
		return
	}

	result, err = n.adaptors.Entity.GetById(ctx, entity.Id)
	if err != nil {
		return
	}

	n.eventBus.Publish("system/models/entities/"+entity.Id.String(), events.EventUpdatedEntityModel{
		EntityId: result.Id,
	})

	return
}

// List ...
func (n *EntityEndpoint) List(ctx context.Context, pagination common.PageParams, query, plugin *string, areaId *int64, tags *[]string) (entities []*m.Entity, total int64, err error) {
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
func (n *EntityEndpoint) Delete(ctx context.Context, id common.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *m.Entity
	entity, err = n.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Entity.Delete(ctx, entity.Id); err != nil {
		return
	}

	n.eventBus.Publish("system/models/entities/"+id.String(), events.CommandUnloadEntity{
		EntityId: id,
	})

	return
}

// Search ...
func (n *EntityEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Entity, total int64, err error) {

	result, total, err = n.adaptors.Entity.Search(ctx, query, limit, offset)
	return
}

// Enable ...
func (n *EntityEndpoint) Enable(ctx context.Context, id common.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *m.Entity
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

	return
}

// Disable ...
func (n *EntityEndpoint) Disable(ctx context.Context, id common.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *m.Entity
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

	return
}
