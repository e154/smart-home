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

package rbac

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/jwt_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"regexp"
)

var (
	log = common.MustGetLogger("rbac")
)

// AccessFilter ...
type AccessFilter struct {
	adaptors            *adaptors.Adaptors
	jwtManager          jwt_manager.JwtManager
	accessListService   access_list.AccessListService
	internalServerError error
}

// NewAccessFilter ...
func NewAccessFilter(adaptors *adaptors.Adaptors,
	jwtManager jwt_manager.JwtManager,
	accessListService access_list.AccessListService) *AccessFilter {
	return &AccessFilter{
		adaptors:            adaptors,
		jwtManager:          jwtManager,
		accessListService:   accessListService,
		internalServerError: status.Error(codes.Unauthenticated, "UNAUTHORIZED"),
	}
}

func (f *AccessFilter) accessDecision(params, method string, accessList access_list.AccessList) bool {

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

// AuthInterceptor ...
func (f *AccessFilter) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	meta, ok := metadata.FromIncomingContext(ctx)

	switch info.FullMethod {
	case "/api.AuthService/Signin":
		return handler(ctx, req)
	}

	if !ok {
		return nil, f.internalServerError
	}
	if len(meta["x-api-key"]) != 1 {
		return nil, f.internalServerError
	}

	// get access token from meta
	var accessToken = meta["x-api-key"][0]

	claims, err := f.jwtManager.Verify(accessToken)
	if err != nil {
		log.Error(err.Error())
		return nil, f.internalServerError
	}

	// если id == 1 is admin
	if claims.UserId == 1 || claims.RoleName == "admin" {
		return f.getUser(claims.UserId, handler, ctx, req)
	}

	// check access filter
	var accessList access_list.AccessList
	if accessList, err = f.accessListService.GetFullAccessList(claims.RoleName); err != nil {
		return nil, f.internalServerError
	}

	const method = "post"
	var requestURI = info.FullMethod

	if ret := f.accessDecision(requestURI, method, accessList); ret {
		return f.getUser(claims.UserId, handler, ctx, req)
	}

	log.Warnf(fmt.Sprintf("access denied: role(%s) [%s] url(%s)", claims.RoleName, method, requestURI))

	return nil, f.internalServerError
}

func (f *AccessFilter) getUser(userId int64, handler grpc.UnaryHandler, ctx context.Context, req interface{}) (interface{}, error) {
	user, err := f.adaptors.User.GetById(userId)
	if err != nil {
		return nil, f.internalServerError
	}
	return handler(context.WithValue(ctx, "currentUser", user), req)
}
