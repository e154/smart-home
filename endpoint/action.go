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
	"fmt"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
)

// ActionEndpoint ...
type ActionEndpoint struct {
	*CommonEndpoint
}

// NewActionEndpoint ...
func NewActionEndpoint(common *CommonEndpoint) *ActionEndpoint {
	return &ActionEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *ActionEndpoint) Add(ctx context.Context, action *m.Action) (result *m.Action, err error) {

	if ok, errs := n.validation.Valid(action); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if action.Id, err = n.adaptors.Action.Add(ctx, action); err != nil {
		return
	}

	if result, err = n.adaptors.Action.GetById(ctx, action.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/actions/%d", result.Id), events.EventAddedActionModel{
		Id: action.Id,
	})

	return
}

// GetById ...
func (n *ActionEndpoint) GetById(ctx context.Context, id int64) (result *m.Action, err error) {

	result, err = n.adaptors.Action.GetById(ctx, id)

	return
}

// Update ...
func (n *ActionEndpoint) Update(ctx context.Context, params *m.Action) (action *m.Action, err error) {

	_, err = n.adaptors.Action.GetById(ctx, params.Id)
	if err != nil {
		return
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Action.Update(ctx, params); err != nil {
		return
	}

	if action, err = n.adaptors.Action.GetById(ctx, params.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/actions/%d", action.Id), events.EventUpdatedActionModel{
		Id:     action.Id,
		Action: action,
	})

	return
}

// GetList ...
func (n *ActionEndpoint) GetList(ctx context.Context, pagination common.PageParams, ids *[]uint64) (result []*m.Action, total int64, err error) {

	result, total, err = n.adaptors.Action.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, ids)

	return
}

// Delete ...
func (n *ActionEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = n.adaptors.Action.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Action.Delete(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/actions/%d", id), events.EventRemovedActionModel{
		Id: id,
	})

	return
}

// Search ...
func (n *ActionEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Action, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Action.Search(ctx, query, int(limit), int(offset))

	return
}
