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
	"github.com/e154/smart-home/pkg/common"
	"github.com/labstack/echo/v4"
)

// ControllerDeveloperTools ...
type ControllerDeveloperTools struct {
	*ControllerCommon
}

// NewControllerDeveloperTools ...
func NewControllerDeveloperTools(common *ControllerCommon) *ControllerDeveloperTools {
	return &ControllerDeveloperTools{
		ControllerCommon: common,
	}
}

// EntitySetStateName ...
func (c ControllerDeveloperTools) DeveloperToolsServiceEntitySetState(ctx echo.Context, _ stub.DeveloperToolsServiceEntitySetStateParams) error {

	obj := &stub.ApiEntityRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.DeveloperTools.EntitySetStateName(ctx.Request().Context(), common.EntityId(obj.Id), obj.Name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ReloadEntity ...
func (c ControllerDeveloperTools) DeveloperToolsServiceReloadEntity(ctx echo.Context, _ stub.DeveloperToolsServiceReloadEntityParams) error {

	obj := &stub.ApiReloadRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.DeveloperTools.ReloadEntity(ctx.Request().Context(), common.EntityId(obj.Id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// CallTrigger ...
func (c ControllerDeveloperTools) DeveloperToolsServiceCallTrigger(ctx echo.Context, _ stub.DeveloperToolsServiceCallTriggerParams) error {

	obj := &stub.ApiAutomationRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	if err := c.endpoint.DeveloperTools.TaskCallTrigger(ctx.Request().Context(), obj.Id, obj.Name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// CallAction ...
func (c ControllerDeveloperTools) DeveloperToolsServiceCallAction(ctx echo.Context, _ stub.DeveloperToolsServiceCallActionParams) error {

	obj := &stub.ApiAutomationRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	if err := c.endpoint.DeveloperTools.TaskCallAction(ctx.Request().Context(), obj.Id, obj.Name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// GetEventBusStateList ...
func (c ControllerDeveloperTools) DeveloperToolsServiceGetEventBusStateList(ctx echo.Context, params stub.DeveloperToolsServiceGetEventBusStateListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.DeveloperTools.GetEventBusState(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.DeveloperTools.ToListResult(items), total, pagination))
}
