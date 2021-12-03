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

package triggers

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.triggers")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	mu       *sync.Mutex
	triggers map[string]ITrigger
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin:   plugins.NewPlugin(),
		mu:       &sync.Mutex{},
		triggers: make(map[string]ITrigger),
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.attachTrigger()

	return
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}
	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) attachTrigger() {

	p.mu.Lock()
	defer p.mu.Unlock()

	// init triggers ...
	p.triggers[StateChangeName] = NewStateChangedTrigger(p.EventBus)
	p.triggers[SystemName] = NewSystemTrigger(p.EventBus)
	p.triggers[TimeName] = NewTimeTrigger(p.EventBus)

	wg := &sync.WaitGroup{}

	for _, tr := range p.triggers {
		wg.Add(1)
		go func(tr ITrigger, wg *sync.WaitGroup) {
			log.Infof("register trigger '%s'", tr.Name())
			tr.AsyncAttach(wg)
		}(tr, wg)
	}

	wg.Wait()
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// GetTrigger ...
func (p *plugin) GetTrigger(name string) (trigger ITrigger, err error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	var ok bool
	if trigger, ok = p.triggers[name]; !ok {
		err = fmt.Errorf("not found trigger with name(%s)", name)
	}
	return
}

// RegisterTrigger ...
func (p *plugin) RegisterTrigger(tr ITrigger) (err error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.triggers[tr.Name()]; ok {
		err = fmt.Errorf("trigger with name %s is registerred", tr.Name())
		return
	}

	p.triggers[tr.Name()] = tr
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Infof("register trigger '%s'", tr.Name())
	go tr.AsyncAttach(wg)
	wg.Wait()

	return
}

// UnregisterTrigger ...
func (p *plugin) UnregisterTrigger(name string) error {

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.triggers[name]; ok {
		delete(p.triggers, name)
		return nil
	}

	return nil
}

// TriggerList ...
func (p *plugin) TriggerList() (list []string) {

	p.mu.Lock()
	defer p.mu.Unlock()

	list = make([]string, 0, len(p.triggers))
	for name, _ := range p.triggers {
		list = append(list, name)
	}
	return
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers: true,
	}
}
