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
	"context"
	"fmt"
	"github.com/e154/smart-home/common"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/atomic"

	"github.com/e154/bus"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
)

type pluginManager struct {
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

func (p *pluginManager) getPlugin(name string) (plugin Pluggable, err error) {

	if item, ok := pluginList.Load(name); ok {
		plugin = item.(Pluggable)
		return
	}

	err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("name %s", name))

	return
}

func (p *pluginManager) loadPlugins(ctx context.Context) {

	var page int64
	var loadList []*m.Plugin
	const perPage = 500
	var err error

LOOP:
	loadList, _, err = p.adaptors.Plugin.List(context.Background(), perPage, perPage*page, "", "", common.Bool(true), nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, pl := range loadList {
		go func(pl *m.Plugin) {
			if err = p.loadPlugin(ctx, pl.Name); err != nil {
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

func (p *pluginManager) loadPlugin(ctx context.Context, name string) (err error) {

	if p.PluginIsLoaded(name) {
		err = errors.Wrap(apperr.ErrPluginIsLoaded, name)
		return
	}
	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Pluggable)
		log.Infof("load plugin '%v'", plugin.Name())
		if err = plugin.Load(ctx, p.service); err != nil {
			err = errors.Wrap(err, "load plugin")
			return
		}
	} else {
		err = apperr.ErrNotFound
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
		err = errors.Wrap(apperr.ErrPluginNotLoaded, name)
		return
	}

	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Pluggable)
		log.Infof("unload plugin %v", plugin.Name())
		_ = plugin.Unload(ctx)
	} else {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("name %s", name))
	}

	p.enabledPlugins.Store(name, false)

	p.pluginsWg.Done()

	p.eventBus.Publish("system/plugins/+", events.EventPluginUnloaded{
		PluginName: string(name),
	})

	return
}

// Install ...
func (p *pluginManager) Install(ctx context.Context, t string) {

	pl, _ := p.adaptors.Plugin.GetByName(context.Background(), t)
	if pl.Enabled {
		return
	}

	plugin, err := p.getPlugin(t)
	if err != nil {
		return
	}

	if plugin.Type() != PluginInstallable {
		return
	}

	installable, ok := plugin.(Installable)
	if !ok {
		return
	}

	if err := installable.Install(); err != nil {
		log.Error(err.Error())
		return
	}

	_ = p.adaptors.Plugin.CreateOrUpdate(context.Background(), &m.Plugin{
		Name:    plugin.Name(),
		Version: plugin.Version(),
		Enabled: true,
		System:  plugin.Type() == PluginBuiltIn,
	})

	if err = p.loadPlugin(ctx, plugin.Name()); err != nil {
		log.Error(err.Error())
	}
}

// Uninstall ...
func (p *pluginManager) Uninstall(name string) {

}

// EnablePlugin ...
func (p *pluginManager) EnablePlugin(ctx context.Context, name string) (err error) {
	if err = p.loadPlugin(ctx, name); err != nil {
		return
	}
	if _, ok := pluginList.Load(name); !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("name %s", name))
		return
	}
	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(context.Background(), name); err != nil {
		err = errors.Wrap(apperr.ErrPluginGet, fmt.Sprintf("name %s", name))
		return
	}
	plugin.Enabled = true
	if err = p.adaptors.Plugin.CreateOrUpdate(context.Background(), plugin); err != nil {
		err = errors.Wrap(apperr.ErrPluginUpdate, fmt.Sprintf("name %s", name))
	}
	return
}

// DisablePlugin ...
func (p *pluginManager) DisablePlugin(ctx context.Context, name string) (err error) {
	if err = p.unloadPlugin(ctx, name); err != nil {
		return
	}
	if _, ok := pluginList.Load(name); !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("name %s", name))
		return
	}
	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(context.Background(), name); err != nil {
		err = errors.Wrap(apperr.ErrPluginGet, fmt.Sprintf("name %s", name))
		return
	}
	plugin.Enabled = false
	if err = p.adaptors.Plugin.CreateOrUpdate(context.Background(), plugin); err != nil {
		err = errors.Wrap(apperr.ErrPluginUpdate, fmt.Sprintf("name %s", name))
	}
	return
}

// PluginList ...
func (p *pluginManager) PluginList() (list []PluginInfo, total int64, err error) {

	list = make([]PluginInfo, 0)
	pluginList.Range(func(key, value interface{}) bool {
		total++
		plugin := value.(Pluggable)
		list = append(list, PluginInfo{
			Name:    plugin.Name(),
			Version: plugin.Version(),
			Enabled: p.PluginIsLoaded(plugin.Name()),
			System:  plugin.Type() == PluginBuiltIn,
		})
		return true
	})
	return
}

func (p *pluginManager) PluginIsLoaded(name string) (loaded bool) {
	if value, ok := p.enabledPlugins.Load(name); ok {
		loaded = value.(bool)
	}
	return
}

func (p *pluginManager) GetPluginReadme(ctx context.Context, name string, note *string, lang *string) (result []byte, err error) {
	var plugin Pluggable
	plugin, err = p.getPlugin(name)
	if err != nil {
		return
	}
	result, err = plugin.Readme(note, lang)
	return
}
