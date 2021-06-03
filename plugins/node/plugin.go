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

package node

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.node")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[common.EntityId]*Actor
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
}

func New() plugins.Plugable {
	return &plugin{
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()

	p.mqttClient = p.mqttServ.NewClient("plugins.node")
	p.mqttServ.Authenticator().Register(p.Authenticator)

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.node")
	p.mqttServ.Authenticator().Unregister(p.Authenticator)

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	// remove actors
	for entityId, actor := range p.actors {
		actor.destroy()
		delete(p.actors, entityId)
	}

	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		err = fmt.Errorf("the actor with id '%s' has already been created", entity.Id)
		return
	}

	var actor *Actor
	actor = NewActor(entity, p.EntityManager, p.Adaptors, p.ScriptService, p.EventBus, p.mqttClient)
	p.actors[entity.Id] = actor
	p.EntityManager.Spawn(actor.Spawn)

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	actor, ok := p.actors[entityId]
	if !ok {
		err = fmt.Errorf("not found")
		return
	}

	actor.destroy()

	delete(p.actors, entityId)

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

func (p *plugin) pushToNode() {

}

func (p *plugin) Authenticator(login, password string) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		attrs := actor.Settings()

		if attrs[AttrNodeLogin].String() != login {
			continue
		}

		if attrs[AttrNodePass].String() != password {
			continue
		}

		err = nil
		return

		// todo add encripted password
		//if ok := common.CheckPasswordHash(password, settings[AttrNodePass].String()); ok {
		//	return
		//}
	}

	err = mqtt_authenticator.ErrBadLoginOrPassword

	return
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:      true,
		ActorAttrs:  NewAttr(),
		ActorSetts:  NewSettings(),
		ActorStates: entity_manager.ToEntityStateShort(NewStates()),
	}
}
