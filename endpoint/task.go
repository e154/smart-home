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

	"github.com/e154/smart-home/system/event_bus"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
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
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	if result, err = n.adaptors.Task.GetById(task.Id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// Update ...
func (n *TaskEndpoint) Update(ctx context.Context, task *m.Task) (result *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(task); !ok {
		return
	}

	if err = n.adaptors.Task.Update(task); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = n.adaptors.Task.GetById(task.Id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// GetById ...
func (n *TaskEndpoint) GetById(ctx context.Context, id int64) (task *m.Task, errs validator.ValidationErrorsTranslations, err error) {

	if task, err = n.adaptors.Task.GetById(id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// Delete ...
func (n *TaskEndpoint) Delete(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Task.Delete(id); err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			err = errors.Wrap(common.ErrInternal, err.Error())
		}
	}

	return
}

// List ...
func (n *TaskEndpoint) List(ctx context.Context, pagination common.PageParams) (list []*m.Task, total int64, errs validator.ValidationErrorsTranslations, err error) {

	list, total, err = n.adaptors.Task.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// TaskCallTrigger ...
func (n *TaskEndpoint) TaskCallTrigger(ctx context.Context, id int64, name string) (err error) {

	n.eventBus.Publish(event_bus.TopicAutomation, event_bus.EventCallTaskTrigger{
		Id:   id,
		Name: name,
	})

	return
}

// TaskCallAction ...
func (n *TaskEndpoint) TaskCallAction(ctx context.Context, id int64, name string) (err error) {

	n.eventBus.Publish(event_bus.TopicAutomation, event_bus.EventCallTaskAction{
		Id:   id,
		Name: name,
	})

	return
}
