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

package controllers

import (
	"context"

	"github.com/e154/smart-home/api/stub/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerRole ...
type ControllerRole struct {
	*ControllerCommon
}

// NewControllerRole ...
func NewControllerRole(common *ControllerCommon) ControllerRole {
	return ControllerRole{
		ControllerCommon: common,
	}
}

// GetRoleAccessList ...
func (c ControllerRole) GetRoleAccessList(ctx context.Context, req *api.GetRoleAccessListRequest) (*api.RoleAccessListResult, error) {

	accessList, err := c.endpoint.Role.GetAccessList(ctx, req.Name, c.accessList)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Role.ToRoleAccessListResult(accessList), nil
}

// UpdateRoleAccessList ...
func (c ControllerRole) UpdateRoleAccessList(ctx context.Context, req *api.UpdateRoleAccessListRequest) (*api.RoleAccessListResult, error) {

	if err := c.endpoint.Role.UpdateAccessList(ctx, req.Name, c.dto.Role.FromUpdateRoleAccessListRequest(req)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	accessList, err := c.endpoint.Role.GetAccessList(ctx, req.Name, c.accessList)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Role.ToRoleAccessListResult(accessList), nil
}

// AddRole ...
func (c ControllerRole) AddRole(ctx context.Context, req *api.NewRoleRequest) (*api.Role, error) {

	role := c.dto.Role.FromNewRoleRequest(req)

	role, errs, err := c.endpoint.Role.Add(ctx, role)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Role.ToGRole(role), nil
}

// GetRoleByName ...
func (c ControllerRole) GetRoleByName(ctx context.Context, req *api.GetRoleRequest) (*api.Role, error) {

	role, err := c.endpoint.Role.GetByName(ctx, req.Name)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Role.ToGRole(role), nil
}

// UpdateRoleByName ...
func (c ControllerRole) UpdateRoleByName(ctx context.Context, req *api.UpdateRoleRequest) (*api.Role, error) {

	role := c.dto.Role.FromUpdateRoleRequest(req)

	role, errs, err := c.endpoint.Role.Update(ctx, role)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Role.ToGRole(role), nil
}

// GetRoleList ...
func (c ControllerRole) GetRoleList(ctx context.Context, req *api.GetRoleListRequest) (*api.GetRoleListResult, error) {

	pagination := c.Pagination(req.Limit, req.Offset, req.Order, req.SortBy)
	items, total, err := c.endpoint.Role.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Role.ToListResult(items, uint64(total), pagination), nil
}

// SearchRoleByName ...
func (c ControllerRole) SearchRoleByName(ctx context.Context, req *api.SearchRoleRequest) (*api.SearchRoleListResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Role.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Role.ToSearchResult(items), nil
}

// DeleteRoleByName ...
func (c ControllerRole) DeleteRoleByName(ctx context.Context, req *api.DeleteRoleRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Role.Delete(ctx, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
