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
	"github.com/e154/smart-home/internal/api/dto"
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/labstack/echo/v4"
)

// ControllerAutomation ...
type ControllerAutomation struct {
	*ControllerCommon
}

// NewControllerAutomation ...
func NewControllerAutomation(common *ControllerCommon) *ControllerAutomation {
	return &ControllerAutomation{
		ControllerCommon: common,
	}
}

// AddTask ...
func (c ControllerAutomation) AutomationServiceAddTask(ctx echo.Context, _ stub.AutomationServiceAddTaskParams) error {

	obj := &stub.ApiNewTaskRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	task, err := c.endpoint.Task.Add(ctx.Request().Context(), c.dto.Automation.AddTask(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Automation.GetTask(task)))
}

// UpdateTask ...
func (c ControllerAutomation) AutomationServiceUpdateTask(ctx echo.Context, id int64, _ stub.AutomationServiceUpdateTaskParams) error {

	obj := &stub.AutomationServiceUpdateTaskJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	task, err := c.endpoint.Task.Update(ctx.Request().Context(), c.dto.Automation.UpdateTask(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Automation.GetTask(task)))
}

// GetTask ...
func (c ControllerAutomation) AutomationServiceGetTask(ctx echo.Context, id int64) error {

	task, err := c.endpoint.Task.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Automation.GetTask(task)))
}

// GetTaskList ...
func (c ControllerAutomation) AutomationServiceGetTaskList(ctx echo.Context, params stub.AutomationServiceGetTaskListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Task.List(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Automation.ToListResult(items), total, pagination))
}

// DeleteTask ...
func (c ControllerAutomation) AutomationServiceDeleteTask(ctx echo.Context, id int64) error {

	if err := c.endpoint.Task.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// EnableTask ...
func (c ControllerAutomation) AutomationServiceEnableTask(ctx echo.Context, id int64) error {

	if err := c.endpoint.Task.Enable(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DisableTask ...
func (c ControllerAutomation) AutomationServiceDisableTask(ctx echo.Context, id int64) error {

	if err := c.endpoint.Task.Disable(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// ImportTask ...
func (c ControllerAutomation) AutomationServiceImportTask(ctx echo.Context, _ stub.AutomationServiceImportTaskParams) error {

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// AutomationServiceGetStatistic ...
func (c ControllerAutomation) AutomationServiceGetStatistic(ctx echo.Context) error {

	statistic, err := c.endpoint.Automation.Statistic(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, dto.GetStatistic(statistic)))
}
