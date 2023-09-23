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
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.telegram")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	// register telegram provider
	notify.ProviderManager.AddProvider(Name, p)

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)

	notify.ProviderManager.RemoveProvider(Name)

	return nil
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor, err = NewActor(entity, p.Service)
	return
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
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
		ActorStates:        supervisor.ToEntityStateShort(NewStates()),
	}
}

// Save ...
func (p *plugin) Save(msg notify.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = p.Service.Adaptors().Message.Add(context.Background(), message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = []string{attr[AttrName].String()}

	return
}

// Send ...
func (p *plugin) Send(name string, message *m.Message) (err error) {
	params := NewMessageParams()
	_, _ = params.Deserialize(message.Attributes)

	body := params[AttrBody].String()

	entityId := common.EntityId(fmt.Sprintf("telegram.%s", name))
	value, ok := p.Actors.Load(entityId)
	if !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("bot \"%s\"", name))
		return
	}
	actor := value.(*Actor)
	_ = actor.Send(body)

	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
	case events.EventCallEntityAction:
		value, ok := p.Actors.Load(v.EntityId)
		if !ok {
			return
		}
		actor := value.(*Actor)
		actor.addAction(v)
	}
}
