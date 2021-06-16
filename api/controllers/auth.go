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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ControllerAuth struct {
	*ControllerCommon
}

func NewControllerAuth(common *ControllerCommon) ControllerAuth {
	return ControllerAuth{
		ControllerCommon: common,
	}
}

func (a ControllerAuth) Signin(ctx context.Context, _ *emptypb.Empty) (resp *api.SigninResponse, err error) {

	var internalServerError = status.Error(codes.Unauthenticated, "INTERNAL_SERVER_ERROR")
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, internalServerError
	}
	if len(meta["authorization"]) != 1 {
		return nil, internalServerError
	}

	username, pass, _ := a.parseBasicAuth(meta["authorization"][0])

	var user *m.User
	var accessToken string
	if user, accessToken, err = a.endpoint.Auth.SignIn(username, pass, ""); err != nil {
		return
	}

	currentUser := &api.CurrentUser{}
	_ = common.Copy(&currentUser, &user, common.JsonEngine)

	resp = &api.SigninResponse{
		CurrentUser: currentUser,
		AccessToken: accessToken,
	}

	return
}

func (a ControllerAuth) Signout(context.Context, *api.SignoutRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a ControllerAuth) AccessList(context.Context, *api.AccessListRequest) (*api.AccessListResponse, error) {
	return &api.AccessListResponse{}, nil
}
