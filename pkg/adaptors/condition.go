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

// ConditionRepo ...
type ConditionRepo interface {
	Add(ctx context.Context, ver *m.Condition) (id int64, err error)
	GetById(ctx context.Context, id int64) (metric *m.Condition, err error)
	Update(ctx context.Context, ver *m.Condition) error
	Delete(ctx context.Context, deviceId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, ids *[]uint64) (list []*m.Condition, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Condition, total int64, err error)
}
