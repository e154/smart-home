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

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/version"
)

type MigrationPlugins struct {
	Common
}

func NewMigrationPlugins(adaptors *adaptors.Adaptors) *MigrationPlugins {
	return &MigrationPlugins{
		Common{
			adaptors: adaptors,
		},
	}
}

func (n *MigrationPlugins) Up(ctx context.Context) error {

	n.addPlugin(ctx, "alexa", false, false, true, version.VersionString)
	n.addPlugin(ctx, "cgminer", false, false, true, version.VersionString)
	n.addPlugin(ctx, "cpuspeed", true, false, false, version.VersionString)
	n.addPlugin(ctx, "email", true, false, true, version.VersionString)
	n.addPlugin(ctx, "messagebird", false, false, true, version.VersionString)
	n.addPlugin(ctx, "modbus_rtu", false, false, true, version.VersionString)
	n.addPlugin(ctx, "modbus_tcp", false, false, true, version.VersionString)
	n.addPlugin(ctx, "moon", true, false, true, version.VersionString)
	n.addPlugin(ctx, "memory", true, false, false, version.VersionString)
	n.addPlugin(ctx, "memory_app", true, false, false, version.VersionString)
	n.addPlugin(ctx, "hdd", true, false, true, version.VersionString)
	n.addPlugin(ctx, "logs", true, false, false, version.VersionString)
	n.addPlugin(ctx, "version", true, false, false, version.VersionString)
	n.addPlugin(ctx, "node", true, true, true, version.VersionString)
	n.addPlugin(ctx, "notify", true, true, false, version.VersionString)
	n.addPlugin(ctx, "scene", true, false, true, version.VersionString)
	n.addPlugin(ctx, "sensor", true, false, true, version.VersionString)
	n.addPlugin(ctx, "slack", true, false, true, version.VersionString)
	n.addPlugin(ctx, "sun", true, false, true, version.VersionString)
	n.addPlugin(ctx, "telegram", true, false, true, version.VersionString)
	n.addPlugin(ctx, "triggers", true, true, false, version.VersionString)
	n.addPlugin(ctx, "twilio", false, false, true, version.VersionString)
	n.addPlugin(ctx, "updater", true, false, false, version.VersionString)
	n.addPlugin(ctx, "uptime", true, false, false, version.VersionString)
	n.addPlugin(ctx, "weather_met", false, false, true, version.VersionString)
	n.addPlugin(ctx, "weather_owm", false, false, true, version.VersionString)
	n.addPlugin(ctx, "mqtt", true, false, true, version.VersionString)
	n.addPlugin(ctx, "zigbee2mqtt", false, false, true, version.VersionString)
	n.addPlugin(ctx, "html5_notify", true, false, false, version.VersionString)
	n.addPlugin(ctx, "webpush", true, false, false, version.VersionString)
	n.addPlugin(ctx, "onvif", false, false, true, version.VersionString)
	return nil
}
