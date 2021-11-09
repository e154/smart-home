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
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"go.uber.org/fx"
	"strings"
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

func (j *jwtManager) Start() (err error) {
	_, err = j.secretKey()
	return
}

func (j *jwtManager) Generate(user *m.User) (accessToken string, err error) {

	now := time.Now()
	data := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.AddDate(0, 1, 0).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "server",
			NotBefore: now.Unix(),
		},
		UserId:   user.Id,
		Username: user.Nickname,
		RoleName: user.RoleName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	if accessToken, err = token.SignedString(j.hmacKey); err != nil {
		return
	}

	return
}

func (j *jwtManager) Verify(accessToken string) (claims *UserClaims, err error) {

	if len(strings.Split(accessToken, ".")) != 3 {
		err = errors.New("access token invalid")
		return
	}

	var token *jwt.Token
	token, err = jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return j.hmacKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	var ok bool
	if claims, ok = token.Claims.(*UserClaims); !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if err = claims.Valid(); err != nil {
		claims = nil
	}

	return
}

func (j *jwtManager) secretKey() (hmacKey []byte, err error) {

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
