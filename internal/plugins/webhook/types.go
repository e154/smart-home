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

package webhook

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

const (
	// Name ...
	Name = "webhook"
)

const ()

const (
	AttrToken = "token"
	AttrPath  = "path"

	AttrHost    = "host"
	AttrSize    = "size"
	AttrBody    = "body"
	AttrUrl     = "url"
	AttrMethod  = "method"
	AttrHeaders = "headers"
)

const (
	// FuncEntityAction ...
	FuncEntityAction = "webhookAction"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrHost: {
			Name: AttrHost,
			Type: common.AttributeString,
		},
		AttrSize: {
			Name: AttrSize,
			Type: common.AttributeInt,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
		AttrUrl: {
			Name: AttrUrl,
			Type: common.AttributeString,
		},
		AttrMethod: {
			Name: AttrMethod,
			Type: common.AttributeString,
		},
		AttrHeaders: {
			Name: AttrHeaders,
			Type: common.AttributeString,
		},
	}
}

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrToken: {
			Name: AttrToken,
			Type: common.AttributeEncrypted,
		},
		AttrPath: {
			Name: AttrPath,
			Type: common.AttributeString,
		},
	}
}

// NewStates ...
func NewStates() (states map[string]plugins.ActorState) {
	return nil
}
