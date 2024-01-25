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

package email

import (
	"context"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	notifyCommon "github.com/e154/smart-home/plugins/notify/common"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	adaptors *adaptors.Adaptors
	notify   *notify.Notify
	Auth     string
	Pass     string
	Smtp     string
	Port     int64
	Sender   string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		Auth:      entity.Settings[AttrAuth].String(),
		Pass:      entity.Settings[AttrPass].Decrypt(),
		Smtp:      entity.Settings[AttrSmtp].String(),
		Port:      entity.Settings[AttrPort].Int64(),
		Sender:    entity.Settings[AttrSender].String(),
		notify:    notify.NewNotify(service.Adaptors()),
	}

	return actor
}

func (e *Actor) Destroy() {
	e.Service.EventBus().Unsubscribe(notify.TopicNotify, e.eventHandler)
	e.notify.Shutdown()
}

// Spawn ...
func (e *Actor) Spawn() {
	e.Service.EventBus().Subscribe(notify.TopicNotify, e.eventHandler, false)
	e.notify.Start()
}

// Send ...
func (e *Actor) Send(address string, message *m.Message) error {

	if e.Auth == "" || e.Pass == "" || e.Smtp == "" || e.Port == 0 || e.Sender == "" {
		return apperr.ErrBadActorSettingsParameters
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)
	subject := attr[AttrSubject].String()

	defer func() {
		//go func() { _ = e.UpdateStatus() }()
		log.Infof("Sent email '%s' to: '%s'", subject, address)
	}()

	if common.TestMode() {
		return nil
	}

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
		return err
	}

	return nil
}

// UpdateStatus ...
func (e *Actor) UpdateStatus() (err error) {

	var attributeValues = make(m.AttributeValue)
	// ...

	e.AttrMu.Lock()
	var changed bool
	if changed, err = e.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}
	}
	e.AttrMu.Unlock()

	e.SaveState(false, true)

	return
}

// Save ...
func (e *Actor) Save(msg notifyCommon.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = e.Service.Adaptors().Message.Add(context.Background(), message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrAddresses].String(), ",")
	return
}

// MessageParams ...
func (e *Actor) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (e *Actor) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case notifyCommon.Message:
		if v.EntityId != nil && *v.EntityId == e.Id {
			e.notify.SaveAndSend(v, e)
		}
	}
}
