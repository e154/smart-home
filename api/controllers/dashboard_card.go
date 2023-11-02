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
	"github.com/e154/smart-home/api/stub"
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/dto"
)

// ControllerDashboardCard ...
type ControllerDashboardCard struct {
	*ControllerCommon
}

// NewControllerDashboardCard ...
func NewControllerDashboardCard(common *ControllerCommon) *ControllerDashboardCard {
	return &ControllerDashboardCard{
		ControllerCommon: common,
	}
}

// AddDashboardCard ...
func (c ControllerDashboardCard) DashboardCardServiceAddDashboardCard(ctx echo.Context, _ stub.DashboardCardServiceAddDashboardCardParams) error {

	obj := &stub.ApiNewDashboardCardRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	card, err := c.endpoint.DashboardCard.Add(ctx.Request().Context(), c.dto.DashboardCard.AddDashboardCard(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, dto.ToDashboardCard(card)))
}

// UpdateDashboardCard ...
func (c ControllerDashboardCard) DashboardCardServiceUpdateDashboardCard(ctx echo.Context, id int64, _ stub.DashboardCardServiceUpdateDashboardCardParams) error {

	obj := &stub.DashboardCardServiceUpdateDashboardCardJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	card, err := c.endpoint.DashboardCard.Update(ctx.Request().Context(), c.dto.DashboardCard.UpdateDashboardCard(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToDashboardCard(card)))
}

// GetDashboardCardById ...
func (c ControllerDashboardCard) DashboardCardServiceGetDashboardCardById(ctx echo.Context, id int64) error {

	card, err := c.endpoint.DashboardCard.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToDashboardCard(card)))
}

// GetDashboardCardList ...
func (c ControllerDashboardCard) DashboardCardServiceGetDashboardCardList(ctx echo.Context, params stub.DashboardCardServiceGetDashboardCardListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.DashboardCard.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.DashboardCard.ToListResult(items), total, pagination))
}

// DeleteDashboardCard ...
func (c ControllerDashboardCard) DashboardCardServiceDeleteDashboardCard(ctx echo.Context, id int64) error {

	if err := c.endpoint.DashboardCard.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ImportDashboardCard ...
func (c ControllerDashboardCard) DashboardCardServiceImportDashboardCard(ctx echo.Context, _ stub.DashboardCardServiceImportDashboardCardParams) error {

	obj := &stub.ApiDashboardCard{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	card, err := c.endpoint.DashboardCard.Import(ctx.Request().Context(), dto.ImportDashboardCard(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.ToDashboardCard(card)))
}
