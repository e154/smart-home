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

package endpoint

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
)

// TaskEndpoint ...
type TaskEndpoint struct {
	*CommonEndpoint
}

// NewTaskEndpoint ...
func NewTaskEndpoint(common *CommonEndpoint) *TaskEndpoint {
	return &TaskEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *TaskEndpoint) Add(ctx context.Context, task *m.NewTask) (result *m.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = n.adaptors.Task.Add(ctx, task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", result.Id), events.EventAddedTask{
		Id: id,
	})

	return
}

// Import ...
func (n *TaskEndpoint) Import(ctx context.Context, task *m.Task) (result *m.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	for _, condition := range task.Conditions {
		if condition.Script != nil {
			var engine *scripts.Engine
			if engine, err = n.scriptService.NewEngine(condition.Script); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}

			if err = engine.Compile(); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}
		}
	}
	for _, trigger := range task.Triggers {
		if trigger.Script != nil {
			var engine *scripts.Engine
			if engine, err = n.scriptService.NewEngine(trigger.Script); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}

			if err = engine.Compile(); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}
		}
	}
	for _, action := range task.Actions {
		if action.Script != nil {
			var engine *scripts.Engine
			if engine, err = n.scriptService.NewEngine(action.Script); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}

			if err = engine.Compile(); err != nil {
				err = errors.Wrap(apperr.ErrInternal, err.Error())
				return
			}
		}
	}

	if err = n.adaptors.Task.Import(ctx, task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(ctx, task.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", result.Id), events.EventAddedTask{
		Id: task.Id,
	})

	return
}

// Update ...
func (n *TaskEndpoint) Update(ctx context.Context, task *m.UpdateTask) (result *m.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Task.Update(ctx, task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(ctx, task.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", result.Id), events.EventUpdateTask{
		Id: task.Id,
	})

	return
}

// GetById ...
func (n *TaskEndpoint) GetById(ctx context.Context, id int64) (task *m.Task, err error) {

	if task, err = n.adaptors.Task.GetById(ctx, id); err != nil {
		return
	}
	task.IsLoaded = n.automation.TaskIsLoaded(id)
	task.Telemetry = n.automation.TaskTelemetry(id)
	return
}

// Delete ...
func (n *TaskEndpoint) Delete(ctx context.Context, id int64) (err error) {
	if err = n.adaptors.Task.Delete(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", id), events.EventRemoveTask{
		Id: id,
	})
	return
}

// Enable ...
func (n *TaskEndpoint) Enable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Enable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", id), events.EventEnableTask{
		Id: id,
	})
	return
}

// Disable ...
func (n *TaskEndpoint) Disable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Disable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", id), events.EventDisableTask{
		Id: id,
	})
	return
}

// List ...
func (n *TaskEndpoint) List(ctx context.Context, pagination common.PageParams) (tasks []*m.Task, total int64, err error) {

	if tasks, total, err = n.adaptors.Task.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false); err != nil {
		return
	}
	for _, task := range tasks {
		task.IsLoaded = n.automation.TaskIsLoaded(task.Id)
	}
	return
}
