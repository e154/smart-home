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
	"fmt"

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
	if err = p.pluginManager.EnablePlugin(pluginName); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// Disable ...
func (p *PluginEndpoint) Disable(ctx context.Context, pluginName string) (err error) {
	if err = p.pluginManager.DisablePlugin(pluginName); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// GetList ...
func (p *PluginEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []m.Plugin, total int64, err error) {
	list, total, err = p.adaptors.Plugin.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// GetOptions ...
func (p *PluginEndpoint) GetOptions(ctx context.Context, pluginName string) (options m.PluginOptions, err error) {

	var pl interface{}
	if pl, err = p.pluginManager.GetPlugin(pluginName); err != nil {
		err = errors.Wrap(common.ErrInternal, fmt.Sprintf("name %s", pluginName))
		return
	}

	plugin, ok := pl.(plugins.Plugable)
	if !ok {
		err = errors.Wrap(common.ErrInternal, fmt.Sprintf("name %s", pluginName))
		return
	}

	options = plugin.Options()

	return
}

// Search ...
func (n *PluginEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []m.Plugin, total int64, err error) {

	result, total, err = n.adaptors.Plugin.Search(query, limit, offset)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}
