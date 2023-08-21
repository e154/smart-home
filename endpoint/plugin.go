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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/plugins"
	"github.com/pkg/errors"
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
	err = p.pluginManager.EnablePlugin(pluginName)
	return
}

// Disable ...
func (p *PluginEndpoint) Disable(ctx context.Context, pluginName string) (err error) {
	err = p.pluginManager.DisablePlugin(pluginName)
	return
}

// GetList ...
func (p *PluginEndpoint) GetList(ctx context.Context, pagination common.PageParams) (plugins []*m.Plugin, total int64, err error) {
	if plugins, total, err = p.adaptors.Plugin.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy); err != nil {
		return
	}
	for _, plugin := range plugins {
		plugin.IsLoaded = p.pluginManager.IsLoaded(plugin.Name)
	}
	return
}

// GetOptions ...
func (p *PluginEndpoint) GetOptions(ctx context.Context, pluginName string) (options m.PluginOptions, err error) {

	var pl interface{}
	if pl, err = p.pluginManager.GetPlugin(pluginName); err != nil {
		return
	}

	plugin, ok := pl.(plugins.Plugable)
	if !ok {
		return
	}

	options = plugin.Options()

	return
}

// GetByName ...
func (p *PluginEndpoint) GetByName(ctx context.Context, pluginName string) (plugin *m.Plugin, err error) {
	if plugin, err = p.adaptors.Plugin.GetByName(pluginName); err != nil {
		return
	}
	plugin.IsLoaded = p.pluginManager.IsLoaded(plugin.Name)
	return
}

// Search ...
func (n *PluginEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Plugin, total int64, err error) {

	result, total, err = n.adaptors.Plugin.Search(query, limit, offset)
	if err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}
	return
}

// UpdateSettings ...
func (n *PluginEndpoint) UpdateSettings(ctx context.Context, name string, settings m.Attributes) (err error) {

	var plugin *m.Plugin
	if plugin, err = n.adaptors.Plugin.GetByName(name); err != nil {
		return
	}

	plugin.Settings = settings.Serialize()

	if err = n.adaptors.Plugin.Update(plugin); err != nil {
		return
	}

	if err = n.pluginManager.DisablePlugin(name); err != nil {
		return
	}

	err = n.pluginManager.EnablePlugin(name)

	return
}
