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

package vosk

import (
	"context"
	"embed"

	"github.com/e154/bus"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.vosk")
)

var _ supervisor.Pluggable = (*plugin)(nil)

//go:embed *.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	registrar triggers.IRegistrar
	trigger   *Trigger
	stt       STT
	settings  m.Attributes
	msgQueue  bus.Bus
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin:   supervisor.NewPlugin(),
		msgQueue: bus.NewBus(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	// load settings
	p.settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		p.settings = NewSettings()
	}

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			p.trigger = NewTrigger(p.msgQueue)
			if err = p.registrar.RegisterTrigger(p.trigger); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}

	var sttModel = defaultModel
	if p.settings[AttrModel] != nil && p.settings[AttrModel].String() != "" {
		sttModel = p.settings[AttrModel].String()
	}

	p.stt = NewVosk(modelPath, sttModel, p.Service.Crawler())
	p.stt.Start()

	_ = p.Service.EventBus().Subscribe("system/stt", p.eventHandler)
	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	_ = p.Service.EventBus().Unsubscribe("system/stt", p.eventHandler)
	err = p.Plugin.Unload(ctx)

	p.trigger.Shutdown()

	if err = p.registrar.UnregisterTrigger(Name); err != nil {
		log.Error(err.Error())
		return err
	}

	p.stt.Shutdown()

	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	return
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.CommandSTT:
		text, err := p.stt.STT(v.Payload, false)
		if err != nil {
			log.Error(err.Error())
			return
		}
		p.msgQueue.Publish("/", text)
	}
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return Version
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers: true,
		Setts:    NewSettings(),
		Javascript: m.PluginOptionsJs{
			Methods:   map[string]string{},
			Variables: nil,
		},
		TriggerParams: NewTriggerParams(),
	}
}
