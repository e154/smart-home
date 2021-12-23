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
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
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

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.RWMutex{},
		actors:     make(map[common.EntityId]*Actor),
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

	// get provider registrar
	var pl interface{}
	if pl, err = p.GetPlugin(notify.Name); err != nil {
		return
	}

	var ok bool
	p.notify, ok = pl.(notify.ProviderRegistrar)
	if !ok {
		err = errors.Wrap(common.ErrInternal, "can`t static cast to notify.ProviderRegistrar")
		return
	}

	// register telegram provider
	p.notify.AddProvider(Name, p)

	p.EventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)

	return
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.EventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

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
		Actors:             true,
		ActorCustomActions: true,
		ActorCustomStates:  true,
		ActorCustomAttrs:   true,
		ActorAttrs:         NewAttr(),
		ActorSetts:         NewSettings(),
	}
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, p.EntityManager, p.ScriptService, p.EventBus, p.Adaptors); err != nil {
		return
	}
	p.actors[entity.Id] = actor
	p.EntityManager.Spawn(actor.Spawn)
	actor.Start()
	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId]; !ok {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("entityId \"%s\"", entityId))
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

	actor, ok := p.actors[common.EntityId(fmt.Sprintf("telegram.%s", name))]
	if ok {
		actor.Send(body)
	} else {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("bot \"%s\"", name))
	}

	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventStateChanged:
	case event_bus.EventCallAction:
		actor, ok := p.actors[v.EntityId]
		if !ok {
			return
		}
		actor.addAction(v)
	}

	return
}
