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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	Name          = "twilio"
	AttrFrom      = "from"
	AttrSid       = "sid"
	AttrAuthToken = "authToken"
	AttrPhone     = "phone"
	AttrBody      = "body"
	AttrAmount    = "amount"
	AttrCurrency  = "currency"
)

func NewAttr() m.Attributes {
	return map[string]*m.Attribute{
		AttrAmount: {
			Name: AttrAmount,
			Type: common.AttributeFloat,
		},
		AttrSid: {
			Name: AttrSid,
			Type: common.AttributeFloat,
		},
		AttrCurrency: {
			Name: AttrCurrency,
			Type: common.AttributeFloat,
		},
	}
}

func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrPhone: {
			Name: AttrPhone,
			Type: common.AttributeString,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
	}
}

func NewSetts() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrFrom: {
			Name: AttrFrom,
			Type: common.AttributeString,
		},
		AttrSid: {
			Name: AttrSid,
			Type: common.AttributeString,
		},
		AttrAuthToken: {
			Name: AttrAuthToken,
			Type: common.AttributeString,
		},
	}
}

// Balance ...
type Balance struct {
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	AccountSid string `json:"account_sid"`
}

const (
	// StatusAccepted ...
	StatusAccepted = "accepted"
	// StatusQueued ...
	StatusQueued = "queued"
	// StatusSending ...
	StatusSending = "sending"
	// StatusReceiving ...
	StatusReceiving = "receiving"
	// StatusReceived ...
	StatusReceived = "received"
	// StatusDelivered ...
	StatusDelivered = "delivered"
	// StatusUndelivered ...
	StatusUndelivered = "undelivered"
	// StatusSent ...
	StatusSent = "sent"
	// StatusFailed ...
	StatusFailed = "failed"
)
