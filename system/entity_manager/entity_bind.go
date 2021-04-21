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

package entity_manager

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Javascript Binding
//
// Entity
//  .SetState(stateName)
//  .SetAttributes(key, value, args)
//  .SetMetric(name, value)
//
type EntityBind struct {
	id      common.EntityId
	manager *EntityManager
}

func NewEntityBind(id common.EntityId, manager *EntityManager) *EntityBind {
	return &EntityBind{
		id:      id,
		manager: manager,
	}
}

func (e *EntityBind) SetState(stateName string) {
	e.manager.SetState(e.id, EntityStateParams{
		NewState: common.String(stateName),
	})
}

func (e *EntityBind) SetAttributes(params m.EntityAttributeValue) {
	e.manager.SetState(e.id, EntityStateParams{
		AttributeValues: params,
	})
}

func (e *EntityBind) SetMetric(id common.EntityId, name string, value map[string]interface{}) {
	e.manager.SetMetric(id, name, value)
}
