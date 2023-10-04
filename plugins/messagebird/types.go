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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "messagebird"
	// AttrAccessKey ...
	AttrAccessKey = "access_key"
	// AttrName ...
	AttrName = "name"
	// AttrPhone ...
	AttrPhone = "phone"
	// AttrBody ...
	AttrBody = "body"
	// AttrPayment ...
	AttrPayment = "Payment"
	// AttrType ...
	AttrType = "Type"
	// AttrAmount ...
	AttrAmount = "Amount"

	Version = "0.0.1"
)

const (
	// StatusDelivered ...
	StatusDelivered = "delivered"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return map[string]*m.Attribute{
		AttrPayment: {
			Name: AttrPayment,
			Type: common.AttributeString,
		},
		AttrType: {
			Name: AttrType,
			Type: common.AttributeString,
		},
		AttrAmount: {
			Name: AttrAmount,
			Type: common.AttributeFloat,
		},
	}
}

// NewMessageParams ...
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

// NewSettings ...
func NewSettings() m.Attributes {
	return map[string]*m.Attribute{
		AttrAccessKey: {
			Name: AttrAccessKey,
			Type: common.AttributeEncrypted,
		},
		AttrName: {
			Name: AttrName,
			Type: common.AttributeString,
		},
	}
}

// Balance ...
type Balance struct {
	Payment string
	Type    string
	Amount  float32
}
