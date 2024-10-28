// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

	m "github.com/e154/smart-home/pkg/models"
)

// UserRepo ...
type UserRepo interface {
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
}
