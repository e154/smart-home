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

package controllers

import (
	"context"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerPlugin ...
type ControllerPlugin struct {
	*ControllerCommon
}

// NewControllerPlugin ...
func NewControllerPlugin(common *ControllerCommon) ControllerPlugin {
	return ControllerPlugin{
		ControllerCommon: common,
	}
}

// GetPluginList ...
func (c ControllerPlugin) GetPluginList(ctx context.Context, req *api.GetPluginListRequest) (*api.GetPluginListResult, error) {

	items, total, err := c.endpoint.Plugin.GetList(int64(req.Limit), int64(req.Offset), req.Order, req.SortBy)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Plugin.ToPluginListResult(items, uint32(total), req.Limit, req.Offset), nil
}

// EnablePlugin ...
func (c ControllerPlugin) EnablePlugin(ctx context.Context, req *api.EnablePluginRequest) (*api.EnablePluginResult, error) {

	if err := c.endpoint.Plugin.Enable(req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &api.EnablePluginResult{}, nil
}

// DisablePlugin ...
func (c ControllerPlugin) DisablePlugin(ctx context.Context, req *api.DisablePluginRequest) (*api.DisablePluginResult, error) {

	if err := c.endpoint.Plugin.Disable(req.Name); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &api.DisablePluginResult{}, nil
}

// GetPluginOptions ...
func (c ControllerPlugin) GetPluginOptions(ctx context.Context, req *api.GetPluginOptionsRequest) (*api.GetPluginOptionsResult, error) {

	options, err := c.endpoint.Plugin.GetOptions(req.Name)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Plugin.Options(options), nil
}
