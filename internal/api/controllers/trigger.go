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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerTrigger ...
type ControllerTrigger struct {
	*ControllerCommon
}

// NewControllerTrigger ...
func NewControllerTrigger(common *ControllerCommon) *ControllerTrigger {
	return &ControllerTrigger{
		ControllerCommon: common,
	}
}

// AddTrigger ...
func (c ControllerTrigger) TriggerServiceAddTrigger(ctx echo.Context, _ stub.TriggerServiceAddTriggerParams) error {

	obj := &stub.ApiNewTriggerRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	trigger, err := c.endpoint.Trigger.Add(ctx.Request().Context(), c.dto.Trigger.AddTrigger(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Trigger.ToTrigger(trigger)))
}

// UpdateTrigger ...
func (c ControllerTrigger) TriggerServiceUpdateTrigger(ctx echo.Context, id int64, _ stub.TriggerServiceUpdateTriggerParams) error {

	obj := &stub.TriggerServiceUpdateTriggerJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	trigger, err := c.endpoint.Trigger.Update(ctx.Request().Context(), c.dto.Trigger.UpdateTrigger(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Trigger.ToTrigger(trigger)))
}

// GetTriggerById ...
func (c ControllerTrigger) TriggerServiceGetTriggerById(ctx echo.Context, id int64) error {

	trigger, err := c.endpoint.Trigger.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Trigger.ToTrigger(trigger)))
}

// GetTriggerList ...
func (c ControllerTrigger) TriggerServiceGetTriggerList(ctx echo.Context, params stub.TriggerServiceGetTriggerListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Trigger.GetList(ctx.Request().Context(), pagination, params.Ids)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Trigger.ToListResult(items), total, pagination))
}

// DeleteTrigger ...
func (c ControllerTrigger) TriggerServiceDeleteTrigger(ctx echo.Context, id int64) error {

	if err := c.endpoint.Trigger.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchTrigger ...
func (c ControllerTrigger) TriggerServiceSearchTrigger(ctx echo.Context, params stub.TriggerServiceSearchTriggerParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Trigger.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Trigger.ToSearchResult(items))
}

// EnableTrigger ...
func (c ControllerTrigger) TriggerServiceEnableTrigger(ctx echo.Context, id int64) error {

	if err := c.endpoint.Trigger.Enable(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DisableTrigger ...
func (c ControllerTrigger) TriggerServiceDisableTrigger(ctx echo.Context, id int64) error {

	if err := c.endpoint.Trigger.Disable(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
