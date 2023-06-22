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

package endpoint

import (
	"context"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/go-playground/validator/v10"
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
func (n *TaskEndpoint) Add(ctx context.Context, task *m.Task) (result *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(task); !ok {
		return
	}

	if err = n.adaptors.Task.Add(task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(task.Id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventAddedTask{
		Id: task.Id,
	})

	return
}

// Import ...
func (n *TaskEndpoint) Import(ctx context.Context, task *m.Task) (result *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(task); !ok {
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

	if err = n.adaptors.Task.Import(task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(task.Id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventAddedTask{
		Id: task.Id,
	})

	return
}

// Update ...
func (n *TaskEndpoint) Update(ctx context.Context, task *m.Task) (result *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(task); !ok {
		return
	}

	if err = n.adaptors.Task.Update(task); err != nil {
		return
	}

	if result, err = n.adaptors.Task.GetById(task.Id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventUpdateTask{
		Id: task.Id,
	})

	return
}

// GetById ...
func (n *TaskEndpoint) GetById(ctx context.Context, id int64) (task *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	task, err = n.adaptors.Task.GetById(id)

	return
}

// Delete ...
func (n *TaskEndpoint) Delete(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Delete(id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventRemoveTask{
		Id: id,
	})

	return
}

// Enable ...
func (n *TaskEndpoint) Enable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Enable(id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventEnableTask{
		Id: id,
	})

	return
}

// Disable ...
func (n *TaskEndpoint) Disable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Disable(id); err != nil {
		return
	}

	n.eventBus.Publish(bus.TopicAutomation, events.EventDisableTask{
		Id: id,
	})

	return
}

// List ...
func (n *TaskEndpoint) List(ctx context.Context, pagination common.PageParams) (list []*m.Task, total int64, errs validator.ValidationErrorsTranslations, err error) {

	list, total, err = n.adaptors.Task.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false)

	return
}
