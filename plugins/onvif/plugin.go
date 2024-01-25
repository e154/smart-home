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

package onvif

import (
	"context"
	"embed"
	"sync"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.onvif")
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
	actorsLock *sync.Mutex
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
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
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}
	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)
	p.Service.ScriptService().PushFunctions("OnvifGetSnapshotUri", GetSnapshotUriBind(p))
	p.Service.ScriptService().PushFunctions("DownloadSnapshot", DownloadSnapshotBind(p))
	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	p.Service.ScriptService().PopFunction("OnvifGetSnapshotUri")
	p.Service.ScriptService().PopFunction("DownloadSnapshot")
	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)
	err = p.Plugin.Unload(ctx)
	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
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
		value, ok := p.Actors.Load(v.EntityId)
		if !ok {
			return
		}
		actor := value.(*Actor)
		actor.addAction(v)
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
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorAttrs:         NewAttr(),
		ActorCustomActions: true,
		ActorActions:       supervisor.ToEntityActionShort(NewActions()),
		ActorStates:        supervisor.ToEntityStateShort(NewStates()),
		ActorSetts:         NewSettings(),
		Javascript: m.PluginOptionsJs{
			Methods: nil,
			Objects: map[string]string{
				"Camera": "",
			},
			Variables: nil,
		},
	}
}

func (p *plugin) GetSnapshotUri(entityId common.EntityId) string {
	if value, ok := p.Actors.Load(entityId); ok {
		actor := value.(*Actor)
		return actor.GetSnapshotUri()
	}
	return ""
}

// experimental method ...
func (p *plugin) DownloadSnapshotDigest(entityId common.EntityId) (filePath string) {
	value, ok := p.Actors.Load(entityId)
	if !ok {
		return
	}
	actor := value.(*Actor)

	crawler := web.New()
	crawler.DigestAuth(actor.Setts[AttrUserName].String(),
		actor.Setts[AttrPassword].Decrypt())

	var err error
	filePath, err = crawler.Download(web.Request{Method: "GET", Url: actor.GetSnapshotUri()})
	if err != nil {
		log.Error(err.Error())
	}

	return
}
