// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package events

import (
	"reflect"

	m "github.com/e154/smart-home/pkg/models"

	"github.com/iancoleman/strcase"
)

type OwnerType string

const (
	OwnerUser   = OwnerType("user")
	OwnerSystem = OwnerType("system")
)

type Common struct {
	Owner     OwnerType `json:"owner"`
	SessionID string    `json:"session_id"`
	User      *m.User   `json:"user"`
}

func (c Common) UserId() int64 {
	if c.User != nil {
		return c.User.Id
	} else {
		return 0
	}
}

func EventName(event interface{}) string {
	if t := reflect.TypeOf(event); t.Kind() == reflect.Ptr {
		return strcase.ToSnake(t.Elem().Name())
	} else {
		return strcase.ToSnake(t.Name())
	}
}
