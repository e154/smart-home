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

package alexa

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.server")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	server     IServer
	actorsLock *sync.Mutex
	registrar  triggers.IRegistrar
}

func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			if err = p.registrar.RegisterTrigger(NewTrigger(p.EventBus)); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}

	// run server
	p.server = NewServer(p.Adaptors,
		NewConfig(service.AppConfig()),
		p.ScriptService,
		service.GateClient(),
		p.EventBus)

	p.server.Start()

	p.EventBus.Subscribe(TopicPluginAlexa, p.eventHandler)

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.EventBus.Unsubscribe(TopicPluginAlexa, p.eventHandler)

	p.server.Stop()
	p.server = nil

	if err = p.registrar.UnregisterTrigger(TriggerName); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	return
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return []string{"triggers"}
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Server() IServer {
	return p.server
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case EventAlexaAddSkill:
		p.server.AddSkill(v.Skill)
	case EventAlexaUpdateSkill:
		p.server.UpdateSkill(v.Skill)
	case EventAlexaDeleteSkill:
		p.server.DeleteSkill(v.Skill)
	}
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:   false,
		Triggers: true,
	}
}
