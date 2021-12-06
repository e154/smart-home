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

package jwt_manager

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/fx"
	"time"
)

var (
	log = common.MustGetLogger("jwt")
)

type jwtManager struct {
	adaptors      *adaptors.Adaptors
	tokenDuration time.Duration
	hmacKey       []byte
}

// NewJwtManager ...
func NewJwtManager(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors) (mananger JwtManager) {

	mananger = &jwtManager{adaptors: adaptors}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return mananger.Start()
		},
	})

	return
}

// Start ...
func (j *jwtManager) Start() (err error) {
	_, err = j.getSecretKey()
	return
}

// Generate ...
func (j *jwtManager) Generate(user *m.User, opts ...*time.Time) (accessToken string, err error) {

	var exp *int64
	if len(opts) > 0 {
		exp = common.Int64(opts[0].Unix())
	}

	now := time.Now()
	if exp == nil {
		exp = common.Int64(now.AddDate(0, 1, 0).Unix())
	}

	data := jwt.MapClaims{
		"exp": exp,
		"iat": now.Unix(),
		"iss": "server",
		"nbf": now.Unix(),
		"i":   user.Id,
		"n":   user.Nickname,
		"r":   user.RoleName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	accessToken, err = token.SignedString(j.hmacKey)

	return
}

// Verify ...
func (j *jwtManager) Verify(accessToken string) (claims *UserClaims, err error) {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.hmacKey, nil
	})

	if token == nil {
		err = errors.New("invalid access token")
		return
	}

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err = mapClaims.Valid(); err != nil {
			return
		}

		claims = &UserClaims{}
		err = common.Copy(claims, mapClaims, common.JsonEngine)
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}

	return
}

func (j *jwtManager) getSecretKey() (hmacKey []byte, err error) {

	if j.hmacKey != nil && len(j.hmacKey) > 0 {
		return j.hmacKey, nil
	}

	var variable m.Variable
	if variable, err = j.adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = m.Variable{
			Name:  "hmacKey",
			Value: common.ComputeHmac256(),
		}
		if err = j.adaptors.Variable.Add(variable); err != nil {
			log.Error(err.Error())
		}
	}

	if hmacKey, err = hex.DecodeString(variable.Value); err != nil {
		return
	}

	j.hmacKey = hmacKey

	return
}

// SetHmacKey ...
func (j *jwtManager) SetHmacKey(hmacKey []byte) {
	j.hmacKey = hmacKey
}
