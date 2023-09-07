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
		Items: make([]*api.PluginShort, 0, len(items)),
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}

	for _, item := range items {
		result.Items = append(result.Items, &api.PluginShort{
			Name:     item.Name,
			Version:  item.Version,
			Enabled:  item.Enabled,
			System:   item.System,
			IsLoaded: common.Bool(item.IsLoaded),
		})
	}

	return
}

// Options ...
func (p Plugin) Options(options m.PluginOptions) (result *api.PluginOptionsResult) {

	var actions = make(map[string]*api.PluginOptionsResult_EntityAction)
	for k, v := range options.ActorActions {
		actions[k] = &api.PluginOptionsResult_EntityAction{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	var states = make(map[string]*api.PluginOptionsResult_EntityState)
	for k, v := range options.ActorStates {
		states[k] = &api.PluginOptionsResult_EntityState{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	result = &api.PluginOptionsResult{
		Triggers:           options.Triggers,
		Actors:             options.Actors,
		ActorCustomAttrs:   options.ActorCustomAttrs,
		ActorAttrs:         AttributeToApi(options.ActorAttrs),
		ActorCustomActions: options.ActorCustomActions,
		ActorActions:       actions,
		ActorCustomStates:  options.ActorCustomStates,
		ActorStates:        states,
		ActorCustomSetts:   options.ActorCustomSetts,
		ActorSetts:         AttributeToApi(options.ActorSetts),
		Setts:              AttributeToApi(options.Setts),
	}
	return
}

// ToSearchResult ...
func (p Plugin) ToSearchResult(list []*m.Plugin) *api.SearchPluginResult {

	items := make([]*api.PluginShort, 0, len(list))

	for _, i := range list {
		items = append(items, &api.PluginShort{
			Name:    i.Name,
			Version: i.Version,
			Enabled: i.Enabled,
			System:  i.System,
		})
	}

	return &api.SearchPluginResult{
		Items: items,
	}
}

func (p Plugin) ToGetPlugin(plugin *m.Plugin, options m.PluginOptions) (result *api.Plugin) {

	var settings = make(map[string]*api.Attribute)
	if options.Setts != nil && plugin.Settings != nil {
		setts := options.Setts.Copy()
		setts.Deserialize(plugin.Settings)
		settings = AttributeToApi(setts)
	}
	result = &api.Plugin{
		Name:     plugin.Name,
		Version:  plugin.Version,
		Enabled:  plugin.Enabled,
		System:   plugin.System,
		Actor:    plugin.Actor,
		Settings: settings,
		Options:  p.Options(options),
		IsLoaded: common.Bool(plugin.IsLoaded),
	}
	return
}
