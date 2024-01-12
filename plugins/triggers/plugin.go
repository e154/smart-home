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

package triggers

import (
	"context"
	"embed"
	"fmt"
	"sync"

	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common/logger"

	"github.com/pkg/errors"

	m "github.com/e154/smart-home/models"
)

var (
	log = logger.MustGetLogger("plugins.triggers")
)

var _ supervisor.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	mu       *sync.Mutex
	triggers map[string]ITrigger
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin:   supervisor.NewPlugin(),
		mu:       &sync.Mutex{},
		triggers: make(map[string]ITrigger),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	p.attachTrigger()

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
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
	p.triggers[StateChangeName] = NewStateChangedTrigger(p.Service.EventBus())
	p.triggers[SystemName] = NewSystemTrigger(p.Service.EventBus())
	p.triggers[TimeName] = NewTimeTrigger(p.Service.EventBus(), p.Service.Scheduler())

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
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginBuiltIn
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return Version
}

// GetTrigger ...
func (p *plugin) GetTrigger(name string) (trigger ITrigger, err error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	var ok bool
	if trigger, ok = p.triggers[name]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("trigger name '%s'", name))
	}
	return
}

// RegisterTrigger ...
func (p *plugin) RegisterTrigger(tr ITrigger) (err error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.triggers[tr.Name()]; ok {
		err = errors.Wrap(apperr.ErrInternal, fmt.Sprintf("trigger '%s' is registerred", tr.Name()))
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
	for name := range p.triggers {
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
