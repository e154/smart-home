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

package html5_notify

import (
	"strconv"
	"strings"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"

	"github.com/e154/smart-home/common/logger"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = logger.MustGetLogger("plugins.html5_notify")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	go func() {
		if err = p.asyncLoad(); err != nil {
			log.Error(err.Error())
		}
	}()

	return nil
}

func (p *plugin) asyncLoad() (err error) {

	// register html5_notify provider
	notify.ProviderManager.AddProvider(Name, p)

	return
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	notify.ProviderManager.RemoveProvider(Name)

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginInstallable
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
		Setts: NewSettings(),
	}
}

// Save ...
func (p *plugin) Save(msg notify.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = p.Adaptors.Message.Add(message); err != nil {
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

	p.EventBus.Publish(bus.TopicEntities, events.EventDirectMessage{
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