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

package telegram

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.telegram")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	notify     notify.ProviderRegistrar
	actorsLock *sync.RWMutex
	actors     map[common.EntityId]*Actor
}

func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.RWMutex{},
		actors:     make(map[common.EntityId]*Actor),
	}
}

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

	// get provider registrar
	var pl interface{}
	if pl, err = p.GetPlugin(notify.Name); err != nil {
		return
	}

	var ok bool
	p.notify, ok = pl.(notify.ProviderRegistrar)
	if !ok {
		err = errors.New("fail static cast to notify.ProviderRegistrar")
		return
	}

	// register telegram provider
	p.notify.AddProvider(Name, p)

	return
}

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

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return []string{notify.Name}
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:     true,
		ActorAttrs: NewAttr(),
		ActorSetts: NewSettings(),
	}
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		err = fmt.Errorf("the actor with id '%s' has already been created", entity.Id)
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, p.EntityManager, p.EventBus, p.Adaptors); err != nil {
		return
	}
	p.actors[entity.Id] = actor
	p.EntityManager.Spawn(actor.Spawn)
	actor.Start()
	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId]; !ok {
		err = fmt.Errorf("not found")
		return
	}

	p.actors[entityId].Stop()
	delete(p.actors, entityId)
	return
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

	addresses = []string{attr[AttrName].String()}

	return
}

// Send ...
func (p *plugin) Send(name string, message m.Message) (err error) {
	params := NewMessageParams()
	params.Deserialize(message.Attributes)

	body := params[AttrBody].String()

	p.actorsLock.RLock()
	defer p.actorsLock.RUnlock()

	if actor, ok := p.actors[common.EntityId(fmt.Sprintf("telegram.%s", name))]; ok {
		actor.Send(body)
	}

	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}
