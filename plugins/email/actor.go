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
	"fmt"
	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
	"gopkg.in/gomail.v2"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	adaptors *adaptors.Adaptors
	Auth     string
	Pass     string
	Smtp     string
	Port     int64
	Sender   string
}

// NewActor ...
func NewActor(settings m.Attributes,
	service supervisor.Service) *Actor {

	entity := &m.Entity{
		Id:         common.EntityId(fmt.Sprintf("%s.%s", Name, Name)),
		PluginName: Name,
	}

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		Auth:      settings[AttrAuth].String(),
		Pass:      settings[AttrPass].String(),
		Smtp:      settings[AttrSmtp].String(),
		Port:      settings[AttrPort].Int64(),
		Sender:    settings[AttrSender].String(),
	}

	return actor
}

func (e *Actor) Destroy() {

}

// Spawn ...
func (e *Actor) Spawn() {

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
		go func() { _ = e.UpdateStatus() }()
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

	oldState := e.GetEventState()
	now := e.Now(oldState)

	var attributeValues = make(m.AttributeValue)
	// ...

	e.AttrMu.Lock()
	var changed bool
	if changed, err = e.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return
			}
		}
	}
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		StorageSave: true,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
	})

	return
}
