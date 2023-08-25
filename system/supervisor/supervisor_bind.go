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

package supervisor

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// SupervisorBind ...
type SupervisorBind struct {
	manager Supervisor
}

// NewSupervisorBind ...
func NewSupervisorBind(manager Supervisor) *SupervisorBind {
	return &SupervisorBind{manager: manager}
}

// GetEntity ...
func (e *SupervisorBind) GetEntity(id common.EntityId) *EntityBind {
	return NewEntityBind(id, e.manager)
}

// SetState ...
func (e *SupervisorBind) SetState(id common.EntityId, stateName string) {
	_ = e.manager.SetState(id, EntityStateParams{
		NewState: common.String(stateName),
	})
}

// SetAttribute ...
func (e *SupervisorBind) SetAttribute(id common.EntityId, params m.AttributeValue) {
	_ = e.manager.SetState(id, EntityStateParams{
		AttributeValues: params,
	})
}

// SetMetric ...
func (e *SupervisorBind) SetMetric(id common.EntityId, name string, value map[string]float32) {
	e.manager.SetMetric(id, name, value)
}

// CallAction ...
func (e *SupervisorBind) CallAction(id common.EntityId, action string, arg map[string]interface{}) {
	e.manager.CallAction(id, action, arg)
}

// CallScene ...
func (e *SupervisorBind) CallScene(id common.EntityId, arg map[string]interface{}) {
	e.manager.CallScene(id, arg)
}
