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

func (c ControllerAutomation) AddTask(ctx context.Context, req *api.NewTaskRequest) (*api.Task, error) {

	task := c.dto.Automation.AddTask(req)

	task, errs, err := c.endpoint.Task.Add(ctx, task)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Automation.ToTask(task), nil
}

func (c ControllerAutomation) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.Task, error) {

	task := c.dto.Automation.UpdateTask(req)

	task, errs, err := c.endpoint.Task.Update(ctx, task)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Automation.ToTask(task), nil
}

func (c ControllerAutomation) GetTask(ctx context.Context, req *api.GetTaskRequest) (*api.Task, error) {

	task, errs, err := c.endpoint.Task.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Automation.ToTask(task), nil
}

func (c ControllerAutomation) GetTaskList(ctx context.Context, req *api.PaginationRequest) (*api.GetTaskListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, errs, err := c.endpoint.Task.List(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Automation.ToListResult(items, uint64(total), pagination), nil
}

func (c ControllerAutomation) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Task.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

func (c ControllerAutomation) TaskCallTrigger(ctx context.Context, req *api.CallRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Task.TaskCallAction(ctx, req.Id, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

func (c ControllerAutomation) TaskCallAction(ctx context.Context, req *api.CallRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Task.TaskCallAction(ctx, req.Id, req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
