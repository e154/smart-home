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
	"fmt"
	"sync"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
	"github.com/messagebird/go-rest-api/sms"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	AccessToken string
	Name        string
	balanceLock *sync.Mutex
}

// NewActor ...
func NewActor(settings m.Attributes,
	service supervisor.Service) *Actor {

	accessToken := settings[AttrAccessKey].Decrypt()

	entity := &m.Entity{
		Id:         common.EntityId(fmt.Sprintf("%s.%s", Name, Name)),
		PluginName: Name,
		Attributes: NewAttr(),
	}

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		AccessToken: accessToken,
		Name:        settings[AttrName].String(),
		balanceLock: &sync.Mutex{},
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {

}

// Send ...
func (e *Actor) Send(phone string, message *m.Message) (err error) {

	params := &sms.Params{
		Type:       "sms",
		DataCoding: "unicode",
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	var msg *sms.Message
	if common.TestMode() {
		msg = &sms.Message{
			ID: "123",
		}
	} else {
		var client *messagebird.Client
		if client, err = e.client(); err != nil {
			return
		}

		if msg, err = sms.Create(client, e.Name, []string{phone}, attr[AttrBody].String(), params); err != nil {
			mbErr, ok := err.(messagebird.ErrorResponse)
			if !ok {
				err = errors.Wrap(err, "can`t static cast to messagebird.ErrorResponse")
				return
			}

			err = errors.Wrap(err, mbErr.Errors[0].Description)
			return
		}
	}

	defer func() {
		go func() { _, _ = e.UpdateBalance() }()
	}()

	log.Infof("SMS id(%s) successfully sent to phone '%s'", msg.ID, phone)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var status string

	i := 0
	for range ticker.C {
		if i > 15 {
			return
		}
		if status, err = e.GetStatus(msg.ID); err != nil {
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

// GetStatus ...
func (e *Actor) GetStatus(smsId string) (string, error) {

	if common.TestMode() {
		return StatusDelivered, nil
	}

	client, err := e.client()
	if err != nil {
		return "", err
	}

	msg, err := sms.Read(client, smsId)
	if err != nil {
		return "", errors.Wrap(err, "failed read sms")
	}

	return msg.Recipients.Items[0].Status, nil
}

// UpdateBalance ...
func (e *Actor) UpdateBalance() (bal Balance, err error) {

	e.balanceLock.Lock()
	defer e.balanceLock.Unlock()

	oldState := e.GetEventState()
	now := e.Now(oldState)

	var b *balance.Balance
	if common.TestMode() {
		b = &balance.Balance{
			Payment: "prepaid",
			Type:    "euros",
			Amount:  68.93,
		}
	} else {
		var client *messagebird.Client
		if client, err = e.client(); err != nil {
			return
		}
		if b, err = balance.Read(client); err != nil {
			return
		}
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrPayment] = b.Payment
	attributeValues[AttrType] = b.Type
	attributeValues[AttrAmount] = b.Amount

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

func (e *Actor) client() (client *messagebird.Client, err error) {
	if e.AccessToken == "" {
		err = apperr.ErrBadActorSettingsParameters
		return
	}
	client = messagebird.New(e.AccessToken)
	return
}
