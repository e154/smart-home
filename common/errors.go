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

package common

import (
	"errors"
	"fmt"
)

var (

	// ErrInternal ...
	ErrInternal = errors.New("internal error")
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
	// ErrTransactionError ...
	ErrTransactionError = errors.New("transaction error")
	// ErrUnimplemented ...
	ErrUnimplemented = errors.New("unimplemented")

	// ErrPassNotValid ...
	ErrPassNotValid = errors.New("password not valid")
	// ErrAccountIsBlocked ...
	ErrAccountIsBlocked = errors.New("account is blocked")
	// ErrNotAuthorized ...
	ErrNotAuthorized = errors.New("not authorized")
	// ErrTokenIsDeprecated ...
	ErrTokenIsDeprecated = errors.New("token is deprecated")
	// ErrBadLoginOrPassword ...
	ErrBadLoginOrPassword = fmt.Errorf("bad login or password")

	// ErrMqttServerNoWorked ...
	ErrMqttServerNoWorked = errors.New("mqtt server not worked")

	// ErrBadRequestParams ...
	ErrBadRequestParams = errors.New("bad request params")

	// ErrBadActorSettingsParameters ...
	ErrBadActorSettingsParameters = errors.New("bad actor settings parameters")

	// ErrTimeout ...
	ErrTimeout = errors.New("timeout")

	// ErrProviderIsEmpty ...
	ErrProviderIsEmpty = errors.New("provider is empty")

	// ErrUnknownPrototype ...
	ErrUnknownPrototype = errors.New("unknown prototype")
)
