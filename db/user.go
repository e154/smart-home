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

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Users ...
type Users struct {
	Db *gorm.DB
}

// User ...
type User struct {
	Id                  int64 `gorm:"primary_key"`
	Nickname            string
	FirstName           string
	LastName            string
	EncryptedPassword   string
	Email               string
	Status              string
	ResetPasswordToken  string
	AuthenticationToken *string
	Image               *Image
	ImageId             sql.NullInt64
	SignInCount         int64
	CurrentSignInIp     string
	LastSignInIp        string
	Lang                string
	User                *User
	UserId              sql.NullInt64
	Role                *Role
	RoleName            string
	Meta                []*UserMeta `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE;"`
	ResetPasswordSentAt *time.Time
	CurrentSignInAt     *time.Time
	LastSignInAt        *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
	History             json.RawMessage `gorm:"type:jsonb;not null"`
}

// TableName ...
func (u *User) TableName() string {
	return "users"
}

// Add ...
func (u *Users) Add(user *User) (id int64, err error) {

	if err = u.Db.Create(&user).Error; err != nil {
		err = errors.Wrap(apperr.ErrUserAdd, err.Error())
		return
	}
	id = user.Id
	return
}

// GetById ...
func (u *Users) GetById(userId int64) (user *User, err error) {

	user = &User{}
	err = u.Db.Model(user).
		Where("id = ? and deleted_at isnull", userId).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Find(user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrUserNotFound, fmt.Sprintf("id \"%d\"", userId))
			return
		}
		err = errors.Wrap(apperr.ErrUserGet, err.Error())
	}
	return
}

// GetByEmail ...
func (u *Users) GetByEmail(email string) (user *User, err error) {

	user = &User{}
	err = u.Db.Model(user).
		Where("email = ?", email).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Find(user).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserGet, err.Error())
	}
	return
}

// GetByNickname ...
func (u *Users) GetByNickname(nickname string) (user *User, err error) {

	user = &User{}
	err = u.Db.Model(user).
		Where("nickname = ?", nickname).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Find(user).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserGet, err.Error())
	}
	return
}

// GetByAuthenticationToken ...
func (u *Users) GetByAuthenticationToken(token string) (user *User, err error) {

	user = &User{}
	err = u.Db.Model(user).
		Where("authentication_token = ?", token).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Find(user).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserGet, err.Error())
	}
	return
}

// GetByResetPassToken ...
func (u *Users) GetByResetPassToken(token string) (user *User, err error) {

	user = &User{}
	err = u.Db.Model(user).
		Where("reset_password_token = ?", token).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Find(user).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserGet, err.Error())
	}
	return
}

// Update ...
func (u *Users) Update(user *User) (err error) {

	q := map[string]interface{}{
		"nickname":               user.Nickname,
		"first_name":             user.FirstName,
		"last_name":              user.LastName,
		"email":                  user.Email,
		"status":                 user.Status,
		"reset_password_token":   user.ResetPasswordToken,
		"authentication_token":   user.AuthenticationToken,
		"image_id":               user.ImageId,
		"sign_in_count":          user.SignInCount,
		"current_sign_in_ip":     user.CurrentSignInIp,
		"last_sign_in_ip":        user.LastSignInIp,
		"lang":                   user.Lang,
		"user_id":                user.UserId,
		"role_name":              user.RoleName,
		//"meta":                   user.Meta, //todo fix
		"reset_password_sent_at": user.ResetPasswordSentAt,
		"current_sign_in_at":     user.CurrentSignInAt,
		"last_sign_in_at":        user.LastSignInAt,
		"created_at":             user.CreatedAt,
		"updated_at":             user.UpdatedAt,
		"deleted_at":             user.DeletedAt,
	}
	if user.EncryptedPassword != "" {
		q["encrypted_password"] = user.EncryptedPassword
	}
	if err = u.Db.Model(&User{Id: user.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrUserUpdate, err.Error())
	}
	return
}

// NewResetPassToken ...
func (u *Users) NewResetPassToken(userId int64, token string) (err error) {
	err = u.Db.Model(&User{Id: userId}).Updates(map[string]interface{}{
		"reset_password_token":   token,
		"reset_password_sent_at": time.Now(),
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserUpdate, err.Error())
	}
	return
}

// ClearResetPassToken ...
func (u *Users) ClearResetPassToken(userId int64) (err error) {
	err = u.Db.Model(&User{Id: userId}).Updates(map[string]interface{}{
		"reset_password_token":   "",
		"reset_password_sent_at": nil,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserUpdate, err.Error())
	}
	return
}

// ClearToken ...
func (u *Users) ClearToken(userId int64) (err error) {
	err = u.Db.Model(&User{Id: userId}).Updates(map[string]interface{}{
		"authentication_token": "",
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserUpdate, err.Error())
	}
	return
}

// UpdateAuthenticationToken ...
func (u *Users) UpdateAuthenticationToken(userId int64, token string) (err error) {
	err = u.Db.Model(&User{Id: userId}).Updates(map[string]interface{}{
		"authentication_token": token,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserUpdate, err.Error())
	}
	return
}

// Delete ...
func (u *Users) Delete(userId int64) (err error) {
	err = u.Db.Model(&User{Id: userId}).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserDelete, err.Error())
	}
	return
}

// List ...
func (n *Users) List(limit, offset int, orderBy, sort string) (list []*User, total int64, err error) {

	if err = n.Db.Model(User{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrUserList, err.Error())
		return
	}

	list = make([]*User, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Preload("Image").
		Preload("Meta").
		Preload("Role").
		Preload("User").
		Where("deleted_at isnull").
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrUserList, err.Error())
	}
	return
}
