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

package plugins

import (
	"context"
	"fmt"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"go.uber.org/fx"
)

var (
	log = common.MustGetLogger("plugins.manager")
)

type pluginManager struct {
	adaptors       *adaptors.Adaptors
	isStarted      *atomic.Bool
	service        *service
	eventBus       event_bus.EventBus
	enabledPlugins map[string]bool
}

// NewPluginManager ...
func NewPluginManager(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	bus event_bus.EventBus,
	entityManager entity_manager.EntityManager,
	mqttServ mqtt.MqttServ,
	scriptService scripts.ScriptService,
	appConfig *m.AppConfig,
	gateClient *gate_client.GateClient,
	eventBus event_bus.EventBus) common.PluginManager {
	pluginManager := &pluginManager{
		adaptors:       adaptors,
		isStarted:      atomic.NewBool(false),
		eventBus:       eventBus,
		enabledPlugins: make(map[string]bool),
	}
	pluginManager.service = &service{
		pluginManager: pluginManager,
		bus:           bus,
		entityManager: entityManager,
		mqttServ:      mqttServ,
		adaptors:      adaptors,
		scriptService: scriptService,
		appConfig:     appConfig,
		gateClient:    gateClient,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (err error) {
			pluginManager.Shutdown()
			return nil
		},
	})
	return pluginManager
}

// Start ...
func (p *pluginManager) Start() {
	if p.isStarted.Load() {
		return
	}
	p.isStarted.Store(true)

	log.Info("Start")

	p.loadPlugins()
}

// Shutdown ...
func (p *pluginManager) Shutdown() {

	if !p.isStarted.Load() {
		return
	}
	p.isStarted.Store(false)

	for name, ok := range p.enabledPlugins {
		if !ok {
			continue
		}
		log.Infof("unload plugin '%s'", name)
		if item, ok := pluginList.Load(name); ok {
			plugin := item.(Plugable)
			plugin.Unload()
		}
		p.enabledPlugins[name] = false
	}

	log.Info("Shutdown")
}

// GetPlugin ...
func (p *pluginManager) GetPlugin(t string) (plugin interface{}, err error) {

	plugin, err = p.getPlugin(t)

	return
}

func (p *pluginManager) getPlugin(name string) (plugin Plugable, err error) {

	if item, ok := pluginList.Load(name); ok {
		plugin = item.(Plugable)
		return
	}

	err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("name %s", name))

	return
}

func (p *pluginManager) loadPlugins() {

	var page int64
	var loadList []m.Plugin
	const perPage = 100
	var err error

LOOP:
	loadList, _, err = p.adaptors.Plugin.List(perPage, perPage*page, "", "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, pl := range loadList {
		if !pl.Enabled {
			continue
		}
		if err = p.loadPlugin(pl.Name); err != nil {
			log.Errorf("plugin name '%s', %s", pl.Name, err.Error())
		}
	}

	if len(loadList) != 0 {
		page++
		goto LOOP
	}

	log.Info("all plugins loaded ...")
}

func (p *pluginManager) loadPlugin(name string) (err error) {

	if p.enabledPlugins[name] {
		err = errors.Wrap(ErrPluginIsLoaded, name)
		return
	}
	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Plugable)
		log.Infof("load plugin %v", plugin.Name())
		if err = plugin.Load(p.service); err != nil {
			err = errors.Wrap(err, "load plugin")
			return
		}
	} else {
		err = common.ErrNotFound
	}

	p.enabledPlugins[name] = true

	p.eventBus.Publish(event_bus.TopicPlugins, event_bus.EventLoadedPlugin{
		PluginName: string(name),
	})

	return
}

func (p *pluginManager) unloadPlugin(name string) (err error) {

	if !p.enabledPlugins[name] {
		err = errors.Wrap(ErrPluginNotLoaded, name)
		return
	}

	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Plugable)
		log.Infof("unload plugin %v", plugin.Name())
		plugin.Unload()
	} else {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("name %s", name))
	}

	p.enabledPlugins[name] = false

	p.eventBus.Publish(event_bus.TopicPlugins, event_bus.EventUnloadedPlugin{
		PluginName: string(name),
	})

	return
}

// Install ...
func (p *pluginManager) Install(t string) {

	pl, _ := p.adaptors.Plugin.GetByName(t)
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

	p.adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    plugin.Name(),
		Version: plugin.Version(),
		Enabled: true,
		System:  plugin.Type() == PluginBuiltIn,
	})

	if err = p.loadPlugin(plugin.Name()); err != nil {
		log.Error(err.Error())
	}
}

// Uninstall ...
func (p *pluginManager) Uninstall(name string) {

}

// EnablePlugin ...
func (p *pluginManager) EnablePlugin(name string) (err error) {
	if err = p.loadPlugin(name); err != nil {
		return
	}
	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Plugable)
		err = p.adaptors.Plugin.CreateOrUpdate(m.Plugin{
			Name:    plugin.Name(),
			Version: plugin.Version(),
			Enabled: true,
			System:  plugin.Type() == PluginBuiltIn,
		})
	} else {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("name %s", name))
	}
	return
}

// DisablePlugin ...
func (p *pluginManager) DisablePlugin(name string) (err error) {
	if err = p.unloadPlugin(name); err != nil {
		return
	}
	if item, ok := pluginList.Load(name); ok {
		plugin := item.(Plugable)
		err = p.adaptors.Plugin.CreateOrUpdate(m.Plugin{
			Name:    plugin.Name(),
			Version: plugin.Version(),
			Enabled: plugin.Type() == PluginBuiltIn,
			System:  plugin.Type() == PluginBuiltIn,
		})
	} else {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("name %s", name))
	}
	return
}

// PluginList ...
func (p *pluginManager) PluginList() (list []common.PluginInfo, total int64, err error) {

	list = make([]common.PluginInfo, 0)
	pluginList.Range(func(key, value interface{}) bool {
		total++
		plugin := value.(Plugable)
		list = append(list, common.PluginInfo{
			Name:    plugin.Name(),
			Version: plugin.Version(),
			Enabled: p.enabledPlugins[plugin.Name()],
			System:  plugin.Type() == PluginBuiltIn,
		})
		return true
	})
	return
}
