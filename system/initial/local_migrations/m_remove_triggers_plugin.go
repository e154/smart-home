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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/alexa"
	"github.com/e154/smart-home/plugins/ble"
	"github.com/e154/smart-home/plugins/time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/plugins/state_change"
	"github.com/e154/smart-home/plugins/system"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationRemoveTriggersPlugin struct {
	Common
}

func NewMigrationRemoveTriggersPlugin(adaptors *adaptors.Adaptors) *MigrationRemoveTriggersPlugin {
	return &MigrationRemoveTriggersPlugin{
		Common{
			adaptors: adaptors,
		},
	}
}

func (n *MigrationRemoveTriggersPlugin) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	err := n.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:     "state_change",
		Version:  state_change.Version,
		Enabled:  true,
		System:   true,
		Actor:    false,
		Triggers: true,
	})
	So(err, ShouldBeNil)

	err = n.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:     "system",
		Version:  system.Version,
		Enabled:  true,
		System:   true,
		Actor:    false,
		Triggers: true,
	})
	So(err, ShouldBeNil)

	err = n.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:     "time",
		Version:  time.Version,
		Enabled:  true,
		System:   true,
		Actor:    false,
		Triggers: true,
	})
	So(err, ShouldBeNil)

	err = n.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:     "ble",
		Version:  ble.Version,
		Enabled:  false,
		System:   false,
		Actor:    true,
		Triggers: true,
	})
	So(err, ShouldBeNil)

	err = n.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:     "alexa",
		Version:  alexa.Version,
		Enabled:  false,
		System:   false,
		Actor:    true,
		Triggers: true,
	})
	So(err, ShouldBeNil)

	return nil
}
