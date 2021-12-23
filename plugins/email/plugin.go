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

package email

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = common.MustGetLogger("plugins.email")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	notify notify.ProviderRegistrar
	actor  *Actor
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	go func() {
		if err = p.asyncLoad(); err != nil {
			log.Error(err.Error())
		}
	}()

	return nil
}

func (p *plugin) asyncLoad() (err error) {

	// load settings
	var settings m.Attributes
	settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		settings = NewSettings()
	}

	if settings == nil {
		settings = NewSettings()
	}

	// get provider registrar
	var pl interface{}
	if pl, err = p.GetPlugin(notify.Name); err != nil {
		return
	}

	var ok bool
	p.notify, ok = pl.(notify.ProviderRegistrar)
	if !ok {
		err = errors.Wrap(common.ErrInternal, "can`t static cast to  notify.ProviderRegistrar")
		return
	}

	// add actor
	p.actor = NewActor(settings, p.EntityManager, p.EventBus, p.Adaptors)
	p.EntityManager.Spawn(p.actor.Spawn)

	// register email provider
	p.notify.AddProvider(Name, p)

	return
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	if p.notify == nil {
		return
	}
	p.notify.RemoveProvider(Name)

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{notify.Name}
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Setts: NewSettings(),
	}
}

// Save ...
func (p *plugin) Save(msg notify.Message) (addresses []string, message m.Message) {
	message = m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = p.Adaptors.Message.Add(message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrAddresses].String(), ",")
	return
}

// Send ...
func (p *plugin) Send(address string, message m.Message) (err error) {
	err = p.actor.Send(address, message)
	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}
