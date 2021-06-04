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

package notify

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = common.MustGetLogger("plugins.notify")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	notify Notify
}


func New() plugins.Plugable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.notify = NewNotify(p.Adaptors, p.ScriptService)
	p.notify.Start()

	p.EventBus.Subscribe(TopicNotify, p.eventHandler)

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.EventBus.Unsubscribe(TopicNotify, p.eventHandler)

	p.notify.Shutdown()

	return nil
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case Message:
		p.notify.Send(v)
	}

	return
}

func (p *plugin) Name() string {
	return Name
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

// AddProvider ...
func (p *plugin) AddProvider(name string, provider Provider) {
	p.notify.AddProvider(name, provider)
}

// RemoveProvider ...
func (p *plugin) RemoveProvider(name string) {
	p.notify.RemoveProvider(name)
}

func (p *plugin) Provider(name string) (provider Provider, err error) {
	panic("implement me")
}
