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

package supervisor

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

// EntityBind ...
type EntityBind struct {
	m.EntityShort
	manager plugins.Supervisor
}

// NewEntityBind ...
func NewEntityBind(id common.EntityId, manager plugins.Supervisor) *EntityBind {
	entity, _ := manager.GetEntityById(id)
	return &EntityBind{
		EntityShort: entity,
		manager:     manager,
	}
}

func GetEntityBind(manager plugins.Supervisor) func(entityId string) *EntityBind {
	return func(entityId string) *EntityBind {
		return NewEntityBind(common.EntityId(entityId), manager)
	}
}
