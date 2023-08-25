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

	events2 "github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/go-playground/validator/v10"
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
	states, err = d.supervisor.List()
	if err != nil {
		return
	}
	total = int64(len(states))
	return
}

// SetEntityState ...
func (d DeveloperToolsEndpoint) SetEntityState(ctx context.Context, entityId string, newState *string, attrs map[string]interface{}) (errs validator.ValidationErrorsTranslations, err error) {

	_, err = d.adaptors.Entity.GetById(common.EntityId(entityId))
	if err != nil {
		return
	}

	d.eventBus.Publish(bus.TopicEntities, events2.EventEntitySetState{
		EntityId:        common.EntityId(entityId),
		NewState:        newState,
		AttributeValues: attrs,
	})

	return
}

// EventList ...
func (d DeveloperToolsEndpoint) EventList(ctx context.Context) (events []bus.Stat, total int64, err error) {
	events, total, err = d.eventBus.Stat()
	return
}

// TaskCallTrigger ...
func (d *DeveloperToolsEndpoint) TaskCallTrigger(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Task.GetById(id); err != nil {
		return
	}

	d.eventBus.Publish(bus.TopicAutomation, events2.EventCallTaskTrigger{
		Id:   id,
		Name: name,
	})

	return
}

// TaskCallAction ...
func (d *DeveloperToolsEndpoint) TaskCallAction(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Task.GetById(id); err != nil {
		return
	}

	d.eventBus.Publish(bus.TopicAutomation, events2.EventCallTaskAction{
		Id:   id,
		Name: name,
	})

	return
}

// ReloadEntity ...
func (d *DeveloperToolsEndpoint) ReloadEntity(ctx context.Context, id common.EntityId) (err error) {

	_, err = d.adaptors.Entity.GetById(id)
	if err != nil {
		return
	}

	d.eventBus.Publish(bus.TopicEntities, events2.EventUpdatedEntity{
		EntityId: id,
	})

	return
}

// EntitySetState ...
func (d *DeveloperToolsEndpoint) EntitySetState(ctx context.Context, id common.EntityId, name string) (err error) {

	_, err = d.adaptors.Entity.GetById(id)
	if err != nil {
		return
	}

	d.eventBus.Publish(bus.TopicEntities, events2.EventEntitySetState{
		EntityId: id,
		NewState: common.String(name),
	})

	return
}
