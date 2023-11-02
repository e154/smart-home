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

package controllers

import (
	"github.com/e154/smart-home/api/stub"
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/dto"
)

// ControllerPlugin ...
type ControllerPlugin struct {
	*ControllerCommon
}

// NewControllerPlugin ...
func NewControllerPlugin(common *ControllerCommon) *ControllerPlugin {
	return &ControllerPlugin{
		ControllerCommon: common,
	}
}

// GetPluginList ...
func (c ControllerPlugin) PluginServiceGetPluginList(ctx echo.Context, params stub.PluginServiceGetPluginListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Plugin.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Plugin.ToPluginListResult(items), total, pagination))
}

// EnablePlugin ...
func (c ControllerPlugin) PluginServiceEnablePlugin(ctx echo.Context, name string) error {

	if err := c.endpoint.Plugin.Enable(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// DisablePlugin ...
func (c ControllerPlugin) PluginServiceDisablePlugin(ctx echo.Context, name string) error {

	if err := c.endpoint.Plugin.Disable(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// SearchPlugin ...
func (c ControllerPlugin) PluginServiceSearchPlugin(ctx echo.Context, params stub.PluginServiceSearchPluginParams) error {

	search := c.Search(params.Query, params.Limit, params.Offset)
	items, _, err := c.endpoint.Plugin.Search(ctx.Request().Context(), search.Query, search.Limit, search.Offset)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, c.dto.Plugin.ToSearchResult(items))
}

// GetPlugin ...
func (c ControllerPlugin) PluginServiceGetPlugin(ctx echo.Context, name string) error {

	plugin, err := c.endpoint.Plugin.GetByName(ctx.Request().Context(), name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	options, err := c.endpoint.Plugin.GetOptions(ctx.Request().Context(), name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Plugin.ToGetPlugin(plugin, options)))
}

// UpdatePluginSettings ...
func (c ControllerPlugin) PluginServiceUpdatePluginSettings(ctx echo.Context, name string, _ stub.PluginServiceUpdatePluginSettingsParams) error {

	obj := &stub.PluginServiceUpdatePluginSettingsJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	if err := c.endpoint.Plugin.UpdateSettings(ctx.Request().Context(), name, dto.AttributeFromApi(obj.Settings)); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
