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

type ControllerScript struct {
	*ControllerCommon
}

func NewControllerScript(common *ControllerCommon) ControllerScript {
	return ControllerScript{
		ControllerCommon: common,
	}
}

func (c ControllerScript) AddScript(_ context.Context, req *api.NewScriptRequest) (*api.Script, error) {

	script, errs, err := c.endpoint.Script.Add(c.dto.Script.FromNewScriptRequest(req))
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToGScript(script), nil
}

func (c ControllerScript) GetScriptById(_ context.Context, req *api.GetScriptRequest) (*api.Script, error) {

	script, err := c.endpoint.Script.GetById(int64(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToGScript(script), nil
}

func (c ControllerScript) UpdateScriptById(_ context.Context, req *api.UpdateScriptRequest) (*api.Script, error) {

	script, errs, err := c.endpoint.Script.Update(c.dto.Script.FromUpdateScriptRequest(req))
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToGScript(script), nil
}

func (c ControllerScript) GetScriptList(_ context.Context, req *api.GetScriptListRequest) (*api.GetScriptListResult, error) {

	items, total, err := c.endpoint.Script.GetList(int64(req.Limit), int64(req.Offset), req.Order, req.SortBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToListResult(items, uint32(total), req.Limit, req.Offset), nil
}

func (c ControllerScript) SearchScriptById(_ context.Context, req *api.SearchScriptRequest) (*api.SearchScriptListResult, error) {

	items, _, err := c.endpoint.Script.Search(req.Query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToSearchResult(items), nil
}

func (c ControllerScript) DeleteScriptById(_ context.Context, req *api.DeleteScriptRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Script.DeleteScriptById(int64(req.Id)); err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (c ControllerScript) ExecScriptById(_ context.Context, req *api.ExecScriptRequest) (*api.ExecScriptResult, error) {

	result, err := c.endpoint.Script.Execute(int64(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ExecScriptResult{Result: result}, nil
}

func (c ControllerScript) ExecSrcScriptById(_ context.Context, req *api.ExecSrcScriptRequest) (*api.ExecScriptResult, error) {

	result, err := c.endpoint.Script.ExecuteSource(c.dto.Script.FromExecSrcScriptRequest(req))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ExecScriptResult{Result: result}, nil
}

func (c ControllerScript) CopyScriptById(_ context.Context, req *api.CopyScriptRequest) (*api.Script, error) {

	script, err := c.endpoint.Script.Copy(int64(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Script.ToGScript(script), nil
}
