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

package plugin_manager

import (
	"container/list"
	"context"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.manager")
)

type pluginManager struct {
	adaptors  *adaptors.Adaptors
	isStarted *atomic.Bool
	loadLock  *sync.Mutex
	plugins   *list.List
}

func NewPluginManager(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors, ) (manager PluginManager) {
	manager = &pluginManager{
		adaptors:  adaptors,
		isStarted: atomic.NewBool(false),
		loadLock:  &sync.Mutex{},
		plugins:   list.New(),
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (err error) {
			manager.Shutdown()
			return nil
		},
	})
	return
}

func (p *pluginManager) Register(plugin Plugable) {
	p.loadLock.Lock()
	defer p.loadLock.Unlock()
	p.plugins.PushBack(pluginListItem{
		Name:   plugin.Name(),
		Plugin: plugin,
	})
}

func (p *pluginManager) Start() {
	if p.isStarted.Load() {
		return
	}
	p.isStarted.Store(true)

	if err := p.loadPlugins(); err != nil {
		log.Error(err.Error())
	}
}

func (p *pluginManager) Shutdown() {
	log.Info("Shutdown")

	if !p.isStarted.Load() {
		return
	}
	p.isStarted.Store(false)

	if err := p.unloadPlugins(); err != nil {
		log.Error(err.Error())
	}
}

func (p *pluginManager) loadPlugins() (err error) {
	p.loadLock.Lock()
	defer p.loadLock.Unlock()

	log.Info("Start")
	for e := p.plugins.Front(); e != nil; e = e.Next() {
		i, ok := e.Value.(pluginListItem)
		if !ok {
			continue
		}
		log.Infof("load plugin %v", i.Name)
		plugins := make(map[string]interface{})
		if len(i.Plugin.Depends()) > 0 {
			var plugin Plugable
			for _, dep := range i.Plugin.Depends() {
				if plugin, err = p.unsafeGetPlugin(dep); err != nil {
					log.Error(err.Error())
					continue
				}
				plugins[dep] = plugin
			}
		}
		if err := i.Plugin.Load(p, plugins); err != nil {
			log.Error(err.Error())
		}
	}

	log.Info("loaded ...")

	return
}

func (p *pluginManager) unloadPlugins() (err error) {
	p.loadLock.Lock()
	defer p.loadLock.Unlock()

	for e := p.plugins.Front(); e != nil; e = e.Next() {
		i, ok := e.Value.(pluginListItem)
		if !ok {
			continue
		}

		log.Infof("unload plugin %v", i.Name)
		if err := i.Plugin.Unload(); err != nil {
			log.Error(err.Error())
		}
	}
	for e := p.plugins.Front(); e != nil; e = e.Next() {
		p.plugins.Remove(e)
	}
	p.plugins = list.New()
	return
}

func (p *pluginManager) GetPlugin(t string) (plugin Plugable, err error) {
	p.loadLock.Lock()
	defer p.loadLock.Unlock()

	plugin, err = p.unsafeGetPlugin(t)

	return
}

func (p *pluginManager) unsafeGetPlugin(t string) (plugin Plugable, err error) {

	for e := p.plugins.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(pluginListItem)
		if !ok {
			continue
		}
		if p.Name == t {
			plugin = p.Plugin
			return
		}
	}

	err = fmt.Errorf("plugin with name '%v' not found", t)

	return
}
