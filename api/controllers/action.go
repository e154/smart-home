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

// ControllerAction ...
type ControllerAction struct {
	*ControllerCommon
}

// NewControllerAction ...
func NewControllerAction(common *ControllerCommon) ControllerAction {
	return ControllerAction{
		ControllerCommon: common,
	}
}

// AddAction ...
func (c ControllerAction) AddAction(ctx context.Context, req *api.NewActionRequest) (*api.Action, error) {

	action := c.dto.Action.AddAction(req)

	action, errs, err := c.endpoint.Action.Add(ctx, action)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Action.ToAction(action), nil
}

// UpdateAction ...
func (c ControllerAction) UpdateAction(ctx context.Context, req *api.UpdateActionRequest) (*api.Action, error) {

	action := c.dto.Action.UpdateAction(req)

	action, errs, err := c.endpoint.Action.Update(ctx, action)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Action.ToAction(action), nil
}

// GetActionById ...
func (c ControllerAction) GetActionById(ctx context.Context, req *api.GetActionRequest) (*api.Action, error) {

	action, err := c.endpoint.Action.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Action.ToAction(action), nil
}

// GetActionList ...
func (c ControllerAction) GetActionList(ctx context.Context, req *api.PaginationRequest) (*api.GetActionListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Action.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Action.ToListResult(items, uint64(total), pagination), nil
}

// DeleteAction ...
func (c ControllerAction) DeleteAction(ctx context.Context, req *api.DeleteActionRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Action.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// SearchAction ...
func (c ControllerAction) SearchAction(ctx context.Context, req *api.SearchRequest) (*api.SearchActionResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Action.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Action.ToSearchResult(items), nil
}
