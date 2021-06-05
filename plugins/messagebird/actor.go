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

package messagebird

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
	"github.com/messagebird/go-rest-api/sms"
	"strings"
	"sync"
	"time"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus    event_bus.EventBus
	adaptors    *adaptors.Adaptors
	AccessToken string
	Name        string
	balanceLock *sync.Mutex
}

// NewActor ...
func NewActor(settings m.Attributes,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) *Actor {

	accessToken := settings[AttrAccessKey].String()

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", Name, Name)),
			Name:       Name,
			EntityType: Name,
			AttrMu:     &sync.RWMutex{},
			Attrs:      NewAttr(),
			Manager:    entityManager,
		},
		eventBus:    eventBus,
		adaptors:    adaptors,
		AccessToken: accessToken,
		Name:        settings[AttrName].String(),
		balanceLock: &sync.Mutex{},
	}

	return actor
}

func (p *Actor) Spawn() entity_manager.PluginActor {
	return p
}

// Save ...
func (p *Actor) Save(msg notify.Message) (addresses []string, message m.Message) {
	message = m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = p.adaptors.Message.Add(message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrPhone].String(), ",")
	return
}

// Send ...
func (p *Actor) Send(phone string, message m.Message) (err error) {

	if p.AccessToken == "" {
		return errors.New("bad settings parameters")
	}

	defer func() {
		go p.UpdateBalance()
	}()

	params := &sms.Params{
		Type:       "sms",
		DataCoding: "unicode",
	}

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)

	client := messagebird.New(p.AccessToken)

	var msg *sms.Message
	if common.TestMode() {
		msg = &sms.Message{
			ID: "123",
		}
	} else {
		if msg, err = sms.Create(client, p.Name, []string{phone}, attr[AttrBody].String(), params); err != nil {
			mbErr, ok := err.(messagebird.ErrorResponse)
			if !ok {
				err = errors.New(err.Error())
				return
			}

			err = errors.New(mbErr.Errors[0].Description)
			return
		}
	}

	log.Infof("SMS id(%s) successfully sent to phone '%s'", msg.ID, phone)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var status string

	i := 0
	for range ticker.C {
		if i > 15 {
			return
		}
		if status, err = p.GetStatus(msg.ID); err != nil {
			return
		}
		if status == StatusDelivered {
			err = nil
			return
		}
		i++
	}

	return
}

// MessageParams ...
func (p *Actor) MessageParams() m.Attributes {
	return NewMessageParams()
}

// GetStatus ...
func (p *Actor) GetStatus(smsId string) (string, error) {

	if common.TestMode() {
		return StatusDelivered, nil
	}

	client := messagebird.New(p.AccessToken)
	msg, err := sms.Read(client, smsId)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return msg.Recipients.Items[0].Status, nil
}

// UpdateBalance ...
func (p *Actor) UpdateBalance() (bal Balance, err error) {

	p.balanceLock.Lock()
	defer p.balanceLock.Lock()

	oldState := p.GetEventState(p)
	now := p.Now(oldState)

	var b *balance.Balance
	if common.TestMode() {
		b = &balance.Balance{
			Payment: "prepaid",
			Type:    "euros",
			Amount:  68.93,
		}
	} else {
		client := messagebird.New(p.AccessToken)
		if b, err = balance.Read(client); err != nil {
			return
		}
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrPayment] = b.Payment
	attributeValues[AttrType] = b.Type
	attributeValues[AttrAmount] = b.Amount

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
		Type:        p.Id.Type(),
		EntityId:    p.Id,
		OldState:    oldState,
		NewState:    p.GetEventState(p),
	})

	return
}
