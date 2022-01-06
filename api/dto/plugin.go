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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Plugin ...
type Plugin struct{}

// NewPluginDto ...
func NewPluginDto() Plugin {
	return Plugin{}
}

// ToPluginListResult ...
func (p Plugin) ToPluginListResult(items []*m.Plugin, total uint64, pagination common.PageParams) (result *api.GetPluginListResult) {

	result = &api.GetPluginListResult{
		Items: make([]*api.Plugin, 0, len(items)),
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}

	for _, item := range items {
		result.Items = append(result.Items, &api.Plugin{
			Name:     item.Name,
			Version:  item.Version,
			Enabled:  item.Enabled,
			System:   item.System,
			Settings: AttributeToApi(item.Settings),
		})
	}

	return
}

// Options ...
func (p Plugin) Options(options m.PluginOptions) (result *api.GetPluginOptionsResult) {

	var actions = make(map[string]*api.GetPluginOptionsResult_EntityAction)
	for k, v := range options.ActorActions {
		actions[k] = &api.GetPluginOptionsResult_EntityAction{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	var states = make(map[string]*api.GetPluginOptionsResult_EntityState)
	for k, v := range options.ActorStates {
		states[k] = &api.GetPluginOptionsResult_EntityState{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	result = &api.GetPluginOptionsResult{
		Triggers:           options.Triggers,
		Actors:             options.Actors,
		ActorCustomAttrs:   options.ActorCustomAttrs,
		ActorAttrs:         AttributeToApi(options.ActorAttrs),
		ActorCustomActions: options.ActorCustomActions,
		ActorActions:       actions,
		ActorCustomStates:  options.ActorCustomStates,
		ActorStates:        states,
		ActorSetts:         AttributeToApi(options.ActorSetts),
		Setts:              AttributeToApi(options.Setts),
	}
	return
}
