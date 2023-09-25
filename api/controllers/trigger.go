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

// ControllerTrigger ...
type ControllerTrigger struct {
	*ControllerCommon
}

// NewControllerTrigger ...
func NewControllerTrigger(common *ControllerCommon) ControllerTrigger {
	return ControllerTrigger{
		ControllerCommon: common,
	}
}

// AddTrigger ...
func (c ControllerTrigger) AddTrigger(ctx context.Context, req *api.NewTriggerRequest) (*api.Trigger, error) {

	trigger := c.dto.Trigger.AddTrigger(req)

	trigger, errs, err := c.endpoint.Trigger.Add(ctx, trigger)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Trigger.ToTrigger(trigger), nil
}

// UpdateTrigger ...
func (c ControllerTrigger) UpdateTrigger(ctx context.Context, req *api.UpdateTriggerRequest) (*api.Trigger, error) {

	trigger := c.dto.Trigger.UpdateTrigger(req)

	trigger, errs, err := c.endpoint.Trigger.Update(ctx, trigger)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Trigger.ToTrigger(trigger), nil
}

// GetTriggerById ...
func (c ControllerTrigger) GetTriggerById(ctx context.Context, req *api.GetTriggerRequest) (*api.Trigger, error) {

	trigger, err := c.endpoint.Trigger.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Trigger.ToTrigger(trigger), nil
}

// GetTriggerList ...
func (c ControllerTrigger) GetTriggerList(ctx context.Context, req *api.PaginationRequest) (*api.GetTriggerListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Trigger.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Trigger.ToListResult(items, uint64(total), pagination), nil
}

// DeleteTrigger ...
func (c ControllerTrigger) DeleteTrigger(ctx context.Context, req *api.DeleteTriggerRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Trigger.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// SearchTrigger ...
func (c ControllerTrigger) SearchTrigger(ctx context.Context, req *api.SearchRequest) (*api.SearchTriggerResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Trigger.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Trigger.ToSearchResult(items), nil
}

// EnableTrigger ...
func (c ControllerTrigger) EnableTrigger(ctx context.Context, req *api.EnableTriggerRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Trigger.Enable(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// DisableTrigger ...
func (c ControllerTrigger) DisableTrigger(ctx context.Context, req *api.DisableTriggerRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Trigger.Disable(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
