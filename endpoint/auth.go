// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"time"
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
		err = errors.New("user not found")
		return
	} else if !user.CheckPass(password) {
		err = errors.New("password not valid")
		return
	} else if user.Status == "blocked" && user.Id != AdminId {
		err = errors.New("account is blocked")
		return
	}

	if err = a.adaptors.User.SignIn(user, ip); err != nil {
		return
	}

	//if _, err = a.adaptors.User.NewToken(user); err != nil {
	//	return
	//}

	// ger hmac key
	var variable m.Variable
	if variable, err = a.adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = m.Variable{
			Name:  "hmacKey",
			Value: common.ComputeHmac256(),
		}
		if err = a.adaptors.Variable.Add(variable); err != nil {
			log.Error(err.Error())
		}
	}

	var hmacKey []byte
	hmacKey, err = hex.DecodeString(variable.Value)
	if err != nil {
		return
	}

	now := time.Now()
	data := map[string]interface{}{
		"userId": user.Id,
		"iss":    "server",
		"nbf":    now.Unix(),
		"iat":    now.Unix(),
		"exp":    now.AddDate(0, 1, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))

	if accessToken, err = token.SignedString(hmacKey); err != nil {
		return
	}

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
