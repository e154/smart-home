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
)

type MigrationDashboard struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationDashboard(adaptors *adaptors.Adaptors) *MigrationDashboard {
	return &MigrationDashboard{
		adaptors: adaptors,
	}
}

func (n *MigrationDashboard) addDashboard(ctx context.Context, name, src string) error {

	//req := &api.Dashboard{}
	//_ = json.Unmarshal([]byte(src), req)
	//
	//board := dto.ImportDashboard(req)
	//
	//var err error
	//if board.Id, err = n.adaptors.Dashboard.Import(ctx, board); err != nil {
	//	return err
	//}
	//
	//err = n.adaptors.Variable.CreateOrUpdate(ctx, m.Variable{
	//	Name:   name,
	//	Value:  fmt.Sprintf("%d", board.Id),
	//	System: true,
	//})

	return nil
}

func (n *MigrationDashboard) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	//err := n.addDashboard(ctx, "devDashboard", devDashboardRaw)
	//So(err, ShouldBeNil)
	//
	//err = n.addDashboard(ctx, "mainDashboard", mainDashboardRaw)
	//So(err, ShouldBeNil)

	return nil
}
