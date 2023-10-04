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

package controllers

import (
	"context"

	"github.com/e154/smart-home/api/stub/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerCondition ...
type ControllerCondition struct {
	*ControllerCommon
}

// NewControllerCondition ...
func NewControllerCondition(common *ControllerCommon) ControllerCondition {
	return ControllerCondition{
		ControllerCommon: common,
	}
}

// AddCondition ...
func (c ControllerCondition) AddCondition(ctx context.Context, req *api.NewConditionRequest) (*api.Condition, error) {

	condition := c.dto.Condition.AddCondition(req)

	condition, errs, err := c.endpoint.Condition.Add(ctx, condition)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Condition.ToCondition(condition), nil
}

// UpdateCondition ...
func (c ControllerCondition) UpdateCondition(ctx context.Context, req *api.UpdateConditionRequest) (*api.Condition, error) {

	condition := c.dto.Condition.UpdateCondition(req)

	condition, errs, err := c.endpoint.Condition.Update(ctx, condition)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Condition.ToCondition(condition), nil
}

// GetConditionById ...
func (c ControllerCondition) GetConditionById(ctx context.Context, req *api.GetConditionRequest) (*api.Condition, error) {

	condition, err := c.endpoint.Condition.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Condition.ToCondition(condition), nil
}

// GetConditionList ...
func (c ControllerCondition) GetConditionList(ctx context.Context, req *api.PaginationRequest) (*api.GetConditionListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Condition.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Condition.ToListResult(items, uint64(total), pagination), nil
}

// DeleteCondition ...
func (c ControllerCondition) DeleteCondition(ctx context.Context, req *api.DeleteConditionRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Condition.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// SearchCondition ...
func (c ControllerCondition) SearchCondition(ctx context.Context, req *api.SearchRequest) (*api.SearchConditionResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Condition.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Condition.ToSearchResult(items), nil
}
