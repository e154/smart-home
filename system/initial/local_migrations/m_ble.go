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
	"github.com/e154/smart-home/plugins/ble"
)

type MigrationBle struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationBle(adaptors *adaptors.Adaptors) *MigrationBle {
	return &MigrationBle{
		adaptors: adaptors,
	}
}

func (n *MigrationBle) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	n.addPlugin("ble", false, false, true, ble.Version)

	return nil
}

func (n *MigrationBle) addPlugin(name string, enabled, system, actor bool, version string) (node *m.Plugin) {
	_ = n.adaptors.Plugin.CreateOrUpdate(context.Background(), &m.Plugin{
		Name:    name,
		Version: version,
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
	return
}
