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

package controllers

import (
	"context"

	"github.com/e154/smart-home/api/stub/api"
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerUser ...
type ControllerUser struct {
	*ControllerCommon
}

// NewControllerUser ...
func NewControllerUser(common *ControllerCommon) ControllerUser {
	return ControllerUser{
		ControllerCommon: common,
	}
}

// AddUser ...
func (c ControllerUser) AddUser(ctx context.Context, req *api.NewtUserRequest) (userFull *api.UserFull, err error) {

	user := c.dto.User.AddUserRequest(req)

	if req.Password == req.PasswordRepeat {
		_ = user.SetPass(req.Password)
	}

	var currentUser *m.User
	if currentUser, err = c.currentUser(ctx); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var errs validator.ValidationErrorsTranslations
	var userF *m.User
	userF, errs, err = c.endpoint.User.Add(ctx, user, currentUser)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.User.ToUserFull(userF), nil
}

// GetUserById ...
func (c ControllerUser) GetUserById(ctx context.Context, req *api.GetUserByIdRequest) (*api.UserFull, error) {

	user, err := c.endpoint.User.GetById(ctx, int64(req.Id))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.User.ToUserFull(user), nil
}

// UpdateUserById ...
func (c ControllerUser) UpdateUserById(ctx context.Context, req *api.UpdateUserRequest) (*api.UserFull, error) {

	user := c.dto.User.UpdateUserByIdRequest(req)

	if req.Password != req.PasswordRepeat {
		st := status.New(codes.InvalidArgument, "One or more fields are invalid")
		st, _ = st.WithDetails(&errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: "field not valid",
		})
		return nil, st.Err()
	}

	if req.Password != "" {
		_ = user.SetPass(req.Password)
	}

	user, errs, err := c.endpoint.User.Update(ctx, user)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.User.ToUserFull(user), nil
}

// GetUserList ...
func (c ControllerUser) GetUserList(ctx context.Context, req *api.PaginationRequest) (*api.GetUserListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.User.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.User.ToListResult(items, uint64(total), pagination), nil
}

// DeleteUserById ...
func (c ControllerUser) DeleteUserById(ctx context.Context, req *api.DeleteUserRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.User.Delete(ctx, int64(req.Id)); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}
