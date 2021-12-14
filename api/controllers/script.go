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

// ControllerScript ...
type ControllerScript struct {
	*ControllerCommon
}

// NewControllerScript ...
func NewControllerScript(common *ControllerCommon) ControllerScript {
	return ControllerScript{
		ControllerCommon: common,
	}
}

// AddScript ...
func (c ControllerScript) AddScript(ctx context.Context, req *api.NewScriptRequest) (*api.Script, error) {

	script, errs, err := c.endpoint.Script.Add(ctx, c.dto.Script.FromNewScriptRequest(req))
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Script.ToGScript(script), nil
}

// GetScriptById ...
func (c ControllerScript) GetScriptById(ctx context.Context, req *api.GetScriptRequest) (*api.Script, error) {

	script, err := c.endpoint.Script.GetById(ctx, int64(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Script.ToGScript(script), nil
}

// UpdateScriptById ...
func (c ControllerScript) UpdateScriptById(ctx context.Context, req *api.UpdateScriptRequest) (*api.Script, error) {

	script, errs, err := c.endpoint.Script.Update(ctx, c.dto.Script.FromUpdateScriptRequest(req))
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Script.ToGScript(script), nil
}

// GetScriptList ...
func (c ControllerScript) GetScriptList(ctx context.Context, req *api.GetScriptListRequest) (*api.GetScriptListResult, error) {

	pagination := c.Pagination(req.Limit, req.Offset, req.Order, req.SortBy)
	items, total, err := c.endpoint.Script.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Script.ToListResult(items, uint64(total), req.Limit, req.Offset), nil
}

// SearchScriptById ...
func (c ControllerScript) SearchScriptById(ctx context.Context, req *api.SearchScriptRequest) (*api.SearchScriptListResult, error) {

	search := c.Search(req.Query, req.Limit, req.Offset)
	items, _, err := c.endpoint.Script.Search(ctx, search.Query, search.Limit, search.Offset)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Script.ToSearchResult(items), nil
}

// DeleteScriptById ...
func (c ControllerScript) DeleteScriptById(ctx context.Context, req *api.DeleteScriptRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Script.DeleteScriptById(ctx, int64(req.Id)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// ExecScriptById ...
func (c ControllerScript) ExecScriptById(ctx context.Context, req *api.ExecScriptRequest) (*api.ExecScriptResult, error) {

	result, err := c.endpoint.Script.Execute(ctx, int64(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &api.ExecScriptResult{Result: result}, nil
}

// ExecSrcScriptById ...
func (c ControllerScript) ExecSrcScriptById(ctx context.Context, req *api.ExecSrcScriptRequest) (*api.ExecScriptResult, error) {

	result, err := c.endpoint.Script.ExecuteSource(ctx, c.dto.Script.FromExecSrcScriptRequest(req))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &api.ExecScriptResult{Result: result}, nil
}

// CopyScriptById ...
func (c ControllerScript) CopyScriptById(ctx context.Context, req *api.CopyScriptRequest) (*api.Script, error) {

	script, err := c.endpoint.Script.Copy(ctx, int64(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Script.ToGScript(script), nil
}
