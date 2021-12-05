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

package endpoint

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/pkg/errors"
)

const (
	// AdminId ...
	AdminId = 1
)

// AuthEndpoint ...
type AuthEndpoint struct {
	*CommonEndpoint
}

// NewAuthEndpoint ...
func NewAuthEndpoint(common *CommonEndpoint) *AuthEndpoint {
	return &AuthEndpoint{
		CommonEndpoint: common,
	}
}

// SignIn ...
func (a *AuthEndpoint) SignIn(email, password string, ip string) (user *m.User, accessToken string, err error) {

	if user, err = a.adaptors.User.GetByEmail(email); err != nil {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("email %s", email))
		return
	} else if !user.CheckPass(password) {
		err = common.ErrPassNotValid
		return
	} else if user.Status == "blocked" && user.Id != AdminId {
		err = common.ErrAccountIsBlocked
		return
	}

	if accessToken, err = a.jwtManager.Generate(user); err != nil {
		return
	}

	err = a.adaptors.User.SignIn(user, ip)

	log.Infof("Successful login, user: %s", user.Email)

	return
}

// SignOut ...
func (a *AuthEndpoint) SignOut(user *m.User) (err error) {
	err = a.adaptors.User.ClearToken(user)
	return
}

// Recovery ...
func (a *AuthEndpoint) Recovery() {}

// Reset ...
func (a *AuthEndpoint) Reset() {}

// AccessList ...
func (a *AuthEndpoint) AccessList(user *m.User, accessListService access_list.AccessListService) (accessList *access_list.AccessList, err error) {
	accessList = accessListService.List()
	return
}
