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

package pachka

import (
	"context"
	"embed"
	"fmt"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/notify/common"
	"github.com/e154/smart-home/system/supervisor"
	"strconv"
)

var (
	log = logger.MustGetLogger("plugins.pachka")
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
	client *Client
	notify *notify.Notify
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
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	p.notify = notify.NewNotify(service.Adaptors())
	p.notify.Start()

	// load settings
	var settings m.Attributes
	settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		settings = NewSettings()
	}

	if settings == nil {
		settings = NewSettings()
	}

	if settings[AttrToken].Decrypt() != "" {
		p.client = NewClient(settings[AttrToken].Decrypt())
	} else {
		log.Warn("empty access token")
		p.client = NewClient("NoToken")
	}

	_ = p.Service.EventBus().Subscribe(notify.TopicNotify, p.eventHandler, false)

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.notify.Shutdown()

	_ = p.Service.EventBus().Unsubscribe(notify.TopicNotify, p.eventHandler)

	return nil
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
	return Version
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:     true,
		ActorAttrs: NewAttr(),
		Setts:      NewSettings(),
	}
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case common.Message:
		if v.Type == Name {
			p.notify.SaveAndSend(v, p)
		}
	}
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (p *plugin) Save(msg common.Message) (addresses []string, message *m.Message) {
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

	addresses = []string{fmt.Sprintf("%d", attr[AttrEntityID].Int64())}
	return
}

// Send ...
func (p *plugin) Send(address string, message *m.Message) (err error) {

	attr := NewMessageParams()
	if _, err = attr.Deserialize(message.Attributes); err != nil {
		log.Error(err.Error())
		return
	}

	chatID, _ := strconv.ParseInt(address, 0, 64)

	go func() {
		if _, _, err = p.client.SendMsg(attr[AttrBody].String(), chatID); err != nil {
			log.Error(err.Error())
		}
	}()
	return
}
