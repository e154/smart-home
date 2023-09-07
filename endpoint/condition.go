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
	"fmt"
	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
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
func (n *ConditionEndpoint) Add(ctx context.Context, condition *m.Condition) (result *m.Condition, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(condition); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetContext(err, errs)
		return
	}

	if condition.Id, err = n.adaptors.Condition.Add(condition); err != nil {
		return
	}

	if result, err = n.adaptors.Condition.GetById(condition.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/conditions/%d", result.Id), events.EventAddedCondition{
		Id: result.Id,
	})

	return
}

// GetById ...
func (n *ConditionEndpoint) GetById(ctx context.Context, id int64) (result *m.Condition, err error) {

	result, err = n.adaptors.Condition.GetById(id)

	return
}

// Update ...
func (n *ConditionEndpoint) Update(ctx context.Context, params *m.Condition) (result *m.Condition, errs validator.ValidationErrorsTranslations, err error) {

	_, err = n.adaptors.Condition.GetById(params.Id)
	if err != nil {
		return
	}

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if err = n.adaptors.Condition.Update(params); err != nil {
		return
	}

	if result, err = n.adaptors.Condition.GetById(params.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/conditions/%d", result.Id), events.EventUpdatedCondition{
		Id: result.Id,
	})

	return
}

// GetList ...
func (n *ConditionEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Condition, total int64, err error) {

	result, total, err = n.adaptors.Condition.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// Delete ...
func (n *ConditionEndpoint) Delete(ctx context.Context, id int64) (err error) {

	var area *m.Condition
	area, err = n.adaptors.Condition.GetById(id)
	if err != nil {
		return
	}

	if err = n.adaptors.Condition.Delete(area.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/conditions/%d", id), events.EventRemovedCondition{
		Id: id,
	})

	return
}

// Search ...
func (n *ConditionEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Condition, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Condition.Search(query, int(limit), int(offset))

	return
}
