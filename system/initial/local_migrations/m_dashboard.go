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
	"fmt"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationDashboard struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationDashboard(adaptors *adaptors.Adaptors) *MigrationDashboard {
	return &MigrationDashboard{
		adaptors: adaptors,
	}
}

func (n *MigrationDashboard) addDashboard(ctx context.Context, name, _n, _d string) error {

	if _, err := n.adaptors.Variable.GetByName(ctx, name); err == nil {
		return nil
	}

	board := &m.Dashboard{
		Name:        _n,
		Description: _d,
	}

	var err error
	if board.Id, err = n.adaptors.Dashboard.Add(ctx, board); err != nil {
		return err
	}

	err = n.adaptors.Variable.Update(ctx, m.Variable{
		Name:   name,
		Value:  fmt.Sprintf("%d", board.Id),
		System: true,
	})

	return err
}

func (n *MigrationDashboard) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	err := n.addDashboard(ctx, "devDashboardLight", "develop (light theme)", "DEVELOP")
	So(err, ShouldBeNil)

	err = n.addDashboard(ctx, "devDashboardDark", "develop (dark theme)", "DEVELOP")
	So(err, ShouldBeNil)

	err = n.addDashboard(ctx, "mainDashboardLight", "main (light theme)", "MAIN")
	So(err, ShouldBeNil)

	err = n.addDashboard(ctx, "mainDashboardDark", "main (dark theme)", "MAIN")
	So(err, ShouldBeNil)

	return nil
}
