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

package html5_notify

import (
	"context"
	"strconv"
	"strings"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/notify/common"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.html5_notify")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	notify *notify.Notify
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}
	p.notify = notify.NewNotify(service.Adaptors())
	p.notify.Start()
	_ = p.Service.EventBus().Subscribe(notify.TopicNotify, p.eventHandler, false)
	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	_ = p.Service.EventBus().Unsubscribe(notify.TopicNotify, p.eventHandler)
	err = p.Plugin.Unload(ctx)
	p.notify.Shutdown()
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
		ActorSetts: NewSettings(),
	}
}

// Save ...
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

	addresses = strings.Split(attr[AttrUserIDS].String(), ",")
	return
}

// Send ...
func (p *plugin) Send(address string, message *m.Message) (err error) {

	userID, _ := strconv.ParseInt(address, 0, 64)

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	p.Service.EventBus().Publish("system/plugins/html5_notify", events.EventDirectMessage{
		UserID: userID,
		Query:  "html5_notify",
		Message: Notification{
			Title: attr[AttrTitle].String(),
			Options: &NotificationOptions{
				Badge:              attr[AttrBadge].String(),
				Body:               attr[AttrBody].String(),
				Data:               attr[AttrData].String(),
				Dir:                attr[AttrDir].String(),
				Icon:               attr[AttrIcon].String(),
				Image:              attr[AttrImage].String(),
				Lang:               attr[AttrLang].String(),
				Renotify:           attr[AttrRenotify].Bool(),
				RequireInteraction: attr[AttrRequireInteraction].Bool(),
				Silent:             attr[AttrSilent].Bool(),
				Tag:                attr[AttrTag].String(),
				Timestamp:          attr[AttrTimestamp].Int64(),
			},
		},
	},
	)

	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case common.Message:
		if v.Type == Name {
			p.notify.SaveAndSend(v, p)
		}
	}
}
