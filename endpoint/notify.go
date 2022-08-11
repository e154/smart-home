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

package endpoint

import (
	"strings"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/email"
	"github.com/e154/smart-home/plugins/messagebird"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/slack"
	"github.com/e154/smart-home/plugins/telegram"
)

// NotifyEndpoint ...
type NotifyEndpoint struct {
	*CommonEndpoint
}

// NewNotifyEndpoint ...
func NewNotifyEndpoint(common *CommonEndpoint) *NotifyEndpoint {
	return &NotifyEndpoint{
		CommonEndpoint: common,
	}
}

// Repeat ...
func (n *NotifyEndpoint) Repeat(id int64) (err error) {

	var msg m.MessageDelivery
	msg, err = n.adaptors.MessageDelivery.GetById(id)
	if err != nil {
		return
	}

	n.eventBus.Publish(notify.TopicNotify, notify.Message{
		Type:       msg.Message.Type,
		Attributes: msg.Message.Attributes,
	})

	return
}

// Send ...
func (n *NotifyEndpoint) Send(params *m.NewNotifrMessage) (err error) {

	var render *m.TemplateRender
	if params.BodyType == "template" && params.Template != nil && params.Params != nil {
		if render, err = n.adaptors.Template.Render(common.StringValue(params.Template), params.Params); err != nil {
			return
		}
	}

	switch params.Type {
	case "email":
		var body = common.StringValue(params.EmailBody)
		if render != nil {
			body = render.Body
		}
		n.eventBus.Publish(notify.TopicNotify, notify.Message{
			Type: email.Name,
			Attributes: map[string]interface{}{
				"addresses": params.Address,
				"subject":   common.StringValue(params.EmailSubject),
				"body":      body,
			},
		})
	case "sms":
		var body = common.StringValue(params.SmsText)
		if render != nil {
			body = render.Body
		}
		for _, address := range strings.Split(params.Address, ",") {
			phone := strings.Replace(address, " ", "", -1)
			if phone == "" {
				continue
			}
			n.eventBus.Publish(notify.TopicNotify, notify.Message{
				Type: messagebird.Name,
				Attributes: map[string]interface{}{
					messagebird.AttrPhone: phone,
					messagebird.AttrBody:  body,
				},
			})
		}
		//todo add sms service....
	case "slack":
		var body = common.StringValue(params.SlackText)
		if render != nil {
			body = render.Body
		}
		n.eventBus.Publish(notify.TopicNotify, notify.Message{
			Type: slack.Name,
			Attributes: map[string]interface{}{
				slack.AttrChannel: params.Address,
				slack.AttrText:    body,
			},
		})
	case "telegram_notify":
		var body = common.StringValue(params.TelegramText)
		if render != nil {
			body = render.Body
		}
		n.eventBus.Publish(notify.TopicNotify, notify.Message{
			Type: telegram.Name,
			Attributes: map[string]interface{}{
				telegram.AttrName: params.Address,
				telegram.AttrBody: body,
			},
		})
	}

	return
}
