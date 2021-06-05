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
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/sfreiberg/gotwilio"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Provider ...
type Provider struct {
	adaptors  *adaptors.Adaptors
	from      string
	sid       string
	authToken string
	client    *gotwilio.Twilio
	sync.RWMutex
	balance Balance
}

// NewProvider ...
func NewProvider(attrs m.Attributes,
	adaptors *adaptors.Adaptors) (p *Provider, err error) {

	sid := attrs[AttrSid].String()
	authToken := attrs[AttrAuthToken].String()
	p = &Provider{
		adaptors:  adaptors,
		sid:       sid,
		from:      attrs[AttrFrom].String(),
		authToken: authToken,
		client:    gotwilio.NewTwilioClient(sid, authToken),
	}

	p.UpdateBalance()

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

	addresses = strings.Split(attr[AttrPhone].String(), ",")
	return
}

// Send ...
func (e *Provider) Send(phone string, message m.Message) (err error) {

	attr := NewMessageParams()
	attr.Deserialize(message.Attributes)

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception

	if !strings.Contains(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}

	resp, ex, err = e.client.SendSMS(e.from, phone, attr[AttrBody].String(), "", "")
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
			err = errors.New("status timeout")
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

// MessageParams ...
// Channel
// Text
func (e *Provider) MessageParams() m.Attributes {
	return NewMessageParams()
}

// GetStatus ...
func (e *Provider) GetStatus(smsId string) (string, error) {

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception
	var err error

	resp, ex, err = e.client.GetSMS(smsId)
	if err != nil {
		return "", err
	}

	if ex != nil {
		return "", errors.New(ex.Message)
	}

	return resp.Status, nil
}

// Balance ...
func (e *Provider) Balance() (balance Balance, err error) {

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
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&balance)

	return
}

// UpdateBalance ...
func (e *Provider) UpdateBalance() {
	balance, err := e.Balance()
	if err != nil {
		return
	}
	e.Lock()
	e.balance = balance
	e.Unlock()
}

// GetBalance ...
func (p *Provider) GetBalance() Balance {
	p.RLock()
	defer p.RUnlock()
	return p.balance
}
