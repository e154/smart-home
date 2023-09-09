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
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.node")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[common.EntityId]*Actor
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()

	p.mqttClient = p.mqttServ.NewClient("plugins.node")
	_ = p.mqttServ.Authenticator().Register(p.Authenticator)

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.node")
	_ = p.mqttServ.Authenticator().Unregister(p.Authenticator)

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	// remove actors
	for entityId, actor := range p.actors {
		actor.destroy()
		delete(p.actors, entityId)
	}

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		return
	}

	actor := NewActor(entity, p.Supervisor, p.Adaptors, p.ScriptService, p.EventBus, p.mqttClient)
	p.actors[entity.Id] = actor
	p.Supervisor.Spawn(actor.Spawn)

	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	actor, ok := p.actors[entityId]
	if !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId))
		return
	}

	actor.destroy()

	delete(p.actors, entityId)

	return
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
	return "0.0.1"
}

func (p *plugin) pushToNode() {

}

// Authenticator ...
func (p *plugin) Authenticator(login, password string) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		attrs := actor.Settings()

		if _login, ok := attrs[AttrNodeLogin]; !ok || _login.String() != login {
			continue
		}

		if _password, ok := attrs[AttrNodePass]; !ok || _password.String() != password {
			continue
		}

		err = nil
		return

		// todo add encripted password
		//if ok := common.CheckPasswordHash(password, settings[AttrNodePass].String()); ok {
		//	return
		//}
	}

	err = apperr.ErrBadLoginOrPassword

	return
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers:           false,
		Actors:             true,
		ActorCustomAttrs:   false,
		ActorAttrs:         NewAttr(),
		ActorCustomActions: false,
		ActorActions:       nil,
		ActorCustomStates:  false,
		ActorStates:        supervisor.ToEntityStateShort(NewStates()),
		ActorCustomSetts:   false,
		ActorSetts:         NewSettings(),
		Setts:              nil,
	}
}
