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

package adaptors

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IUser ...
type IUser interface {
	Add(ctx context.Context, user *m.User) (id int64, err error)
	GetById(ctx context.Context, userId int64) (user *m.User, err error)
	GetByNickname(ctx context.Context, nick string) (user *m.User, err error)
	GetByEmail(ctx context.Context, email string) (user *m.User, err error)
	GetByAuthenticationToken(ctx context.Context, token string) (user *m.User, err error)
	GetByResetPassToken(ctx context.Context, token string) (user *m.User, err error)
	Update(ctx context.Context, user *m.User) (err error)
	Delete(ctx context.Context, userId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.User, total int64, err error)
	SignIn(ctx context.Context, u *m.User, ipv4 string) (err error)
	GenResetPassToken(ctx context.Context, u *m.User) (token string, err error)
	ClearResetPassToken(ctx context.Context, u *m.User) (err error)
	ClearToken(ctx context.Context, u *m.User) (err error)
	fromDb(dbUser *db.User) (user *m.User)
	toDb(user *m.User) (dbUser *db.User)
}

// User ...
type User struct {
	IUser
	table *db.Users
	db    *gorm.DB
}

// GetUserAdaptor ...
func GetUserAdaptor(d *gorm.DB) IUser {
	return &User{
		table: &db.Users{Db: d},
		db:    d,
	}
}

// Add ...
func (n *User) Add(ctx context.Context, user *m.User) (id int64, err error) {

	dbUser := n.toDb(user)
	_ = dbUser.History.UnmarshalJSON([]byte("[]"))
	if id, err = n.table.Add(ctx, dbUser); err != nil {
		return
	}

	metaAdaptor := GetUserMetaAdaptor(n.db)
	for _, meta := range user.Meta {
		meta.UserId = id
		_, _ = metaAdaptor.UpdateOrCreate(ctx, meta)
	}

	return
}

// GetById ...
func (n *User) GetById(ctx context.Context, userId int64) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetById(ctx, userId); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(ctx, user.Role)

	return
}

// GetByNickname ...
func (n *User) GetByNickname(ctx context.Context, nick string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByNickname(ctx, nick); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(ctx, user.Role)

	return
}

// GetByEmail ...
func (n *User) GetByEmail(ctx context.Context, email string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByEmail(ctx, email); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(ctx, user.Role)

	return
}

// GetByAuthenticationToken ...
func (n *User) GetByAuthenticationToken(ctx context.Context, token string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByAuthenticationToken(ctx, token); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(ctx, user.Role)

	return
}

// GetByResetPassToken ...
func (n *User) GetByResetPassToken(ctx context.Context, token string) (user *m.User, err error) {

	if utf8.RuneCountInString(token) > 255 {
		return
	}

	var dbUser *db.User
	if dbUser, err = n.table.GetByResetPassToken(ctx, token); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	t := time.Now()
	sub := t.Sub(user.ResetPasswordSentAt.Add(time.Hour * 24)).String()
	if !strings.Contains(sub, "-") {
		err = apperr.ErrTokenIsDeprecated
	}

	_ = n.ClearResetPassToken(ctx, user)

	return
}

// Update ...
func (n *User) Update(ctx context.Context, user *m.User) (err error) {

	dbUser := n.toDb(user)
	if err = n.table.Update(ctx, dbUser); err != nil {
		return
	}

	metaAdaptor := GetUserMetaAdaptor(n.db)
	for _, meta := range user.Meta {
		meta.UserId = user.Id
		_, _ = metaAdaptor.UpdateOrCreate(ctx, meta)
	}

	return
}

// Delete ...
func (n *User) Delete(ctx context.Context, userId int64) (err error) {
	err = n.table.Delete(ctx, userId)
	return
}

// List ...
func (n *User) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.User, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.User
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.User, 0)
	for _, dbUser := range dbList {
		user := n.fromDb(dbUser)
		list = append(list, user)
	}

	return
}

// SignIn ...
func (n *User) SignIn(ctx context.Context, u *m.User, ipv4 string) (err error) {

	// update count
	u.SignInCount += 1

	// update time
	lastT := u.CurrentSignInAt
	currentT := time.Now()

	u.LastSignInAt = lastT
	u.CurrentSignInAt = &currentT

	// update ipv4
	lastIp := u.CurrentSignInIp
	currentIp := ipv4
	u.LastSignInIp = lastIp
	u.CurrentSignInIp = currentIp

	u.UpdateHistory(currentT, currentIp)

	dbUser := n.toDb(u)
	err = n.table.Update(ctx, dbUser)

	return
}

// GenResetPassToken ...
func (n *User) GenResetPassToken(ctx context.Context, u *m.User) (token string, err error) {

	for {
		token = common.RandStr(50, common.Alphanum)
		u.ResetPasswordToken = token

		if _, err = n.table.GetByResetPassToken(ctx, token); err != nil {
			break
		}
	}

	err = n.table.NewResetPassToken(ctx, u.Id, u.ResetPasswordToken)

	return
}

// ClearResetPassToken ...
func (n *User) ClearResetPassToken(ctx context.Context, u *m.User) (err error) {

	err = n.table.ClearResetPassToken(ctx, u.Id)
	return
}

// ClearToken ...
func (n *User) ClearToken(ctx context.Context, u *m.User) (err error) {

	err = n.table.ClearToken(ctx, u.Id)

	return
}

func (n *User) fromDb(dbUser *db.User) (user *m.User) {
	user = &m.User{
		Id:                  dbUser.Id,
		Nickname:            dbUser.Nickname,
		FirstName:           dbUser.FirstName,
		LastName:            dbUser.LastName,
		EncryptedPassword:   dbUser.EncryptedPassword,
		Email:               dbUser.Email,
		Status:              dbUser.Status,
		ResetPasswordToken:  dbUser.ResetPasswordToken,
		AuthenticationToken: dbUser.AuthenticationToken,
		ImageId:             dbUser.ImageId,
		SignInCount:         dbUser.SignInCount,
		CurrentSignInIp:     dbUser.CurrentSignInIp,
		LastSignInIp:        dbUser.LastSignInIp,
		Lang:                dbUser.Lang,
		UserId:              dbUser.UserId,
		RoleName:            dbUser.RoleName,
		ResetPasswordSentAt: dbUser.ResetPasswordSentAt,
		CurrentSignInAt:     dbUser.CurrentSignInAt,
		LastSignInAt:        dbUser.LastSignInAt,
		CreatedAt:           dbUser.CreatedAt,
		UpdatedAt:           dbUser.UpdatedAt,
		DeletedAt:           dbUser.DeletedAt,
		Meta:                make([]*m.UserMeta, 0),
	}

	if dbUser.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		user.Image = imageAdaptor.fromDb(dbUser.Image)
	}

	if dbUser.Meta != nil && len(dbUser.Meta) > 0 {
		userMetaAdaptor := GetUserMetaAdaptor(n.db)
		for _, dbMeta := range dbUser.Meta {
			meta := userMetaAdaptor.fromDb(dbMeta)
			user.Meta = append(user.Meta, meta)
		}
	}

	// history
	user.History = make([]*m.UserHistory, 0)
	data, _ := dbUser.History.MarshalJSON()
	_ = json.Unmarshal(data, &user.History)

	// role
	if dbUser.Role != nil {
		roleAdaptor := GetRoleAdaptor(n.db)
		user.Role = roleAdaptor.fromDb(dbUser.Role)
	}

	// created by
	if dbUser.User != nil {
		user.User = n.fromDb(dbUser.User)
	}

	return
}

func (n *User) toDb(user *m.User) (dbUser *db.User) {
	dbUser = &db.User{
		Id:                  user.Id,
		Nickname:            user.Nickname,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		EncryptedPassword:   user.EncryptedPassword,
		Email:               user.Email,
		Status:              user.Status,
		ResetPasswordToken:  user.ResetPasswordToken,
		AuthenticationToken: user.AuthenticationToken,
		ImageId:             user.ImageId,
		SignInCount:         user.SignInCount,
		CurrentSignInIp:     user.CurrentSignInIp,
		LastSignInIp:        user.LastSignInIp,
		Lang:                user.Lang,
		UserId:              user.UserId,
		RoleName:            user.RoleName,
		ResetPasswordSentAt: user.ResetPasswordSentAt,
		CurrentSignInAt:     user.CurrentSignInAt,
		LastSignInAt:        user.LastSignInAt,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		DeletedAt:           user.DeletedAt,
	}

	if user.ImageId.Valid {
		_ = dbUser.ImageId.Scan(user.ImageId.Int64)
	}

	if user.UserId.Valid {
		_ = dbUser.UserId.Scan(user.UserId.Int64)
	}

	return
}
