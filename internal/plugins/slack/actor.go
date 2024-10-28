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

package slack

import (
	"context"
	"strings"

	notify2 "github.com/e154/smart-home/internal/plugins/notify"
	"github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/nlopes/slack"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	notify   *notify2.Notify
	Token    string
	UserName string
	api      *slack.Client
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

	token := entity.Settings[AttrToken].Decrypt()

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		notify:    notify2.NewNotify(service.Adaptors()),
		UserName:  entity.Settings[AttrUserName].String(),
		Token:     token,
		api:       slack.New(token),
	}

	return actor
}

func (e *Actor) Destroy() {
	e.Service.EventBus().Unsubscribe(notify2.TopicNotify, e.eventHandler)
	e.notify.Shutdown()
}

func (e *Actor) Spawn() {
	e.notify.Start()
	e.Service.EventBus().Subscribe(notify2.TopicNotify, e.eventHandler, false)
}

// Save ...
func (e *Actor) Save(msg common.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = e.Service.Adaptors().Message.Add(context.Background(), message); err != nil {
		log.Error(err.Error())
	}

	attr := NewAttr()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrChannel].String(), ",")
	return
}

// Send ...
func (e *Actor) Send(address string, message *m.Message) (err error) {

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
func (e *Actor) MessageParams() m.Attributes {
	return NewAttr()
}

func (e *Actor) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case common.Message:
		if v.EntityId != nil && *v.EntityId == e.Id {
			e.notify.SaveAndSend(v, e)
		}
	}
}
