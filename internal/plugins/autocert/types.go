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

package autocert

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

const (
	// Name ...
	Name                     = "autocert"
	Version                  = "0.0.1"
	ActionRequestCertificate = "RequestCertificate"
	AttrPublicKey            = "publicKey"
	AttrPrivateKey           = "privateKey"
	AttrCloudflareAPIToken   = "CloudflareAPIToken"
	AttrDomains              = "domains"
	AttrEmails               = "emails"
	AttrProduction           = "production"
	StateSuccessfully        = "successfully"
	StateError               = "error"
)

// store entity status in this struct
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPublicKey: {
			Name: AttrPublicKey,
			Type: common.AttributeString,
		},
		AttrPrivateKey: {
			Name: AttrPrivateKey,
			Type: common.AttributeString,
		},
	}
}

// entity settings
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrCloudflareAPIToken: {
			Name: AttrCloudflareAPIToken,
			Type: common.AttributeEncrypted,
		},
		AttrDomains: {
			Name: AttrDomains,
			Type: common.AttributeString,
		},
		AttrEmails: {
			Name: AttrEmails,
			Type: common.AttributeString,
		},
		AttrProduction: {
			Name: AttrProduction,
			Type: common.AttributeBool,
		},
	}
}

// state list entity
func NewStates() (states map[string]plugins.ActorState) {
	states = map[string]plugins.ActorState{
		StateSuccessfully: {
			Name:        StateSuccessfully,
			Description: "successfully",
		},
		StateError: {
			Name:        StateError,
			Description: "error",
		},
	}
	return
}

// entity action list
func NewActions() map[string]plugins.ActorAction {
	return map[string]plugins.ActorAction{
		ActionRequestCertificate: {
			Name:        "RequestCertificate",
			Description: "request certificate",
		},
	}
}
