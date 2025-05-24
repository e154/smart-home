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
	"net/http"

	"github.com/e154/smart-home/internal/api/dto"
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/labstack/echo/v4"
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
	items, total, err := c.endpoint.Plugin.GetList(ctx.Request().Context(), pagination, params.Enabled, params.Triggers)
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

	options, _ := c.endpoint.Plugin.GetOptions(ctx.Request().Context(), name)

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

// PluginServiceGetPluginReadme ...
func (c ControllerPlugin) PluginServiceGetPluginReadme(ctx echo.Context, name string, params stub.PluginServiceGetPluginReadmeParams) error {

	html, err := c.endpoint.Plugin.Readme(ctx.Request().Context(), name, params.Note, params.Lang)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return ctx.HTMLBlob(200, html)
}

// RemovePlugin ...
func (c ControllerPlugin) PluginServiceRemovePlugin(ctx echo.Context, name string) error {

	if err := c.endpoint.Plugin.RemovePlugin(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, "OK")
}

// Upload plugin ...
func (c ControllerPlugin) PluginServiceUploadPlugin(ctx echo.Context, params stub.PluginServiceUploadPluginParams) error {

	r := ctx.Request()

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Error(err.Error())
	}

	form := r.MultipartForm
	if len(form.File) == 0 {
		return c.ERROR(ctx, apperr.ErrInvalidRequest)
	}

	list, errs, err := c.endpoint.Plugin.Upload(r.Context(), form.File)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	var resultPlugins = make([]interface{}, 0)

	for _, file := range list {
		resultPlugins = append(resultPlugins, map[string]string{
			"name": file.Name,
		})
	}

	return c.HTTP200(ctx, map[string]interface{}{
		"files":  resultPlugins,
		"errors": errs,
	})
}

func (c ControllerPlugin) Custom(w http.ResponseWriter, r *http.Request) {
	c.endpoint.Plugin.Custom(w, r)
}
