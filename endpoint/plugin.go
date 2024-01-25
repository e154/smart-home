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

package endpoint

import (
	"context"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// PluginEndpoint ...
type PluginEndpoint struct {
	*CommonEndpoint
}

// NewPluginEndpoint ...
func NewPluginEndpoint(common *CommonEndpoint) *PluginEndpoint {
	return &PluginEndpoint{
		CommonEndpoint: common,
	}
}

// Enable ...
func (p *PluginEndpoint) Enable(ctx context.Context, pluginName string) (err error) {
	err = p.supervisor.EnablePlugin(ctx, pluginName)
	return
}

// Disable ...
func (p *PluginEndpoint) Disable(ctx context.Context, pluginName string) (err error) {
	err = p.supervisor.DisablePlugin(ctx, pluginName)
	return
}

// GetList ...
func (p *PluginEndpoint) GetList(ctx context.Context, pagination common.PageParams) (plugins []*m.Plugin, total int64, err error) {
	if plugins, total, err = p.adaptors.Plugin.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false); err != nil {
		return
	}
	for _, plugin := range plugins {
		plugin.IsLoaded = p.supervisor.PluginIsLoaded(plugin.Name)
	}
	return
}

// GetOptions ...
func (p *PluginEndpoint) GetOptions(ctx context.Context, pluginName string) (options m.PluginOptions, err error) {

	var pl interface{}
	if pl, err = p.supervisor.GetPlugin(pluginName); err != nil {
		return
	}

	plugin, ok := pl.(supervisor.Pluggable)
	if !ok {
		return
	}

	options = plugin.Options()

	return
}

// GetByName ...
func (p *PluginEndpoint) GetByName(ctx context.Context, pluginName string) (plugin *m.Plugin, err error) {
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, pluginName); err != nil {
		return
	}
	plugin.IsLoaded = p.supervisor.PluginIsLoaded(plugin.Name)
	return
}

// Search ...
func (p *PluginEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Plugin, total int64, err error) {

	result, total, err = p.adaptors.Plugin.Search(ctx, query, limit, offset)
	if err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}
	return
}

// UpdateSettings ...
func (p *PluginEndpoint) UpdateSettings(ctx context.Context, name string, settings m.Attributes) (err error) {

	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, name); err != nil {
		return
	}

	plugin.Settings = settings.Serialize()

	if err = p.adaptors.Plugin.Update(ctx, plugin); err != nil {
		return
	}

	if !p.supervisor.PluginIsLoaded(name) {
		return
	}

	if err = p.supervisor.DisablePlugin(ctx, name); err != nil {
		return
	}

	err = p.supervisor.EnablePlugin(ctx, name)

	return
}

func (p *PluginEndpoint) Readme(ctx context.Context, name string, lang *string) (result []byte, err error) {

	result, err = p.supervisor.GetPluginReadme(ctx, name, lang)

	return
}
