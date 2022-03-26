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

package entity_manager

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// EntityManagerBind ...
type EntityManagerBind struct {
	manager EntityManager
}

// NewEntityManagerBind ...
func NewEntityManagerBind(manager EntityManager) *EntityManagerBind {
	return &EntityManagerBind{manager: manager}
}

// GetEntity ...
func (e *EntityManagerBind) GetEntity(id common.EntityId) *EntityBind {
	return NewEntityBind(id, e.manager)
}

// SetState ...
func (e *EntityManagerBind) SetState(id common.EntityId, stateName string) {
	_ = e.manager.SetState(id, EntityStateParams{
		NewState: common.String(stateName),
	})
}

// SetAttribute ...
func (e *EntityManagerBind) SetAttribute(id common.EntityId, params m.AttributeValue) {
	_ = e.manager.SetState(id, EntityStateParams{
		AttributeValues: params,
	})
}

// SetMetric ...
func (e *EntityManagerBind) SetMetric(id common.EntityId, name string, value map[string]interface{}) {
	e.manager.SetMetric(id, name, value)
}

// CallAction ...
func (e *EntityManagerBind) CallAction(id common.EntityId, action string, arg map[string]interface{}) {
	e.manager.CallAction(id, action, arg)
}
