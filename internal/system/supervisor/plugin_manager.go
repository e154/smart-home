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

package supervisor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/e154/bus"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"go.uber.org/atomic"
)

var pluginsDir = path.Join("data", "plugins")

const GoPluginsEnabled = false

type pluginManager struct {
	*ExternalPlugins
	adaptors       *adaptors.Adaptors
	isStarted      *atomic.Bool
	service        *service
	eventBus       bus.Bus
	enabledPlugins sync.Map
	pluginsWg      *sync.WaitGroup
}

// Start ...
func (p *pluginManager) Start(ctx context.Context) {
	if p.isStarted.Load() {
		return
	}
	defer p.isStarted.Store(true)

	p.loadPlugins(ctx)

	log.Info("Started")
}

// Shutdown ...
func (p *pluginManager) Shutdown(ctx context.Context) {

	if !p.isStarted.Load() {
		return
	}
	defer p.isStarted.Store(false)

	p.enabledPlugins.Range(func(name, value any) bool {
		if enabled, _ := value.(bool); !enabled {
			return true
		}
		_ = p.unloadPlugin(ctx, name.(string))
		return true
	})

	p.pluginsWg.Wait()

	log.Info("Shutdown")
}

// GetPlugin ...
func (p *pluginManager) GetPlugin(t string) (plugin interface{}, err error) {

	plugin, err = p.getPlugin(t)

	return
}

func (p *pluginManager) getPlugin(name string) (plugin plugins.Pluggable, err error) {

	if item, ok := IsPluginRegistered(name); ok {
		plugin = item.(plugins.Pluggable)
		return
	}

	err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrNotFound)

	return
}

func (p *pluginManager) loadPlugins(ctx context.Context) {

	var page int64
	var loadList []*m.Plugin
	const perPage = 500
	var err error

	if GoPluginsEnabled {
		p.ExternalPlugins.loadExternalPlugins()
	}

LOOP:
	loadList, _, err = p.adaptors.Plugin.List(context.Background(), perPage, perPage*page, "", "", pkgCommon.Bool(true), nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, pl := range loadList {
		go func(pl *m.Plugin) {
			if err = p.loadPlugin(ctx, pl.Name, pl.External); err != nil {
				log.Errorf("plugin name '%s', %s", pl.Name, err.Error())
			}
		}(pl)
	}

	if len(loadList) != 0 {
		page++
		goto LOOP
	}

	log.Info("all plugins loaded ...")
}

func (p *pluginManager) loadPlugin(ctx context.Context, name string, ext bool) (err error) {

	if p.PluginIsLoaded(name) {
		err = fmt.Errorf("%s: %w", name, apperr.ErrPluginIsLoaded)
		return
	}

	if ext && GoPluginsEnabled {
		if err = p.ExternalPlugins.loadGoPlugin(name); err != nil {
			err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginLoadExternal)
			return
		}
	}

	item, ok := IsPluginRegistered(name)
	if !ok {
		err = apperr.ErrNotFound
		return
	}

	plugin := item.(plugins.Pluggable)
	log.Infof("load plugin '%v'", plugin.Name())
	if err = plugin.Load(ctx, p.service); err != nil {
		err = fmt.Errorf("load plugin: %w", err)
		return
	}

	p.enabledPlugins.Store(name, true)

	p.pluginsWg.Add(1)

	p.eventBus.Publish("system/plugins/"+name, events.EventPluginLoaded{
		PluginName: name,
	})

	return
}

func (p *pluginManager) unloadPlugin(ctx context.Context, name string) (err error) {

	if !p.PluginIsLoaded(name) {
		err = fmt.Errorf("%s: %w", name, apperr.ErrPluginNotLoaded)
		return
	}

	if item, ok := IsPluginRegistered(name); ok {
		plugin := item.(plugins.Pluggable)
		log.Infof("unload plugin %v", plugin.Name())
		_ = plugin.Unload(ctx)
	} else {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrNotFound)
	}

	p.enabledPlugins.Store(name, false)

	p.pluginsWg.Done()

	p.eventBus.Publish("system/plugins/+", events.EventPluginUnloaded{
		PluginName: name,
	})

	return
}

// EnablePlugin ...
func (p *pluginManager) EnablePlugin(ctx context.Context, name string) (err error) {
	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, name); err != nil {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrPluginGet)
		return
	}

	if err = p.loadPlugin(ctx, name, plugin.External); err != nil {
		return
	}
	if _, ok := IsPluginRegistered(name); !ok {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrNotFound)
		return
	}
	plugin.Enabled = true
	if err = p.adaptors.Plugin.CreateOrUpdate(ctx, plugin); err != nil {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrPluginUpdate)
	}
	return
}

// DisablePlugin ...
func (p *pluginManager) DisablePlugin(ctx context.Context, name string) (err error) {
	if err = p.unloadPlugin(ctx, name); err != nil {
		return
	}
	if _, ok := IsPluginRegistered(name); !ok {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrNotFound)
		return
	}
	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(context.Background(), name); err != nil {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrPluginGet)
		return
	}
	plugin.Enabled = false
	if err = p.adaptors.Plugin.CreateOrUpdate(context.Background(), plugin); err != nil {
		err = fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrPluginUpdate)
	}
	return
}

func (p *pluginManager) PluginIsLoaded(name string) (loaded bool) {
	if value, ok := p.enabledPlugins.Load(name); ok {
		loaded = value.(bool)
	}
	return
}

func (p *pluginManager) GetPluginReadme(ctx context.Context, name string, note *string, lang *string) (result []byte, err error) {
	var plugin plugins.Pluggable
	plugin, err = p.getPlugin(name)
	if err != nil {
		return
	}
	result, err = plugin.Readme(note, lang)
	return
}

func (p *pluginManager) RemovePlugin(ctx context.Context, name string) error {

	if !GoPluginsEnabled {
		return errors.New("method not implemented")
	}

	if _, err := p.adaptors.Plugin.GetByName(context.Background(), name); err != nil {
		return fmt.Errorf("%s: %w", fmt.Sprintf("name %s", name), apperr.ErrPluginGet)
	}

	if err := p.unloadPlugin(ctx, name); err != nil {
		log.Warn(err.Error())
	}

	if err := p.ExternalPlugins.loadGoPlugin(name); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginLoadExternal)
		log.Warn(err.Error())
	}

	plugin, ok := IsPluginRegistered(name)
	if ok {
		installable, ok := plugin.(plugins.Installable)
		if ok {
			if err := installable.Uninstall(); err != nil {
				log.Warn(err.Error())
			}
		}
	}

	return p.removeExternalPlugin(ctx, name)
}

func (p *pluginManager) UploadPlugin(ctx context.Context, reader *bufio.Reader) (newPlugin *m.Plugin, err error) {

	if !GoPluginsEnabled {
		return nil, errors.New("method not implemented")
	}

	newPlugin, err = p.ExternalPlugins.uploadPlugin(ctx, reader)
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginUpload)
		return
	}

	item, ok := IsPluginRegistered(newPlugin.Name)
	if !ok {
		err = fmt.Errorf("%s: %w", "it looks like the plugin is loaded, but it didn't work to connect", apperr.ErrPluginUpload)
		return
	}

	installable, ok := item.(plugins.Installable)
	if !ok {
		return
	}

	if err = installable.Install(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginUpload)
		return
	}

	return
}
