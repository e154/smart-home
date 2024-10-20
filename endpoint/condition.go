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

// ConditionEndpoint ...
type ConditionEndpoint struct {
	*CommonEndpoint
}

// NewConditionEndpoint ...
func NewConditionEndpoint(common *CommonEndpoint) *ConditionEndpoint {
	return &ConditionEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *ConditionEndpoint) Add(ctx context.Context, condition *m.Condition) (result *m.Condition, err error) {

	if ok, errs := n.validation.Valid(condition); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if condition.Id, err = n.adaptors.Condition.Add(ctx, condition); err != nil {
		return
	}

	if result, err = n.adaptors.Condition.GetById(ctx, condition.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/conditions/%d", result.Id), events.EventAddedConditionModel{
		Id: result.Id,
	})

	log.Infof("added new condition %s id:(%d)", result.Name, result.Id)

	return
}

// GetById ...
func (n *ConditionEndpoint) GetById(ctx context.Context, id int64) (result *m.Condition, err error) {

	result, err = n.adaptors.Condition.GetById(ctx, id)

	return
}

// Update ...
func (n *ConditionEndpoint) Update(ctx context.Context, params *m.Condition) (result *m.Condition, err error) {

	_, err = n.adaptors.Condition.GetById(ctx, params.Id)
	if err != nil {
		return
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Condition.Update(ctx, params); err != nil {
		return
	}

	if result, err = n.adaptors.Condition.GetById(ctx, params.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/conditions/%d", result.Id), events.EventUpdatedConditionModel{
		Id: result.Id,
	})

	log.Infof("updated condition %s id:(%d)", result.Name, result.Id)

	return
}

// GetList ...
func (n *ConditionEndpoint) GetList(ctx context.Context, pagination common.PageParams, ids *[]uint64) (result []*m.Condition, total int64, err error) {

	result, total, err = n.adaptors.Condition.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, ids)

	return
}

// Delete ...
func (n *ConditionEndpoint) Delete(ctx context.Context, id int64) (err error) {

	var condition *m.Condition
	condition, err = n.adaptors.Condition.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Condition.Delete(ctx, condition.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/conditions/%d", id), events.EventRemovedConditionModel{
		Id: id,
	})

	log.Infof("condition %s id:(%d) was deleted", condition.Name, id)

	return
}

// Search ...
func (n *ConditionEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Condition, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Condition.Search(ctx, query, int(limit), int(offset))

	return
}
