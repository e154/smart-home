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

package twilio

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/e154/smart-home/internal/common"
	notify2 "github.com/e154/smart-home/internal/plugins/notify"
	notifyCommon "github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/pkg/errors"
	"github.com/sfreiberg/gotwilio"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	from      string
	sid       string
	authToken string
	notify    *notify2.Notify
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

	sid := entity.Settings[AttrSid].String()
	authToken := entity.Settings[AttrAuthToken].String()

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		sid:       sid,
		from:      entity.Settings[AttrFrom].String(),
		authToken: authToken,
		notify:    notify2.NewNotify(service.Adaptors()),
	}

	return actor
}

func (e *Actor) Destroy() {
	_ = e.Service.EventBus().Unsubscribe(notify2.TopicNotify, e.eventHandler)
	e.notify.Shutdown()
}

func (e *Actor) Spawn() {
	_ = e.Service.EventBus().Subscribe(notify2.TopicNotify, e.eventHandler, false)
	e.notify.Start()
}

// Send ...
func (e *Actor) Send(phone string, message *m.Message) (err error) {

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
			err = errors.Wrap(apperr.ErrTimeout, "wait ticker")
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
		return "", errors.Wrap(apperr.ErrInternal, ex.Message)
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
func (e *Actor) UpdateBalance() (err error) {

	var balance Balance
	if common.TestMode() {
		balance = Balance{
			Currency:   "euro",
			Balance:    "68.93",
			AccountSid: "XXX",
		}
	} else {
		if balance, err = e.Balance(); err != nil {
			return
		}
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrAmount] = balance.Balance
	attributeValues[AttrSid] = balance.AccountSid
	attributeValues[AttrCurrency] = balance.Currency

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

func (e *Actor) client() (client *gotwilio.Twilio, err error) {
	if e.authToken == "" || e.sid == "" {
		err = apperr.ErrBadActorSettingsParameters
		return
	}
	client = gotwilio.NewTwilioClient(e.sid, e.authToken)
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
