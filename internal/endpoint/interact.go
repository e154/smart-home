// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
)

// InteractEndpoint ...
type InteractEndpoint struct {
	*CommonEndpoint
}

// NewInteractEndpoint ...
func NewInteractEndpoint(common *CommonEndpoint) *InteractEndpoint {
	return &InteractEndpoint{
		CommonEndpoint: common,
	}
}

// EntityCallAction ...
func (d InteractEndpoint) EntityCallAction(ctx context.Context, entityId *string, actionName string, areaId *int64, tags []string, args map[string]interface{}) (err error) {

	if entityId != nil {
		id := common.EntityId(*entityId)
		if _, err = d.adaptors.Entity.GetById(ctx, id); err != nil {
			return
		}

		d.eventBus.Publish("system/entities/"+id.String(), events.EventCallEntityAction{
			PluginName: common.String(id.PluginName()),
			EntityId:   id.Ptr(),
			ActionName: actionName,
			Args:       args,
		})
		return
	}

	d.eventBus.Publish("system/entities/", events.EventCallEntityAction{
		ActionName: actionName,
		Args:       args,
		Tags:       tags,
		AreaId:     areaId,
	})

	return
}
