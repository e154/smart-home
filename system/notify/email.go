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
	"strings"
)

// Email ...
type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"template"`
}

// NewEmail ...
func NewEmail() (email *Email) {
	return &Email{}
}

// SetRender ...
func (e *Email) SetRender(render *m.TemplateRender) {
	e.Subject = render.Subject
	e.Body = render.Body
}

// Save ...
func (e *Email) Save() (addresses []string, message *m.Message) {
	e.To = strings.Replace(e.To, " ", "", -1)
	addresses = strings.Split(e.To, ",")
	message = &m.Message{
		Type:         m.MessageTypeEmail,
		EmailFrom:    common.String(e.From),
		EmailSubject: common.String(e.Subject),
		EmailBody:    common.String(e.Body),
	}
	return
}
