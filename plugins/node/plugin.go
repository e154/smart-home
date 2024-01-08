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

package node

import (
	"context"
	"embed"

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

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin: supervisor.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()
	_ = p.mqttServ.Authenticator().Register(p.Authenticator)

	p.mqttClient = p.mqttServ.NewClient("plugins.node")

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.node")
	_ = p.mqttServ.Authenticator().Unregister(p.Authenticator)

	return nil
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor = NewActor(entity, p.Service, p.mqttClient)
	return
}

// Name ...
func (p *plugin) Name() string {
	return Name
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

	var exist = false
	p.Actors.Range(func(key, value any) bool {
		actor := value.(*Actor)
		attrs := actor.Settings()

		if _login, ok := attrs[AttrNodeLogin]; !ok || _login.String() != login {
			exist = true
			return true
		}

		if _password, ok := attrs[AttrNodePass]; !ok || _password.String() != password {
			exist = true
			return true
		}

		err = nil
		return false

		// todo add encripted password
		//if ok := common.CheckPasswordHash(password, settings[AttrNodePass].String()); ok {
		//	return
		//}
	})

	if !exist {
		err = apperr.ErrBadLoginOrPassword
	}

	return
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:      true,
		ActorAttrs:  NewAttr(),
		ActorStates: supervisor.ToEntityStateShort(NewStates()),
		ActorSetts:  NewSettings(),
	}
}
