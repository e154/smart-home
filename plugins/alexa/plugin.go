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

package alexa

import (
	"context"
	"embed"
	"sync"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.server")
)

var _ supervisor.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	server     IServer
	actorsLock *sync.Mutex
	registrar  triggers.IRegistrar
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			if err = p.registrar.RegisterTrigger(NewTrigger(p.Service.EventBus())); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}

	// run server
	p.server = NewServer(p.Service.Adaptors(),
		NewConfig(service.AppConfig()),
		p.Service.ScriptService(),
		p.Service.EventBus())

	p.server.Start()

	_ = p.Service.EventBus().Subscribe("system/models/alexa/skills/#", p.eventHandler)

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	_ = p.Service.EventBus().Unsubscribe("system/models/alexa/skills/#", p.eventHandler)

	p.server.Stop()
	p.server = nil

	if err = p.registrar.UnregisterTrigger(TriggerName); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	return
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{"triggers"}
}

// Version ...
func (p *plugin) Version() string {
	return Version
}

// Server ...
func (p *plugin) Server() IServer {
	return p.server
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case EventAddedAlexaSkillModel:
		p.server.AddSkill(v.Skill)
	case EventUpdatedAlexaSkillModel:
		p.server.UpdateSkill(v.Skill)
	case EventDeletedAlexaSkill:
		p.server.DeleteSkill(v.Skill)
	}
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:        false,
		Triggers:      true,
		TriggerParams: NewTriggerParams(),
	}
}
