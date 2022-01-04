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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User ...
type User struct{}

// NewUserDto ...
func NewUserDto() User {
	return User{}
}

// FromAddUser ...
func (u User) FromAddUser(req *api.NewtUserRequest) (user *m.User) {
	user = &m.User{}
	common.Copy(&user, req, common.JsonEngine)
	return
}

// ToUserFull ...
func (u User) ToUserFull(user *m.User) (result *api.UserFull) {
	roleDto := NewRoleDto()
	imageDto := NewImageDto()
	result = &api.UserFull{
		Id:                  user.Id,
		Nickname:            user.Nickname,
		FirstName:           common.String(user.FirstName),
		LastName:            common.String(user.LastName),
		Email:               user.Email,
		Status:              user.Status,
		SignInCount:         user.SignInCount,
		Role:                roleDto.ToGRole(user.Role),
		RoleName:            user.RoleName,
		Lang:                user.Lang,
		AuthenticationToken: common.StringValue(user.AuthenticationToken),
		CreatedAt:           timestamppb.New(user.CreatedAt),
		UpdatedAt:           timestamppb.New(user.UpdatedAt),
	}

	if user.LastSignInIp != "" {
		result.LastSignInIp = common.String(user.LastSignInIp)
	}

	if user.CurrentSignInIp != "" {
		result.CurrentSignInIp = common.String(user.CurrentSignInIp)
	}

	// history
	if user.History != nil {
		result.History = make([]*api.UserHistory, 0, len(user.History))
		for _, h := range user.History {
			result.History = append(result.History, &api.UserHistory{
				Ip:   h.Ip,
				Time: timestamppb.New(h.Time),
			})
		}
	}

	// meta
	if user.Meta != nil {
		result.Meta = make([]*api.UserMeta, 0, len(user.Meta))
		for _, m := range user.Meta {
			result.Meta = append(result.Meta, &api.UserMeta{
				Key:   m.Key,
				Value: m.Value,
			})
		}
	}

	// image
	if user.Image != nil {
		result.Image = imageDto.ToImage(user.Image)
	}

	// todo fix time
	// times ...
	if user.CurrentSignInAt != nil {
		result.CurrentSignInAt = timestamppb.New(common.TimeValue(user.CurrentSignInAt))
	}
	// times ...
	if user.LastSignInAt != nil {
		result.LastSignInAt = timestamppb.New(common.TimeValue(user.LastSignInAt))
	}
	// times ...
	if user.ResetPasswordSentAt != nil {
		result.ResetPasswordSentAt = timestamppb.New(common.TimeValue(user.ResetPasswordSentAt))
	}
	// times ...
	if user.DeletedAt != nil {
		result.DeletedAt = timestamppb.New(common.TimeValue(user.DeletedAt))
	}

	// parent
	if user.User != nil {
		result.User = &api.UserFull_Parent{
			Id:       user.User.Id,
			Nickname: user.User.Nickname,
		}
	}
	return
}

// ToUserShot ...
func (u User) ToUserShot(user *m.User) (result *api.UserShot) {

	roleDto := NewRoleDto()
	imageDto := NewImageDto()
	result = &api.UserShot{
		Id:        user.Id,
		Nickname:  user.Nickname,
		FirstName: common.String(user.FirstName),
		LastName:  common.String(user.LastName),
		Email:     user.Email,
		Status:    user.Status,
		Lang:      user.Lang,
		Role:      roleDto.ToGRole(user.Role),
		RoleName:  user.RoleName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	// image
	if user.Image != nil {
		result.Image = imageDto.ToImage(user.Image)
	}

	// parent
	if user.User != nil {
		result.User = &api.UserShot_Parent{
			Id:       user.User.Id,
			Nickname: user.User.Nickname,
		}
	}
	return
}

// ToListResult ...
func (u User) ToListResult(list []*m.User, total uint64, pagination common.PageParams) *api.GetUserListResult {

	items := make([]*api.UserShot, 0, len(list))

	for _, i := range list {
		items = append(items, u.ToUserShot(i))
	}

	return &api.GetUserListResult{
		Items: items,
		Meta: &api.GetUserListResult_Meta{
			Limit:        uint64(pagination.Limit),
			ObjectsCount: total,
			Offset:       uint64(pagination.Offset),
		},
	}
}

// FromUpdateUserRequest ...
func (u User) FromUpdateUserRequest(req *api.UpdateUserRequest) (user *m.User) {
	user = &m.User{}
	common.Copy(&user, req, common.JsonEngine)
	return
}
