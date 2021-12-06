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
	m "github.com/e154/smart-home/models"
)

// NewEntity ...
func NewEntity(a PluginActor) m.EntityShort {

	info := a.Info()
	actions := make([]m.EntityActionShort, len(info.Actions))
	var i int
	for _, a := range info.Actions {
		actions[i] = m.EntityActionShort{
			Name:        a.Name,
			Description: a.Description,
			ImageUrl:    a.ImageUrl,
			Icon:        a.Icon,
		}
		i++
	}

	states := make([]m.EntityStateShort, len(info.States))
	i = 0
	for _, a := range info.States {
		states[i] = m.EntityStateShort{
			Name:        a.Name,
			Description: a.Description,
			ImageUrl:    a.ImageUrl,
			Icon:        a.Icon,
		}
		i++
	}

	attributes := make(m.Attributes, len(a.Attributes()))
	for k, a := range a.Attributes() {
		attributes[k] = &m.Attribute{
			Name:  a.Name,
			Type:  a.Type,
			Value: a.Value,
		}
	}

	entity := m.EntityShort{
		Id:          info.Id,
		Description: info.Description,
		Type:        info.Type,
		Icon:        info.Icon,
		ImageUrl:    info.ImageUrl,
		Actions:     actions,
		States:      states,
		Attributes:  attributes,
		Area:        info.Area,
		Metrics:     a.Metrics(),
		Hidden:      info.Hidde,
	}
	if cs := info.State; cs != nil {
		entity.State = &m.EntityStateShort{
			Name:        cs.Name,
			Description: cs.Description,
			ImageUrl:    cs.ImageUrl,
			Icon:        cs.Icon,
		}
	}

	return entity
}
