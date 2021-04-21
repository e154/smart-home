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

type Entity struct {
	Id          common.EntityId    `json:"unique_id"`
	Type        common.EntityType  `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Icon        *common.Icon            `json:"icon"`
	ImageUrl    *string            `json:"image_url"`
	Actions     []EntityAction     `json:"actions"`
	States      []EntityState      `json:"states"`
	State       *EntityState       `json:"state"`
	Attributes  m.EntityAttributes `json:"attributes"`
	Area        *m.Area            `json:"area"`
	Metrics     []m.Metric         `json:"metrics"`
	Hidden      bool               `json:"hidden"`
}

func NewEntity(a IActor) Entity {

	info := a.Info()
	actions := make([]EntityAction, len(info.Actions))
	var i int
	for _, a := range info.Actions {
		actions[i] = EntityAction{
			Name:        a.Name,
			Description: a.Description,
			ImageUrl:    a.ImageUrl,
			Icon:        a.Icon,
		}
		i++
	}

	states := make([]EntityState, len(info.States))
	i = 0
	for _, a := range info.States {
		states[i] = EntityState{
			Name:        a.Name,
			Description: a.Description,
			ImageUrl:    a.ImageUrl,
			Icon:        a.Icon,
		}
		i++
	}

	attributes := make(m.EntityAttributes, len(a.Attributes()))
	for k, a := range a.Attributes() {
		attributes[k] = &m.EntityAttribute{
			Name:  a.Name,
			Type:  a.Type,
			Value: a.Value,
		}
	}

	entity := Entity{
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
		entity.State = &EntityState{
			Name:        cs.Name,
			Description: cs.Description,
			ImageUrl:    cs.ImageUrl,
			Icon:        cs.Icon,
		}
	}

	return entity
}
