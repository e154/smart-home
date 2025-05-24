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

	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/google/uuid"
)

type MigrationGate struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationGate(adaptors *adaptors.Adaptors) *MigrationGate {
	return &MigrationGate{
		adaptors: adaptors,
	}
}

func (n *MigrationGate) Up(ctx context.Context) error {

	err := AddVariableIfNotExist(n.adaptors, ctx, "gateClientId", uuid.NewString())
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientSecretKey", "")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientServerHost", "gate.e154.ru")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientServerPort", "8443")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientPoolIdleSize", "1")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientPoolMaxSize", "100")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "gateClientTLS", "true")
	So(err, ShouldBeNil)

	return nil
}
