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

// ControllerAutomation ...
type ControllerAutomation struct {
	*ControllerCommon
}

// NewControllerAutomation ...
func NewControllerAutomation(common *ControllerCommon) ControllerAutomation {
	return ControllerAutomation{
		ControllerCommon: common,
	}
}

func (c ControllerAutomation) AddTask(ctx context.Context, request *api.NewTaskRequest) (*api.Task, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerAutomation) UpdateTaskByName(ctx context.Context, request *api.UpdateTaskRequest) (*api.Task, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerAutomation) GetTaskByName(ctx context.Context, request *api.GetTaskRequest) (*api.Task, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerAutomation) GetTaskList(ctx context.Context, request *api.GetTaskListRequest) (*api.GetTaskListResult, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}

func (c ControllerAutomation) DeleteTaskByName(ctx context.Context, request *api.DeleteTaskRequest) (*emptypb.Empty, error) {
	return nil, c.error(ctx, nil, common.ErrUnimplemented)
}
