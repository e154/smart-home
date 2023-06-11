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

package slack

import (
	"strings"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/nlopes/slack"
)

// Provider ...
type Provider struct {
	adaptors *adaptors.Adaptors
	Token    string
	UserName string
	api      *slack.Client
}

// NewProvider ...
func NewProvider(attrs m.Attributes,
	adaptors *adaptors.Adaptors) (p *Provider, err error) {

	token := attrs[AttrToken].String()
	p = &Provider{
		adaptors: adaptors,
		Token:    token,
		UserName: attrs[AttrUserName].String(),
		api:      slack.New(token),
	}

	return
}

// Save ...
func (e *Provider) Save(msg notify.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = e.adaptors.Message.Add(message); err != nil {
		log.Error(err.Error())
	}

	attr := NewAttr()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrChannel].String(), ",")
	return
}

// Send ...
func (e *Provider) Send(address string, message *m.Message) (err error) {

	if e.Token == "" || e.UserName == "" {
		return apperr.ErrBadActorSettingsParameters
	}

	attr := NewAttr()
	_, _ = attr.Deserialize(message.Attributes)

	options := []slack.MsgOption{
		slack.MsgOptionText(attr[AttrText].String(), false),
	}

	if e.UserName != "" {
		options = append(options, slack.MsgOptionUsername(e.UserName))
	}

	var channelID, timestamp string
	if channelID, timestamp, err = e.api.PostMessage(address, options...); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("Message successfully sent to channel '%s' at '%s'", channelID, timestamp)

	return
}

// MessageParams ...
// Channel
// Text
func (e *Provider) MessageParams() m.Attributes {
	return NewAttr()
}
