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

package modbus

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	positionLock        *sync.Mutex
	lat, lon, elevation float64
	solarAzimuth        float64
	solarElevation      float64
	phase               string
	horizonState        string
}

func NewEntityActor(name string, entityManager entity_manager.EntityManager) *EntityActor {

	entity := &EntityActor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntityModbus, name)),
			Name:        name,
			Description: "modbus plugin",
			EntityType:  EntityModbus,
			AttrMu:      &sync.Mutex{},
			Attrs:       NewAttr(),
			Manager:     entityManager,
		},
		positionLock: &sync.Mutex{},
	}

	return entity
}

func (e *EntityActor) Spawn() entity_manager.PluginActor {
	return e
}
