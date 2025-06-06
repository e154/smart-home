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
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

const (
	// Name ...
	Name = "email"
	// AttrAuth ...
	AttrAuth = "auth"
	// AttrPass ...
	AttrPass = "pass"
	// AttrSmtp ...
	AttrSmtp = "smtp"
	// AttrPort ...
	AttrPort = "port"
	// AttrSender ...
	AttrSender = "sender"

	// AttrAddresses ...
	AttrAddresses = "addresses"
	// AttrSubject ...
	AttrSubject = "subject"
	// AttrBody ...
	AttrBody = "body"
)

// NewMessageParams ...
func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrAddresses: {
			Name: AttrAddresses,
			Type: common.AttributeString,
		},
		AttrSubject: {
			Name: AttrSubject,
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
		AttrAuth: {
			Name: AttrAuth,
			Type: common.AttributeString,
		},
		AttrPass: {
			Name: AttrPass,
			Type: common.AttributeEncrypted,
		},
		AttrSmtp: {
			Name: AttrSmtp,
			Type: common.AttributeString,
		},
		AttrPort: {
			Name: AttrPort,
			Type: common.AttributeInt,
		},
		AttrSender: {
			Name: AttrSender,
			Type: common.AttributeString,
		},
	}
}
