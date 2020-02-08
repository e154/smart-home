// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func NewSlackMessage(channel, text string) *SlackMessage {
	return &SlackMessage{Channel: channel, Text: text}
}

func (s *SlackMessage) SetRender(render *m.TemplateRender) {
	s.Text = render.Body
}

func (s *SlackMessage) Save() (addresses []string, message *m.Message) {

	addresses = []string{s.Channel}
	message = &m.Message{
		Type:      m.MessageTypeSlack,
		SlackText: common.String(s.Text),
	}
	return
}
