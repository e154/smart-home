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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// User ...
type User struct{}

// NewUserDto ...
func NewUserDto() User {
	return User{}
}

// AddUserRequest ...
func (u User) AddUserRequest(req *stub.ApiNewtUserRequest) (user *m.User) {
	user = &m.User{}
	_ = common.Copy(&user, req, common.JsonEngine)
	if req.ImageId != nil {
		_ = user.ImageId.Scan(*req.ImageId)
	}
	return
}

// ToUserFull ...
func (u User) ToUserFull(user *m.User) (result *stub.ApiUserFull) {
	roleDto := NewRoleDto()
	imageDto := NewImageDto()
	result = &stub.ApiUserFull{
		Id:                  user.Id,
		Nickname:            user.Nickname,
		FirstName:           common.String(user.FirstName),
		LastName:            common.String(user.LastName),
		Email:               user.Email,
		Status:              user.Status,
		SignInCount:         user.SignInCount,
		Role:                roleDto.GetStubRole(user.Role),
		RoleName:            user.RoleName,
		Lang:                user.Lang,
		AuthenticationToken: common.StringValue(user.AuthenticationToken),
		CurrentSignInAt:     user.CurrentSignInAt,
		LastSignInAt:        user.LastSignInAt,
		ResetPasswordSentAt: user.ResetPasswordSentAt,
		DeletedAt:           user.DeletedAt,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}

	if user.LastSignInIp != "" {
		result.LastSignInIp = common.String(user.LastSignInIp)
	}

	if user.CurrentSignInIp != "" {
		result.CurrentSignInIp = common.String(user.CurrentSignInIp)
	}

	// history
	if user.History != nil {
		result.History = make([]stub.ApiUserHistory, 0, len(user.History))
		for _, h := range user.History {
			result.History = append(result.History, stub.ApiUserHistory{
				Ip:   h.Ip,
				Time: h.Time,
			})
		}
	}

	// meta
	if user.Meta != nil {
		result.Meta = make([]stub.ApiUserMeta, 0, len(user.Meta))
		for _, m := range user.Meta {
			result.Meta = append(result.Meta, stub.ApiUserMeta{
				Key:   m.Key,
				Value: m.Value,
			})
		}
	}

	// image
	if user.Image != nil {
		result.Image = imageDto.ToImage(user.Image)
	}

	// parent
	if user.User != nil {
		result.User = &stub.ApiUserFullParent{
			Id:       user.User.Id,
			Nickname: user.User.Nickname,
		}
	}
	return
}

// ToUserShot ...
func (u User) ToUserShot(user *m.User) (result *stub.ApiUserShot) {

	roleDto := NewRoleDto()
	imageDto := NewImageDto()
	result = &stub.ApiUserShot{
		Id:        user.Id,
		Nickname:  user.Nickname,
		FirstName: common.String(user.FirstName),
		LastName:  common.String(user.LastName),
		Email:     user.Email,
		Status:    user.Status,
		Lang:      user.Lang,
		Role:      roleDto.GetStubRole(user.Role),
		RoleName:  user.RoleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// image
	if user.Image != nil {
		result.Image = imageDto.ToImage(user.Image)
	}

	// parent
	if user.User != nil {
		result.User = &stub.ApiUserShotParent{
			Id:       user.User.Id,
			Nickname: user.User.Nickname,
		}
	}
	return
}

// ToListResult ...
func (u User) ToListResult(list []*m.User) []*stub.ApiUserShot {

	items := make([]*stub.ApiUserShot, 0, len(list))

	for _, i := range list {
		items = append(items, u.ToUserShot(i))
	}

	return items
}

// UpdateUserByIdRequest ...
func (u User) UpdateUserByIdRequest(req *stub.UserServiceUpdateUserByIdJSONBody, id int64) (user *m.User) {
	user = &m.User{}
	_ = common.Copy(&user, req, common.JsonEngine)
	if req.ImageId != nil {
		_ = user.ImageId.Scan(*req.ImageId)
	}
	user.Id = id
	return
}
