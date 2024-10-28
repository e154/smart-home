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

package jwt_manager

import (
	"time"

	m "github.com/e154/smart-home/pkg/models"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidAccessToken ...
	ErrInvalidAccessToken = errors.New("invalid access token")
	// ErrUnexpectedSigningMethod ...
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	// ErrInvalidTokenClaims ...
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)

// JwtManager ...
type JwtManager interface {
	Generate(*m.User, bool, ...*time.Time) (accessToken string, err error)
	Verify(string) (claims *UserClaims, err error)
	SetHmacKey(hmacKey []byte)
}

// UserClaims ...
type UserClaims struct {
	UserId   int64  `json:"i,omitempty"`
	Username string `json:"n,omitempty"`
	RoleName string `json:"r,omitempty"`
	Root     bool   `json:"root,omitempty"`
}
