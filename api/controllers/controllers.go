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

package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
)

// Controllers ...
type Controllers struct {
	Auth   ControllerAuth
	Stream ControllerStream
	User   ControllerUser
}

// NewControllers ...
func NewControllers(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	command *endpoint.Endpoint) *Controllers {
	common := NewControllerCommon(adaptors, accessList, command)
	return &Controllers{
		Auth:   NewControllerAuth(common),
		Stream: NewControllerStream(common),
		User:   NewControllerUser(common),
	}
}