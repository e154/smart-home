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

package triggers

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"go.uber.org/atomic"
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
	isStarted *atomic.Bool
	triggers  map[string]ITrigger
	bus       event_bus.EventBus
}

func New() plugins.Plugable {
	return &plugin{
		isStarted: atomic.NewBool(false),
		triggers:  make(map[string]ITrigger),
	}
}

func (u *plugin) Load(service plugins.Service) error {
	u.bus = service.EventBus()
	return nil
}

func (u *plugin) Start() (err error) {

	if u.isStarted.Load() {
		return
	}

	u.attachTrigger()

	u.isStarted.Store(true)

	return
}

func (u *plugin) Unload() (err error) {

	return
}

func (u plugin) Name() string {
	return Name
}

func (u *plugin) attachTrigger() {

	// init triggers ...
	u.triggers[StateChangeName] = NewStateChangedTrigger(u.bus)
	u.triggers[SystemName] = NewSystemTrigger(u.bus)
	u.triggers[TimeName] = NewTimeTrigger(u.bus)

	wg := &sync.WaitGroup{}

	for _, tr := range u.triggers {
		wg.Add(1)
		go func(tr ITrigger, wg *sync.WaitGroup) {
			log.Infof("attach trigger '%s'", tr.Name())
			tr.AsyncAttach(wg)
		}(tr, wg)
	}

	wg.Wait()
}

func (u *plugin) GetTrigger(name string) (trigger ITrigger, err error) {
	var ok bool
	if trigger, ok = u.triggers[name]; !ok {
		err = fmt.Errorf("not found trigger with name(%s)", name)
	}
	return
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Version() string {
	return "0.0.1"
}
