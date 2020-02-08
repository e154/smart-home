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
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"strings"
)

type SMS struct {
	phones []string
	Text   string `json:"text"`
}

func NewSMS() (sms *SMS) {
	return &SMS{}
}

func (s *SMS) SetRender(render *m.TemplateRender) {
	s.Text = render.Body
}

func (s *SMS) AddPhone(phone string) {
	if !strings.Contains(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}
	s.phones = append(s.phones, phone)
}

func (s *SMS) Save() (addresses []string, message *m.Message) {

	addresses = s.phones
	message = &m.Message{
		Type:    m.MessageTypeSMS,
		SmsText: common.String(s.Text),
	}
	return
}
