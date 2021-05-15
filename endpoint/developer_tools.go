// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/message_queue"
	"github.com/e154/smart-home/system/validation"
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
func (d DeveloperToolsEndpoint) StateList() (states []m.EntityShort, total int64, err error) {
	states, err = d.entityManager.List()
	total = int64(len(states))
	return
}

// UpdateState ...
func (d DeveloperToolsEndpoint) UpdateState(entityId string, state *string, attrs map[string]interface{}) (errs []*validation.Error, err error) {
	err = d.entityManager.SetState(common.EntityId(entityId), entity_manager.EntityStateParams{
		NewState:        state,
		AttributeValues: attrs,
	})
	return
}

// EventList ...
func (d DeveloperToolsEndpoint) EventList() (events []message_queue.Stat, total int64, err error) {
	events, err = d.eventBus.Stat()
	total = int64(len(events))
	return
}
