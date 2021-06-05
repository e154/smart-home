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

package email

import (
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"gopkg.in/gomail.v2"
	"strings"
)

// Provider ...
type Provider struct {
	adaptors *adaptors.Adaptors
	Auth     string
	Pass     string
	Smtp     string
	Port     int64
	Sender   string
}

// NewProvider ...
func NewProvider(attrs m.Attributes,
	adaptors *adaptors.Adaptors) (p *Provider, err error) {

	p = &Provider{
		adaptors: adaptors,
		Auth:     attrs[AttrAuth].String(),
		Pass:     attrs[AttrPass].String(),
		Smtp:     attrs[AttrSmtp].String(),
		Port:     attrs[AttrPort].Int64(),
		Sender:   attrs[AttrSender].String(),
	}

	return
}

// Save ...
func (e *Provider) Save(msg notify.Message) (addresses []string, message m.Message) {
	message = m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = e.adaptors.Message.Add(message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrAddresses].String(), ",")
	return
}

// Send ...
func (e *Provider) Send(address string, message m.Message) error {

	if e.Auth == "" || e.Pass == "" || e.Smtp == "" || e.Port == 0 || e.Sender == "" {
		return errors.New("bad settings parameters")
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)
	subject := attr[AttrSubject].String()

	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":     {e.Sender},
		"Reply-To": {e.Sender},
		"To":       {address},
		"Subject":  {subject},
	})

	m.SetBody("text/html", attr[AttrBody].String())

	d := gomail.NewPlainDialer(e.Smtp, int(e.Port), e.Auth, e.Pass)
	if err := d.DialAndSend(m); err != nil {
		return errors.New(err.Error())
	}

	log.Debugf("Sent email '%s' to: '%s'", subject, address)

	return nil
}

// MessageParams ...
// Addresses
// Subject
// Body
func (e *Provider) MessageParams() m.Attributes {
	return NewMessageParams()
}
