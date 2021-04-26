// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.manager")
)

type pluginManager struct {
	adaptors       *adaptors.Adaptors
	isStarted      *atomic.Bool
	service        *service
	loadLock       *sync.Mutex
	enabledPlugins map[string]bool
}

func NewPluginManager(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	bus event_bus.EventBus,
	entityManager entity_manager.EntityManager,
	mqttServ mqtt.MqttServ,
	scriptService scripts.ScriptService) common.PluginManager {
	pluginManager := &pluginManager{
		adaptors:       adaptors,
		isStarted:      atomic.NewBool(false),
		loadLock:       &sync.Mutex{},
		enabledPlugins: make(map[string]bool),
	}
	pluginManager.service = &service{
		pluginManager: pluginManager,
		bus:           bus,
		entityManager: entityManager,
		mqttServ:      mqttServ,
		adaptors:      adaptors,
		scriptService: scriptService,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (err error) {
			pluginManager.Shutdown()
			return nil
		},
	})
	return pluginManager
}

func (p *pluginManager) Start() {
	if p.isStarted.Load() {
		return
	}
	p.isStarted.Store(true)

	log.Info("Start")

	p.loadPlugins()
}

func (p *pluginManager) Shutdown() {
	log.Info("Shutdown")

	if !p.isStarted.Load() {
		return
	}
	p.isStarted.Store(false)

	for name, ok := range p.enabledPlugins {
		if !ok {
			continue
		}
		log.Infof("unload plugin %v", name)
		if plugin, ok := pluginList[name]; ok {
			plugin.Unload()
		}
	}
}

func (p *pluginManager) GetPlugin(t string) (plugin interface{}, err error) {
	p.loadLock.Lock()
	defer p.loadLock.Unlock()

	plugin, err = p.unsafeGetPlugin(t)

	return
}

func (p *pluginManager) unsafeGetPlugin(t string) (plugin Plugable, err error) {

	if enabled := p.enabledPlugins[t]; !enabled {
		err = fmt.Errorf("plugin '%v' disabled", t)
		return
	}

	var ok bool
	if plugin, ok = pluginList[t]; ok {
		return
	}

	err = fmt.Errorf("plugin '%v' not found", t)

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

		if plugin, ok := pluginList[pl.Name]; ok {
			log.Infof("load plugin %v", plugin.Name())
			plugin.Load(p.service)
		}

		p.enabledPlugins[pl.Name] = true
	}

	if len(loadList) != 0 {
		page++
		goto LOOP
	}

	log.Info("all plugins loaded ...")
}

func (p *pluginManager) Install(t string) {

	pl, _ := p.adaptors.Plugin.GetByName(t)
	if pl.Enabled {
		return
	}

	p.loadLock.Lock()
	defer p.loadLock.Unlock()

	plugin, err := p.unsafeGetPlugin(t)
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

	if err = plugin.Load(p.service); err != nil {
		log.Error(err.Error())
	}
}

func (p *pluginManager) Uninstall(name string) {

}
