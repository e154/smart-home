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

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/event_bus/events"
	"github.com/e154/smart-home/system/message_queue"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// DeveloperToolsEndpoint ...
type DeveloperToolsEndpoint struct {
	*CommonEndpoint
}

// NewDeveloperToolsEndpoint ...
func NewDeveloperToolsEndpoint(common *CommonEndpoint) *DeveloperToolsEndpoint {
	return &DeveloperToolsEndpoint{
		CommonEndpoint: common,
	}
}

// StateList ...
func (d DeveloperToolsEndpoint) StateList(ctx context.Context) (states []m.EntityShort, total int64, err error) {
	states, err = d.entityManager.List()
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}
	total = int64(len(states))
	return
}

// SetEntityState ...
func (d DeveloperToolsEndpoint) SetEntityState(ctx context.Context, entityId string, newState *string, attrs map[string]interface{}) (errs validator.ValidationErrorsTranslations, err error) {

	_, err = d.adaptors.Entity.GetById(common.EntityId(entityId))
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	d.eventBus.Publish(event_bus.TopicEntities, events.EventEntitySetState{
		Id:              common.EntityId(entityId),
		NewState:        newState,
		AttributeValues: attrs,
	})

	return
}

// EventList ...
func (d DeveloperToolsEndpoint) EventList(ctx context.Context) (events []message_queue.Stat, total int64, err error) {
	events, err = d.eventBus.Stat()
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}
	total = int64(len(events))
	return
}

// TaskCallTrigger ...
func (d *DeveloperToolsEndpoint) TaskCallTrigger(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Task.GetById(id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	d.eventBus.Publish(event_bus.TopicAutomation, events.EventCallTaskTrigger{
		Id:   id,
		Name: name,
	})

	return
}

// TaskCallAction ...
func (d *DeveloperToolsEndpoint) TaskCallAction(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Task.GetById(id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	d.eventBus.Publish(event_bus.TopicAutomation, events.EventCallTaskAction{
		Id:   id,
		Name: name,
	})

	return
}

// ReloadEntity ...
func (d *DeveloperToolsEndpoint) ReloadEntity(ctx context.Context, id common.EntityId) (err error) {

	_, err = d.adaptors.Entity.GetById(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	d.eventBus.Publish(event_bus.TopicEntities, events.EventUpdatedEntity{
		Id: id,
	})

	return
}

// EntitySetState ...
func (d *DeveloperToolsEndpoint) EntitySetState(ctx context.Context, id common.EntityId, name string) (err error) {

	_, err = d.adaptors.Entity.GetById(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	d.eventBus.Publish(event_bus.TopicEntities, events.EventEntitySetState{
		Id:       id,
		NewState: common.String(name),
	})

	return
}
