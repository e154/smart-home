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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/plugins"
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
func (p *PluginEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.Plugin, total int64, err error) {
	var pluginList []common.PluginInfo
	if pluginList, total, err = p.pluginManager.PluginList(); err != nil {
		return
	}
	list = make([]*m.Plugin, 0, len(pluginList))
	for _, p := range pluginList {
		list = append(list, &m.Plugin{
			Name:    p.Name,
			Version: p.Version,
			Enabled: p.Enabled,
			System:  p.System,
		})
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
