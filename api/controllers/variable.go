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

// ControllerVariable ...
type ControllerVariable struct {
	*ControllerCommon
}

// NewControllerVariable ...
func NewControllerVariable(common *ControllerCommon) ControllerVariable {
	return ControllerVariable{
		ControllerCommon: common,
	}
}

// AddVariable ...
func (c ControllerVariable) AddVariable(ctx context.Context, req *api.NewVariableRequest) (*api.Variable, error) {

	variable := c.dto.Variable.AddVariable(req)

	errs, err := c.endpoint.Variable.Add(ctx, variable)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Variable.ToVariable(variable), nil
}

// UpdateVariable ...
func (c ControllerVariable) UpdateVariable(ctx context.Context, req *api.UpdateVariableRequest) (*api.Variable, error) {

	variable := c.dto.Variable.UpdateVariable(req)

	errs, err := c.endpoint.Variable.Update(ctx, variable)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Variable.ToVariable(variable), nil
}

// GetVariableByName ...
func (c ControllerVariable) GetVariableByName(ctx context.Context, req *api.GetVariableRequest) (*api.Variable, error) {

	variable, err := c.endpoint.Variable.GetById(ctx, req.Name)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Variable.ToVariable(variable), nil
}

// GetVariableList ...
func (c ControllerVariable) GetVariableList(ctx context.Context, req *api.PaginationRequest) (*api.GetVariableListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.Variable.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Variable.ToListResult(items, uint64(total), pagination), nil
}

// DeleteVariable ...
func (c ControllerVariable) DeleteVariable(ctx context.Context, req *api.DeleteVariableRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Variable.Delete(ctx, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}