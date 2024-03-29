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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// UserEndpoint ...
type UserEndpoint struct {
	*CommonEndpoint
}

// NewUserEndpoint ...
func NewUserEndpoint(common *CommonEndpoint) *UserEndpoint {
	return &UserEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *UserEndpoint) Add(ctx context.Context, params *m.User,
	currentUser *m.User) (result *m.User, err error) {

	user := &m.User{}
	if err = common.Copy(&user, &params); err != nil {
		return
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if currentUser != nil {
		_ = user.UserId.Scan(currentUser.Id)
	}

	if params.Image != nil && params.Image.Id != 0 {
		_ = user.ImageId.Scan(params.Image.Id)
	}

	if params.Role != nil {
		user.RoleName = params.Role.Name
	}

	if params.Meta != nil && len(params.Meta) > 0 {
		for _, rMeta := range params.Meta {
			meta := &m.UserMeta{}
			if err = common.Copy(&meta, &rMeta); err != nil {
				return
			}
			user.Meta = append(user.Meta, meta)
		}
	}

	// check user status
	switch user.Status {
	case "active", "blocked":
	default:
		user.Status = "blocked"
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = n.adaptors.User.Add(ctx, user); err != nil {
		return
	}

	result, err = n.GetById(ctx, id)

	return
}

// GetById ...
func (n *UserEndpoint) GetById(ctx context.Context, userId int64) (result *m.User, err error) {

	result, err = n.adaptors.User.GetById(ctx, userId)

	return
}

// Delete ...
func (n *UserEndpoint) Delete(ctx context.Context, userId int64) (err error) {

	var user *m.User
	user, err = n.adaptors.User.GetById(ctx, userId)
	if err != nil {
		return
	}

	if user.Role.Name == "admin" && user.Id == 1 {
		err = apperr.ErrUserDeleteForbidden
		return
	}

	err = n.adaptors.User.Delete(ctx, user.Id)

	return
}

// GetList ...
func (n *UserEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.User, total int64, err error) {

	result, total, err = n.adaptors.User.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// Update ...
func (n *UserEndpoint) Update(ctx context.Context, params *m.User) (result *m.User, err error) {

	var user *m.User
	user, err = n.adaptors.User.GetById(ctx, params.Id)
	if err != nil {
		return
	}

	if user.Id == 1 && user.RoleName == "admin" && n.checkSuperUser(ctx) {
		err = apperr.ErrUserUpdateForbidden
		return
	}

	_ = common.Copy(&user, &params, common.JsonEngine)

	if params.Image != nil && params.Image.Id != 0 {
		_ = user.ImageId.Scan(params.Image.Id)
	}

	if params.Role != nil {
		user.RoleName = params.Role.Name
	}

	if params.Meta != nil && len(params.Meta) > 0 {
		for _, rMeta := range params.Meta {
			meta := &m.UserMeta{}
			if err = common.Copy(&meta, &rMeta); err != nil {
				return
			}
			user.Meta = append(user.Meta, meta)
		}
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.User.Update(ctx, user); err != nil {
		return
	}

	result, err = n.GetById(ctx, user.Id)

	return
}

// UpdateStatus ...
func (n *UserEndpoint) UpdateStatus(ctx context.Context, userId int64, newStatus string) (err error) {

	var user *m.User
	user, err = n.adaptors.User.GetById(ctx, userId)
	if err != nil {
		return
	}

	user.Status = newStatus

	// check user status
	switch user.Status {
	case "active", "blocked":
	default:
		user.Status = "blocked"
	}

	err = n.adaptors.User.Update(ctx, user)

	return
}
