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
	"github.com/e154/smart-home/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerDeveloperTools ...
type ControllerDeveloperTools struct {
	*ControllerCommon
}

// NewControllerDeveloperTools ...
func NewControllerDeveloperTools(common *ControllerCommon) ControllerDeveloperTools {
	return ControllerDeveloperTools{
		ControllerCommon: common,
	}
}

// EntitySetState ...
func (c ControllerDeveloperTools) EntitySetState(ctx context.Context, req *api.EntityRequest) (*emptypb.Empty, error) {

	err := c.endpoint.DeveloperTools.EntitySetState(ctx, common.EntityId(req.Id), req.Name)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// ReloadEntity ...
func (c ControllerDeveloperTools) ReloadEntity(ctx context.Context, req *api.ReloadRequest) (*emptypb.Empty, error) {

	err := c.endpoint.DeveloperTools.ReloadEntity(ctx, common.EntityId(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// TaskCallTrigger ...
func (c ControllerDeveloperTools) TaskCallTrigger(ctx context.Context, req *api.AutomationRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.DeveloperTools.TaskCallAction(ctx, req.Id, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// TaskCallAction ...
func (c ControllerDeveloperTools) TaskCallAction(ctx context.Context, req *api.AutomationRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.DeveloperTools.TaskCallAction(ctx, req.Id, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
