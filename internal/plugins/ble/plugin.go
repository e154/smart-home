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

package ble

import (
	"context"
	"embed"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/plugins/triggers"
)

var (
	log = logger.MustGetLogger("plugins.ble")
)

var _ plugins.Pluggable = (*plugin)(nil)

//go:embed *.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	registrar triggers.IRegistrar
	trigger   *Trigger
}

// New ...
func New() plugins.Pluggable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			p.trigger = NewTrigger(p.Service.EventBus())
			if err = p.registrar.RegisterTrigger(p.trigger); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)
	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)
	err = p.Plugin.Unload(ctx)

	p.trigger.Shutdown()
	if err = p.registrar.UnregisterTrigger(Name); err != nil {
		log.Error(err.Error())
		return err
	}

	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor plugins.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	return
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
	case events.EventCallEntityAction:
		values, ok := p.Check(v)
		if !ok {
			return
		}
		for _, value := range values {
			actor := value.(*Actor)
			actor.addAction(v)
		}
	}
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers:           true,
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorAttrs:         nil,
		ActorCustomActions: true,
		ActorActions:       plugins.ToEntityActionShort(NewActions()),
		ActorCustomStates:  true,
		ActorStates:        nil,
		ActorCustomSetts:   true,
		ActorSetts:         NewSettings(),
		Setts:              nil,
		Javascript: m.PluginOptionsJs{
			Methods: map[string]string{
				"BleRead":  "",
				"BleWrite": "",
			},
			Variables: nil,
		},
		TriggerParams: NewTriggerParams(),
	}
}
