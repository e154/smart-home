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
	"context"
	"fmt"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/email"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/access_list"
	"github.com/pkg/errors"
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
func (a *AuthEndpoint) SignIn(ctx context.Context, email, password string, ip string) (user *m.User, accessToken string, err error) {

	if user, err = a.adaptors.User.GetByEmail(email); err != nil {
		err = errors.Wrap(apperr.ErrUnauthorized, fmt.Sprintf("email %s", email))
		return
	} else if !user.CheckPass(password) {
		err = apperr.ErrPassNotValid
		return
	} else if user.Status == "blocked" && user.Id != AdminId {
		err = apperr.ErrAccountIsBlocked
		return
	}

	if accessToken, err = a.jwtManager.Generate(user); err != nil {
		err = errors.Wrap(apperr.ErrUnauthorized, err.Error())
		return
	}

	if err = a.adaptors.User.SignIn(user, ip); err != nil {
		err = errors.Wrap(apperr.ErrUnauthorized, err.Error())
		return
	}

	log.Infof("Successful login, user: %s", user.Email)

	return
}

// SignOut ...
func (a *AuthEndpoint) SignOut(ctx context.Context, user *m.User) (err error) {
	err = a.adaptors.User.ClearToken(user)
	if err != nil {
		err = errors.Wrap(apperr.ErrNotAllowed, err.Error())
		return
	}
	return
}

// PasswordReset ...
func (a *AuthEndpoint) PasswordReset(ctx context.Context, userEmail string, token, newPassword *string) (err error) {

	if token != nil {

		if newPassword == nil {
			err = errors.New("password is required")
			return
		}

		var user *m.User
		if user, err = a.adaptors.User.GetByResetPassToken(*token); err != nil {
			return
		}

		if err = user.SetPass(*newPassword); err != nil {
			return
		}

		user.ResetPasswordToken = ""
		user.ResetPasswordSentAt = nil
		if err = a.adaptors.User.Update(user); err == nil {
			log.Warnf("The password for the %s user has just been updated", user.Email)
		}

		return
	}

	var user *m.User
	if user, err = a.adaptors.User.GetByEmail(userEmail); err != nil {
		err = errors.Wrap(apperr.ErrNotAllowed, err.Error())
		return
	}

	if user.ResetPasswordSentAt != nil && time.Now().Before(*user.ResetPasswordSentAt) {
		err = errors.Wrap(apperr.ErrNotAllowed, "reset request already exists")
		return
	}

	var resetToken string
	if resetToken, err = a.adaptors.User.GenResetPassToken(user); err != nil {
		err = errors.Wrap(apperr.ErrNotAllowed, err.Error())
		return
	}

	renderParams := map[string]interface{}{
		"site:name":               "Smart home",
		"user:name:first":         user.FirstName,
		"user:name:last":          user.LastName,
		"user:one-time-login-url": fmt.Sprintf("%s/#/password_reset?t=%s", a.appConfig.ApiFullAddress(), resetToken),
	}

	var render *m.TemplateRender
	if render, err = a.adaptors.Template.Render("password_reset", renderParams); err != nil {
		return
	}

	a.eventBus.Publish(notify.TopicNotify, notify.Message{
		Type: email.Name,
		Attributes: map[string]interface{}{
			email.AttrAddresses: user.Email,
			email.AttrSubject:   "Reset your Smart home password",
			email.AttrBody:      render.Body,
		},
	})

	return
}

// AccessList ...
func (a *AuthEndpoint) AccessList(ctx context.Context, user *m.User, accessListService access_list.AccessListService) (accessList *access_list.AccessList, err error) {
	accessList = accessListService.List()
	return
}
