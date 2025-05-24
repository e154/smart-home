// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package apperr

import "errors"

var (
	ErrInternal        = errors.New("internal error")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotAllowed      = errors.New("not allowed")
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUnknownField    = errors.New("unknown field")
	ErrBadJSONRequest  = errors.New("bad JSON request")
	ErrAccessDenied    = errors.New("access denied")
	ErrAccessForbidden = errors.New("access forbidden")
)

var (
	ErrUnimplemented              = errors.New("unimplemented")
	ErrPassNotValid               = errors.New("password not valid")
	ErrAccountIsBlocked           = errors.New("account is blocked")
	ErrTokenIsDeprecated          = errors.New("token is deprecated")
	ErrBadLoginOrPassword         = errors.New("bad login or password")
	ErrMqttServerNoWorked         = errors.New("mqtt server not worked")
	ErrBadRequestParams           = errors.New("bad request params")
	ErrBadActorSettingsParameters = errors.New("bad actor settings parameters")
	ErrTimeout                    = errors.New("timeout")
	ErrProviderIsEmpty            = errors.New("provider is empty")
	ErrBadSettings                = errors.New("there are no settings or the parameters are incorrectly set")
)
