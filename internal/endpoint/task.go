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

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"
	"github.com/pkg/errors"
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
func (n *TaskEndpoint) Add(ctx context.Context, task *models.NewTask) (result *models.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrValidation
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

	n.eventBus.Publish(fmt.Sprintf("system/models/tasks/%d", result.Id), events.EventCreatedTaskModel{
		Id: id,
	})

	log.Infof("added new task %s id:(%d)", result.Name, result.Id)

	return
}

// Import ...
func (n *TaskEndpoint) Import(ctx context.Context, task *models.Task) (result *models.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	for _, condition := range task.Conditions {
		if condition.Script != nil {
			var engine scripts.Engine
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
			var engine scripts.Engine
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
			var engine scripts.Engine
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

	n.eventBus.Publish(fmt.Sprintf("system/models/tasks/%d", result.Id), events.EventCreatedTaskModel{
		Id: task.Id,
	})

	log.Infof("imported task %s id:(%d)", result.Name, result.Id)

	return
}

// Update ...
func (n *TaskEndpoint) Update(ctx context.Context, task *models.UpdateTask) (result *models.Task, err error) {

	if ok, errs := n.validation.Valid(task); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		//triggers
		if err = n.adaptors.Task.DeleteTrigger(ctx, task.Id); err != nil {
			return err
		}

		//conditions
		if err = n.adaptors.Task.DeleteCondition(ctx, task.Id); err != nil {
			return err
		}

		//actions
		if err = n.adaptors.Task.DeleteAction(ctx, task.Id); err != nil {
			return err
		}

		newTask := &models.Task{
			Id:          task.Id,
			Name:        task.Name,
			Description: task.Description,
			Enabled:     task.Enabled,
			Condition:   task.Condition,
			AreaId:      task.AreaId,
		}

		//triggers
		for _, id := range task.TriggerIds {
			newTask.Triggers = append(newTask.Triggers, &models.Trigger{Id: id})
		}

		//conditions
		for _, id := range task.ConditionIds {
			newTask.Conditions = append(newTask.Conditions, &models.Condition{Id: id})
		}

		//actions
		for _, id := range task.ActionIds {
			newTask.Actions = append(newTask.Actions, &models.Action{Id: id})
		}

		return n.adaptors.Task.Update(ctx, newTask)
	})

	if err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(ctx, task.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/tasks/%d", result.Id), events.EventUpdatedTaskModel{
		Id: task.Id,
	})

	log.Infof("updated task %s id:(%d)", result.Name, result.Id)

	return
}

// GetById ...
func (n *TaskEndpoint) GetById(ctx context.Context, id int64) (task *models.Task, err error) {

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

	n.eventBus.Publish(fmt.Sprintf("system/models/tasks/%d", id), events.EventRemovedTaskModel{
		Id: id,
	})

	log.Infof("task id:(%d) was deleted", id)

	return
}

// Enable ...
func (n *TaskEndpoint) Enable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Enable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", id), events.CommandEnableTask{
		Id: id,
	})
	return
}

// Disable ...
func (n *TaskEndpoint) Disable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Disable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", id), events.CommandDisableTask{
		Id: id,
	})

	log.Infof("task %d was deleted", id)

	return
}

// List ...
func (n *TaskEndpoint) List(ctx context.Context, pagination common.PageParams) (tasks []*models.Task, total int64, err error) {

	if tasks, total, err = n.adaptors.Task.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false); err != nil {
		return
	}
	for _, task := range tasks {
		task.IsLoaded = n.automation.TaskIsLoaded(task.Id)
	}
	return
}
