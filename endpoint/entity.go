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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/go-playground/validator/v10"
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
func (n *EntityEndpoint) Add(ctx context.Context, entity *m.Entity) (result *m.Entity, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(entity); !ok {
		return
	}

	if err = n.adaptors.Entity.Add(entity); err != nil {
		return
	}

	result, err = n.adaptors.Entity.GetById(entity.Id)
	if err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicEntities, events.EventCreatedEntity{
		Id: result.Id,
	})

	return
}

// GetById ...
func (n *EntityEndpoint) GetById(ctx context.Context, id common.EntityId) (result *m.Entity, err error) {

	result, err = n.adaptors.Entity.GetById(id)

	return
}

// Update ...
func (n *EntityEndpoint) Update(ctx context.Context, params *m.Entity) (result *m.Entity, errs validator.ValidationErrorsTranslations, err error) {

	var entity *m.Entity
	if entity, err = n.adaptors.Entity.GetById(params.Id); err != nil {
		return
	}

	_ = common.Copy(entity, params, common.JsonEngine)
	entity.Settings = params.Settings
	entity.Attributes = params.Attributes

	var ok bool
	if ok, errs = n.validation.Valid(entity); !ok {
		return
	}

	if err = n.adaptors.Entity.Update(entity); err != nil {
		return
	}

	result, err = n.adaptors.Entity.GetById(entity.Id)
	if err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicEntities, events.EventUpdatedEntity{
		Id: result.Id,
	})

	return
}

// List ...
func (n *EntityEndpoint) List(ctx context.Context, pagination common.PageParams) (result []*m.Entity, total int64, err error) {
	result, total, err = n.adaptors.Entity.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false)

	return
}

// Delete ...
func (n *EntityEndpoint) Delete(ctx context.Context, id common.EntityId) (err error) {

	if id == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var entity *m.Entity
	entity, err = n.adaptors.Entity.GetById(id)
	if err != nil {
		return
	}

	if err = n.adaptors.Entity.Delete(entity.Id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicEntities, events.EventDeletedEntity{
		Id: id,
	})

	return
}

// Search ...
func (n *EntityEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Entity, total int64, err error) {

	result, total, err = n.adaptors.Entity.Search(query, limit, offset)
	return
}
