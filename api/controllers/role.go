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
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/stub"
)

// ControllerRole ...
type ControllerRole struct {
	*ControllerCommon
}

// NewControllerRole ...
func NewControllerRole(common *ControllerCommon) *ControllerRole {
	return &ControllerRole{
		ControllerCommon: common,
	}
}

// GetRoleAccessList ...
func (c ControllerRole) RoleServiceGetRoleAccessList(ctx echo.Context, name string) error {

	accessList, err := c.endpoint.Role.GetAccessList(ctx.Request().Context(), name, c.accessList)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Role.ToRoleAccessListResult(accessList)))
}

// UpdateRoleAccessList ...
func (c ControllerRole) RoleServiceUpdateRoleAccessList(ctx echo.Context, name string, _ stub.RoleServiceUpdateRoleAccessListParams) error {

	obj := &stub.RoleServiceUpdateRoleAccessListJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	if err := c.endpoint.Role.UpdateAccessList(ctx.Request().Context(), name, c.dto.Role.UpdateRoleAccessList(obj, name)); err != nil {
		return c.ERROR(ctx, err)
	}

	accessList, err := c.endpoint.Role.GetAccessList(ctx.Request().Context(), name, c.accessList)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Role.ToRoleAccessListResult(accessList)))
}

// AddRole ...
func (c ControllerRole) RoleServiceAddRole(ctx echo.Context, _ stub.RoleServiceAddRoleParams) error {

	obj := &stub.ApiNewRoleRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	role, err := c.endpoint.Role.Add(ctx.Request().Context(), c.dto.Role.FromNewRoleRequest(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Role.GetStubRole(role)))
}

// GetRoleByName ...
func (c ControllerRole) RoleServiceGetRoleByName(ctx echo.Context, name string) error {

	role, err := c.endpoint.Role.GetByName(ctx.Request().Context(), name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Role.GetStubRole(role)))
}

// UpdateRoleByName ...
func (c ControllerRole) RoleServiceUpdateRoleByName(ctx echo.Context, name string, _ stub.RoleServiceUpdateRoleByNameParams) error {

	obj := &stub.RoleServiceUpdateRoleByNameJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	role, err := c.endpoint.Role.Update(ctx.Request().Context(), c.dto.Role.FromUpdateRoleRequest(obj, name))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Role.GetStubRole(role)))
}

// GetRoleList ...
func (c ControllerRole) RoleServiceGetRoleList(ctx echo.Context, params stub.RoleServiceGetRoleListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Role.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Role.ToListResult(items), total, pagination))
}

// SearchRoleByName ...
func (c ControllerRole) RoleServiceSearchRoleByName(ctx echo.Context, params stub.RoleServiceSearchRoleByNameParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Role.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Role.ToSearchResult(items))
}

// DeleteRoleByName ...
func (c ControllerRole) RoleServiceDeleteRoleByName(ctx echo.Context, name string) error {

	if err := c.endpoint.Role.Delete(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
