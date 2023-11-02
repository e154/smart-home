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

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
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
	//states, err = d.supervisor.List()
	//if err != nil {
	//	return
	//}
	//total = int64(len(states))
	return
}

// EntitySetState ...
func (d DeveloperToolsEndpoint) EntitySetState(ctx context.Context, entityId string, newState *string, attrs map[string]interface{}) (err error) {

	_, err = d.adaptors.Entity.GetById(ctx, common.EntityId(entityId))
	if err != nil {
		return
	}

	d.eventBus.Publish("system/entities/"+entityId, events.EventEntitySetState{
		EntityId:        common.EntityId(entityId),
		NewState:        newState,
		AttributeValues: attrs,
		StorageSave:     true,
	})

	return
}

// TaskCallTrigger ...
func (d *DeveloperToolsEndpoint) TaskCallTrigger(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Trigger.GetById(ctx, id); err != nil {
		return
	}

	d.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", id), events.EventCallTrigger{
		Id:  id,
		Ctx: context.Background(),
	})

	return
}

// TaskCallAction ...
func (d *DeveloperToolsEndpoint) TaskCallAction(ctx context.Context, id int64, name string) (err error) {

	if _, err = d.adaptors.Action.GetById(ctx, id); err != nil {
		return
	}

	d.eventBus.Publish(fmt.Sprintf("system/automation/actions/%d", id), events.EventCallAction{
		Id:  id,
		Ctx: context.Background(),
	})

	return
}

// ReloadEntity ...
func (d *DeveloperToolsEndpoint) ReloadEntity(ctx context.Context, id common.EntityId) (err error) {

	_, err = d.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	d.eventBus.Publish("system/models/entities/"+id.String(), events.EventUpdatedEntityModel{
		EntityId: id,
	})

	return
}

// EntitySetStateName ...
func (d *DeveloperToolsEndpoint) EntitySetStateName(ctx context.Context, id common.EntityId, name string) (err error) {

	_, err = d.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	d.eventBus.Publish("system/entities/"+id.String(), events.EventEntitySetState{
		EntityId:    id,
		NewState:    common.String(name),
		StorageSave: true,
	})

	return
}

// GetEventBusState ...
func (d *DeveloperToolsEndpoint) GetEventBusState(ctx context.Context, pagination common.PageParams) (bus.Stats, int64, error) {
	return d.eventBus.Stat(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
}
