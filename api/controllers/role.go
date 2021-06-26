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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ControllerRole struct {
	*ControllerCommon
}

func NewControllerRole(common *ControllerCommon) ControllerRole {
	return ControllerRole{
		ControllerCommon: common,
	}
}

func (c ControllerRole) GetRoleAccessList(_ context.Context, req *api.GetRoleAccessListRequest) (*api.RoleAccessListResult, error) {

	accessList, err := c.endpoint.Role.GetAccessList(req.Name, c.accessList)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToRoleAccessListResult(accessList), nil
}

func (c ControllerRole) UpdateRoleAccessList(_ context.Context, req *api.UpdateRoleAccessListRequest) (*api.RoleAccessListResult, error) {

	if err := c.endpoint.Role.UpdateAccessList(req.Name, c.dto.Role.FromUpdateRoleAccessListRequest(req)); err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessList, err := c.endpoint.Role.GetAccessList(req.Name, c.accessList)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToRoleAccessListResult(accessList), nil
}

func (c ControllerRole) AddRole(_ context.Context, req *api.NewRoleRequest) (*api.Role, error) {

	role := c.dto.Role.FromNewRoleRequest(req)

	role, errs, err := c.endpoint.Role.Add(role)
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToGRole(role), nil
}

func (c ControllerRole) GetRoleByName(_ context.Context, req *api.GetRoleRequest) (*api.Role, error) {

	role, err := c.endpoint.Role.GetByName(req.Name)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToGRole(role), nil
}

func (c ControllerRole) UpdateRoleByName(_ context.Context, req *api.UpdateRoleRequest) (*api.Role, error) {

	role := c.dto.Role.FromUpdateRoleRequest(req)

	role, errs, err := c.endpoint.Role.Update(role)
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToGRole(role), nil
}

func (c ControllerRole) GetRoleList(_ context.Context, req *api.GetRoleListRequest) (*api.GetRoleListResult, error) {

	items, total, err := c.endpoint.Role.GetList(int64(req.Limit), int64(req.Offset), req.Order, req.SortBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToListResult(items, uint32(total), req.Limit, req.Offset), nil
}

func (c ControllerRole) SearchRoleByName(_ context.Context, req *api.SearchRoleRequest) (*api.SearchRoleListResult, error) {

	items, _, err := c.endpoint.Role.Search(req.Query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Role.ToSearchResult(items), nil
}

func (c ControllerRole) DeleteRoleByName(_ context.Context, req *api.DeleteRoleRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Role.Delete(req.Name); err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
