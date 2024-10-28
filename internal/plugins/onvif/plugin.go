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
	"time"

	web2 "github.com/e154/smart-home/internal/common/web"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/web"
)

var (
	log = logger.MustGetLogger("plugins.onvif")
)

var _ plugins.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
}

// New ...
func New() plugins.Pluggable {
	p := &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
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
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorAttrs:         NewAttr(),
		ActorCustomActions: true,
		ActorActions:       plugins.ToEntityActionShort(NewActions()),
		ActorStates:        plugins.ToEntityStateShort(NewStates()),
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

	crawler := web2.New()
	crawler.DigestAuth(actor.Setts[AttrUserName].String(),
		actor.Setts[AttrPassword].Decrypt())

	var err error
	filePath, err = crawler.Download(web.Request{Method: "GET", Url: actor.GetSnapshotUri(), Timeout: time.Second * 2})
	if err != nil {
		log.Error(err.Error())
	}

	return
}
