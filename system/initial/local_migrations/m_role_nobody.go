// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package local_migrations

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationRoleNobody struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationRoleNobody(adaptors *adaptors.Adaptors) *MigrationRoleNobody {
	return &MigrationRoleNobody{
		adaptors: adaptors,
	}
}

func (r *MigrationRoleNobody) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		r.adaptors = adaptors
	}

	_, err = r.add(ctx)

	return
}

func (r *MigrationRoleNobody) add(ctx context.Context) (nobodyRole *m.Role, err error) {

	if nobodyRole, err = r.adaptors.Role.GetByName(ctx, "nobody"); err != nil {
		nobodyRole = &m.Role{
			Name: "nobody",
		}
		err = r.adaptors.Role.Add(ctx, nobodyRole)
		So(err, ShouldBeNil)
	}

	return
}
