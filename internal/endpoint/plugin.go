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
	"bufio"
	"context"
	"mime/multipart"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/version"
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
	err = p.supervisor.EnablePlugin(ctx, pluginName)
	return
}

// Disable ...
func (p *PluginEndpoint) Disable(ctx context.Context, pluginName string) (err error) {
	err = p.supervisor.DisablePlugin(ctx, pluginName)
	return
}

// GetList ...
func (p *PluginEndpoint) GetList(ctx context.Context, pagination common.PageParams, enabled, triggers *bool) (plugins []*models.Plugin, total int64, err error) {
	if plugins, total, err = p.adaptors.Plugin.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, enabled, triggers); err != nil {
		return
	}
	for _, plugin := range plugins {
		if !plugin.External {
			plugin.Version = version.VersionString
		}
		plugin.IsLoaded = p.supervisor.PluginIsLoaded(plugin.Name)
	}
	return
}

// GetOptions ...
func (p *PluginEndpoint) GetOptions(ctx context.Context, pluginName string) (options models.PluginOptions, err error) {

	var pl interface{}
	if pl, err = p.supervisor.GetPlugin(pluginName); err != nil {
		return
	}

	plugin, ok := pl.(plugins.Pluggable)
	if !ok {
		return
	}

	options = plugin.Options()

	return
}

// GetByName ...
func (p *PluginEndpoint) GetByName(ctx context.Context, pluginName string) (plugin *models.Plugin, err error) {
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, pluginName); err != nil {
		return
	}
	if !plugin.External {
		plugin.Version = version.VersionString
	}
	plugin.IsLoaded = p.supervisor.PluginIsLoaded(plugin.Name)
	return
}

// Search ...
func (p *PluginEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*models.Plugin, total int64, err error) {

	result, total, err = p.adaptors.Plugin.Search(ctx, query, limit, offset)
	if err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}
	return
}

// UpdateSettings ...
func (p *PluginEndpoint) UpdateSettings(ctx context.Context, pluginName string, settings models.Attributes) (err error) {

	var plugin *models.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, pluginName); err != nil {
		return
	}

	plugin.Settings = settings.Serialize()

	if err = p.adaptors.Plugin.Update(ctx, plugin); err != nil {
		return
	}

	if !p.supervisor.PluginIsLoaded(pluginName) {
		return
	}

	if err = p.supervisor.DisablePlugin(ctx, pluginName); err != nil {
		return
	}

	err = p.supervisor.EnablePlugin(ctx, pluginName)

	return
}

func (p *PluginEndpoint) Readme(ctx context.Context, pluginName string, note *string, lang *string) (result []byte, err error) {

	result, err = p.supervisor.GetPluginReadme(ctx, pluginName, note, lang)

	return
}

func (p *PluginEndpoint) RemovePlugin(ctx context.Context, pluginName string) (err error) {

	if p.checkSuperUser(ctx) {
		err = apperr.ErrPluginUploadForbidden
		return
	}

	err = p.supervisor.RemovePlugin(ctx, pluginName)

	return
}

// Upload ...
func (p *PluginEndpoint) Upload(ctx context.Context, files map[string][]*multipart.FileHeader) (pluginList []*models.Plugin, errs []error, err error) {

	if p.checkSuperUser(ctx) {
		err = apperr.ErrPluginUploadForbidden
		return
	}

	pluginList = make([]*models.Plugin, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, _err := fileHeader[0].Open()
		if _err != nil {
			errs = append(errs, _err)
			continue
		}

		reader := bufio.NewReader(file)
		var newPlugin *models.Plugin
		newPlugin, _err = p.supervisor.UploadPlugin(ctx, reader)
		if _err != nil {
			errs = append(errs, _err)
			err = _err
			file.Close()
			return
		}

		file.Close()

		pluginList = append(pluginList, newPlugin)
	}

	return
}
