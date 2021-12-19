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
	"fmt"
	"sync"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"gopkg.in/gomail.v2"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus event_bus.EventBus
	adaptors *adaptors.Adaptors
	Auth     string
	Pass     string
	Smtp     string
	Port     int64
	Sender   string
}

// NewActor ...
func NewActor(settings m.Attributes,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) *Actor {

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", Name, Name)),
			Name:       Name,
			EntityType: Name,
			AttrMu:     &sync.RWMutex{},
			Manager:    entityManager,
		},
		eventBus: eventBus,
		adaptors: adaptors,
		Auth:     settings[AttrAuth].String(),
		Pass:     settings[AttrPass].String(),
		Smtp:     settings[AttrSmtp].String(),
		Port:     settings[AttrPort].Int64(),
		Sender:   settings[AttrSender].String(),
	}

	return actor
}

// Spawn ...
func (p *Actor) Spawn() entity_manager.PluginActor {
	return p
}

// Send ...
func (e *Actor) Send(address string, message m.Message) error {

	if e.Auth == "" || e.Pass == "" || e.Smtp == "" || e.Port == 0 || e.Sender == "" {
		return errors.New("bad settings parameters")
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)
	subject := attr[AttrSubject].String()

	defer func() {
		go e.UpdateStatus()
		log.Debugf("Sent email '%s' to: '%s'", subject, address)
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
		return errors.New(err.Error())
	}

	return nil
}

// UpdateStatus ...
func (p *Actor) UpdateStatus() (err error) {

	oldState := p.GetEventState(p)
	now := p.Now(oldState)

	var attributeValues = make(m.AttributeValue)
	// ...

	p.AttrMu.Lock()
	var changed bool
	if changed, err = p.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				p.AttrMu.Unlock()
				return
			}
		}
	}
	p.AttrMu.Unlock()

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		StorageSave: true,
		PluginName:  p.Id.PluginName(),
		EntityId:    p.Id,
		OldState:    oldState,
		NewState:    p.GetEventState(p),
	})

	return
}
