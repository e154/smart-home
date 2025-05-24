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

// TaskRepo ...
type TaskRepo interface {
	Add(ctx context.Context, ver *m.NewTask) (id int64, err error)
	Import(ctx context.Context, ver *m.Task) (err error)
	Update(ctx context.Context, ver *m.Task) (err error)
	Delete(ctx context.Context, id int64) (err error)
	GetById(ctx context.Context, id int64) (task *m.Task, err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*m.Task, total int64, err error)
	Enable(ctx context.Context, id int64) (err error)
	Disable(ctx context.Context, id int64) (err error)
	DeleteTrigger(ctx context.Context, taskID int64) error
	DeleteCondition(ctx context.Context, taskID int64) error
	DeleteAction(ctx context.Context, taskID int64) error
}
