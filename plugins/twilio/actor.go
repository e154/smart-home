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

package twilio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/e154/smart-home/system/event_bus/events"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/sfreiberg/gotwilio"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus  event_bus.EventBus
	adaptors  *adaptors.Adaptors
	from      string
	sid       string
	authToken string
}

// NewActor ...
func NewActor(settings m.Attributes,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) *Actor {

	sid := settings[AttrSid].String()
	authToken := settings[AttrAuthToken].String()

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", Name, Name)),
			Name:       Name,
			EntityType: Name,
			AttrMu:     &sync.RWMutex{},
			Attrs:      NewAttr(),
			Manager:    entityManager,
		},
		eventBus:  eventBus,
		adaptors:  adaptors,
		sid:       sid,
		from:      settings[AttrFrom].String(),
		authToken: authToken,
	}

	return actor
}

// Spawn ...
func (p *Actor) Spawn() entity_manager.PluginActor {
	return p
}

// Send ...
func (e *Actor) Send(phone string, message m.Message) (err error) {

	defer func() {
		go func() { _ = e.UpdateBalance() }()
	}()

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception

	if !strings.Contains(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}

	var client *gotwilio.Twilio
	if client, err = e.client(); err != nil {
		return
	}

	resp, ex, err = client.SendSMS(e.from, phone, attr[AttrBody].String(), "", "")
	if err != nil {
		return
	}

	if ex != nil {
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var status string

	i := 0
	for range ticker.C {
		if i > 15 {
			err = errors.Wrap(common.ErrTimeout, "wait ticker")
			return
		}
		if status, err = e.GetStatus(resp.Sid); err != nil {
			return
		}
		if status == StatusDelivered {
			err = nil
			return
		}
		i++
	}

	log.Infof("SMS id(%s) successfully sent to phone '%s'", resp.Sid, phone)

	return
}

// GetStatus ...
func (e *Actor) GetStatus(smsId string) (string, error) {

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception
	var err error

	client, err := e.client()
	if err != nil {
		return "", err
	}
	resp, ex, err = client.GetSMS(smsId)
	if err != nil {
		return "", err
	}

	if ex != nil {
		return "", errors.Wrap(common.ErrInternal, ex.Message)
	}

	return resp.Status, nil
}

// Balance ...
func (e *Actor) Balance() (balance Balance, err error) {

	var uri *url.URL
	if uri, err = url.Parse(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Balance.json", e.sid)); err != nil {
		return
	}

	client := &http.Client{}

	var req *http.Request
	if req, err = http.NewRequest("GET", uri.String(), nil); err != nil {
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", e.sid, e.authToken)))
	req.Header.Add("Authorization", "Basic "+auth)

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&balance)

	return
}

// UpdateBalance ...
func (p *Actor) UpdateBalance() (err error) {

	oldState := p.GetEventState(p)
	now := p.Now(oldState)

	var balance Balance
	if common.TestMode() {
		balance = Balance{
			Currency:   "euro",
			Balance:    "68.93",
			AccountSid: "XXX",
		}
	} else {
		if balance, err = p.Balance(); err != nil {
			return
		}
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrAmount] = balance.Balance
	attributeValues[AttrSid] = balance.AccountSid
	attributeValues[AttrCurrency] = balance.Currency

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

	p.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		StorageSave: true,
		PluginName:  p.Id.PluginName(),
		EntityId:    p.Id,
		OldState:    oldState,
		NewState:    p.GetEventState(p),
	})

	return
}

func (e *Actor) client() (client *gotwilio.Twilio, err error) {
	if e.authToken == "" || e.sid == "" {
		err = common.ErrBadActorSettingsParameters
		return
	}
	client = gotwilio.NewTwilioClient(e.sid, e.authToken)
	return
}
