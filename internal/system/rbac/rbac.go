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

package rbac

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/e154/smart-home/internal/system/jwt_manager"
	"github.com/e154/smart-home/internal/system/rbac/access_list"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

var (
	_ plugins.Authorization = (*Authenticator)(nil)
)

var (
	log = logger.MustGetLogger("rbac")
)

type Authenticator struct {
	adaptors          *adaptors.Adaptors
	jwtManager        jwt_manager.JwtManager
	accessListService access_list.AccessListService
	config            *m.AppConfig
}

func NewAuthenticator(adaptors *adaptors.Adaptors,
	jwtManager jwt_manager.JwtManager,
	accessListService access_list.AccessListService,
	config *m.AppConfig) plugins.Authorization {
	return &Authenticator{
		adaptors:          adaptors,
		jwtManager:        jwtManager,
		accessListService: accessListService,
		config:            config,
	}
}

func (a *Authenticator) AuthREST(ctx context.Context, accessToken string, requestURI *url.URL, method string) (*m.User, bool, error) {

	if accessToken == "" {
		return nil, false, apperr.ErrUnauthorized
	}

	claims, err := a.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, false, apperr.ErrUnauthorized
	}

	// if id == 1 is admin
	if claims.UserId == 1 || claims.RoleName == "admin" {
		user, err := a.adaptors.User.GetById(ctx, claims.UserId)
		if err != nil {
			return nil, false, apperr.ErrUnauthorized
		}
		return user, claims.Root, nil
	}

	var accessList access_list.AccessList
	if accessList, err = a.accessListService.GetFullAccessList(ctx, claims.RoleName); err != nil {
		return nil, false, apperr.ErrUnauthorized
	}

	// check access filter
	if ret := a.accessDecision(requestURI.Path, method, accessList); ret {
		user, err := a.adaptors.User.GetById(ctx, claims.UserId)
		if err != nil {
			return nil, false, apperr.ErrUnauthorized
		}
		return user, claims.Root, nil
	}

	log.Warnf(fmt.Sprintf("access denied: role(%s) [%s] url(%s)", claims.RoleName, method, requestURI.Path))

	return nil, false, apperr.ErrUnauthorized
}

func (a *Authenticator) AuthPlain(login, pass string) (*m.User, error) {

	user, err := a.adaptors.User.GetByNickname(context.Background(), login)
	if err != nil {
		if user, err = a.adaptors.User.GetByEmail(context.Background(), login); err != nil {
			err = fmt.Errorf("email %s: %w", login, apperr.ErrUnauthorized)
			log.Warnf("failed login attempt: \"%v\", pass: \"%v\"", login, pass)
			return nil, err
		}
	}

	if !user.CheckPass(pass) {
		return nil, apperr.ErrPassNotValid
	} else if user.Status == "blocked" {
		return nil, apperr.ErrAccountIsBlocked
	}

	return user, nil
}

func (a *Authenticator) accessDecision(params, method string, accessList access_list.AccessList) bool {

	for _, levels := range accessList {
		for _, item := range levels {
			for _, action := range item.Actions {
				if item.Method != method {
					continue
				}

				if ok, _ := regexp.MatchString(action, params); ok {
					return true
				}
			}
		}
	}

	return false
}
