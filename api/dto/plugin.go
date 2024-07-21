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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
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
func (p Plugin) ToPluginListResult(plugins []*m.Plugin) []*stub.ApiPluginShort {

	var items = make([]*stub.ApiPluginShort, 0, len(plugins))

	for _, item := range plugins {
		items = append(items, &stub.ApiPluginShort{
			Name:     item.Name,
			Version:  item.Version,
			Enabled:  item.Enabled,
			System:   item.System,
			IsLoaded: common.Bool(item.IsLoaded),
		})
	}

	return items
}

// Options ...
func (p Plugin) Options(options m.PluginOptions) (result *stub.ApiPluginOptionsResult) {

	var actions = make(map[string]stub.ApiPluginOptionsResultEntityAction)
	for k, v := range options.ActorActions {
		actions[k] = stub.ApiPluginOptionsResultEntityAction{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	var states = make(map[string]stub.ApiPluginOptionsResultEntityState)
	for k, v := range options.ActorStates {
		states[k] = stub.ApiPluginOptionsResultEntityState{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    common.StringValue(v.ImageUrl),
			Icon:        common.StringValue(v.Icon),
		}
	}

	result = &stub.ApiPluginOptionsResult{
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
		TriggerParams:      make(stub.ApiTriggerParams),
	}

	for key, value := range options.TriggerParams {
		if result.TriggerParams[key] == nil {
			result.TriggerParams[key] = make([]stub.ApiTriggerParamsField, 0)
		}
		for _, val := range value {
			field := stub.ApiTriggerParamsField{
				Description: val.Description,
				Title:       val.Title,
				Type:        string(val.Type),
			}
			result.TriggerParams[key] = append(result.TriggerParams[key], field)
		}
	}
	return
}

// ToSearchResult ...
func (p Plugin) ToSearchResult(list []*m.Plugin) *stub.ApiSearchPluginResult {

	items := make([]stub.ApiPluginShort, 0, len(list))

	for _, i := range list {
		items = append(items, stub.ApiPluginShort{
			Name:    i.Name,
			Version: i.Version,
			Enabled: i.Enabled,
			System:  i.System,
		})
	}

	return &stub.ApiSearchPluginResult{
		Items: items,
	}
}

func (p Plugin) ToGetPlugin(plugin *m.Plugin, options m.PluginOptions) (result *stub.ApiPlugin) {

	var settings = make(map[string]stub.ApiAttribute)
	if options.Setts != nil && plugin.Settings != nil {
		setts := options.Setts.Copy()
		setts.Deserialize(plugin.Settings)
		settings = AttributeToApi(setts)
	}
	result = &stub.ApiPlugin{
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
