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

	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// TriggerEndpoint ...
type TriggerEndpoint struct {
	*CommonEndpoint
}

// NewTriggerEndpoint ...
func NewTriggerEndpoint(common *CommonEndpoint) *TriggerEndpoint {
	return &TriggerEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *TriggerEndpoint) Add(ctx context.Context, trigger *m.Trigger) (result *m.Trigger, err error) {

	if ok, errs := n.validation.Valid(trigger); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if trigger.Id, err = n.adaptors.Trigger.Add(ctx, trigger); err != nil {
		return
	}

	if result, err = n.adaptors.Trigger.GetById(ctx, trigger.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", trigger.Id), events.EventAddedTrigger{
		Id: trigger.Id,
	})

	return
}

// GetById ...
func (n *TriggerEndpoint) GetById(ctx context.Context, id int64) (trigger *m.Trigger, err error) {

	trigger, err = n.adaptors.Trigger.GetById(ctx, id)
	trigger.IsLoaded = n.automation.TriggerIsLoaded(id)
	return
}

// Update ...
func (n *TriggerEndpoint) Update(ctx context.Context, trigger *m.Trigger) (result *m.Trigger, err error) {

	_, err = n.adaptors.Trigger.GetById(ctx, trigger.Id)
	if err != nil {
		return
	}

	if ok, errs := n.validation.Valid(trigger); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Trigger.Update(ctx, trigger); err != nil {
		return
	}

	if result, err = n.adaptors.Trigger.GetById(ctx, trigger.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", trigger.Id), events.EventUpdatedTrigger{
		Id: trigger.Id,
	})

	return
}

// GetList ...
func (n *TriggerEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Trigger, total int64, err error) {

	if result, total, err = n.adaptors.Trigger.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false); err != nil {
		return
	}

	for _, trigger := range result {
		trigger.IsLoaded = n.automation.TriggerIsLoaded(trigger.Id)
	}

	return
}

// Delete ...
func (n *TriggerEndpoint) Delete(ctx context.Context, id int64) (err error) {

	var trigger *m.Trigger
	trigger, err = n.adaptors.Trigger.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = n.adaptors.Trigger.Delete(ctx, trigger.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", trigger.Id), events.EventRemovedTrigger{
		Id: trigger.Id,
	})

	return
}

// Search ...
func (n *TriggerEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Trigger, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Trigger.Search(ctx, query, int(limit), int(offset))

	return
}

// Enable ...
func (n *TriggerEndpoint) Enable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Trigger.Enable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", id), events.EventEnableTrigger{
		Id: id,
	})

	return
}

// Disable ...
func (n *TriggerEndpoint) Disable(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Trigger.Disable(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", id), events.EventDisableTrigger{
		Id: id,
	})

	return
}
