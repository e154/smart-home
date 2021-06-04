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

package notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	Name        = "notify"
	TopicNotify = "notify"
)

// Stat ...
type Stat struct {
	Workers int `json:"workers"`
}

type EventNewNotify struct {
	From       common.EntityId  `json:"from"`
	Type       string           `json:"type"`
	Attributes m.AttributeValue `json:"attributes"`
}

type ProviderRegistrar interface {
	AddProvider(name string, provider Provider)
	RemoveProvider(name string)
	Provider(name string) (provider Provider, err error)
}

type Provider interface {
	Save(EventNewNotify) (addresses []string, message m.Message)
	Send(addresses string, message m.Message) error
	Attrs() m.Attributes
}
