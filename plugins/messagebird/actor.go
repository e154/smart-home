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

package messagebird

import (
	"context"
	"strings"
	"sync"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
	"github.com/messagebird/go-rest-api/sms"
	"github.com/pkg/errors"

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
	AccessToken string
	Name        string
	notify      *notify.Notify
	balanceLock *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	var token, name string
	if val, ok := entity.Settings[AttrAccessKey]; ok {
		token = val.Decrypt()
	}
	if val, ok := entity.Settings[AttrName]; ok {
		name = val.Decrypt()
	}

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		AccessToken: token,
		Name:        name,
		notify:      notify.NewNotify(service.Adaptors()),
		balanceLock: &sync.Mutex{},
	}

	return actor
}

func (e *Actor) Destroy() {
	e.Service.EventBus().Unsubscribe(notify.TopicNotify, e.eventHandler)
	e.notify.Shutdown()
}

func (e *Actor) Spawn() {
	e.Service.EventBus().Subscribe(notify.TopicNotify, e.eventHandler, false)
	e.notify.Start()
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
	}
	e.AttrMu.Unlock()

	e.SaveState(false, true)

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

	addresses = strings.Split(attr[AttrPhone].String(), ",")
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
