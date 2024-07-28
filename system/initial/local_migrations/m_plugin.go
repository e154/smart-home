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

	"github.com/e154/smart-home/plugins/onvif"

	"github.com/e154/smart-home/plugins/alexa"
	"github.com/e154/smart-home/plugins/cgminer"
	"github.com/e154/smart-home/plugins/cpuspeed"
	"github.com/e154/smart-home/plugins/email"
	"github.com/e154/smart-home/plugins/hdd"
	"github.com/e154/smart-home/plugins/html5_notify"
	"github.com/e154/smart-home/plugins/logs"
	"github.com/e154/smart-home/plugins/memory"
	"github.com/e154/smart-home/plugins/memory_app"
	"github.com/e154/smart-home/plugins/messagebird"
	"github.com/e154/smart-home/plugins/modbus_rtu"
	"github.com/e154/smart-home/plugins/modbus_tcp"
	"github.com/e154/smart-home/plugins/moon"
	"github.com/e154/smart-home/plugins/mqtt"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/scene"
	"github.com/e154/smart-home/plugins/sensor"
	"github.com/e154/smart-home/plugins/slack"
	"github.com/e154/smart-home/plugins/sun"
	"github.com/e154/smart-home/plugins/telegram"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/plugins/twilio"
	"github.com/e154/smart-home/plugins/updater"
	"github.com/e154/smart-home/plugins/uptime"
	"github.com/e154/smart-home/plugins/version"
	"github.com/e154/smart-home/plugins/weather_met"
	"github.com/e154/smart-home/plugins/weather_owm"
	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/plugins/zigbee2mqtt"

	"github.com/e154/smart-home/adaptors"
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

func (n *MigrationPlugins) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	n.addPlugin(ctx, "alexa", false, false, true, alexa.Version)
	n.addPlugin(ctx, "cgminer", false, false, true, cgminer.Version)
	n.addPlugin(ctx, "cpuspeed", true, false, false, cpuspeed.Version)
	n.addPlugin(ctx, "email", true, false, true, email.Version)
	n.addPlugin(ctx, "messagebird", false, false, true, messagebird.Version)
	n.addPlugin(ctx, "modbus_rtu", false, false, true, modbus_rtu.Version)
	n.addPlugin(ctx, "modbus_tcp", false, false, true, modbus_tcp.Version)
	n.addPlugin(ctx, "moon", true, false, true, moon.Version)
	n.addPlugin(ctx, "memory", true, false, false, memory.Version)
	n.addPlugin(ctx, "memory_app", true, false, false, memory_app.Version)
	n.addPlugin(ctx, "hdd", true, false, true, hdd.Version)
	n.addPlugin(ctx, "logs", true, false, false, logs.Version)
	n.addPlugin(ctx, "version", true, false, false, version.Version)
	n.addPlugin(ctx, "node", true, true, true, node.Version)
	n.addPlugin(ctx, "notify", true, true, false, notify.Version)
	n.addPlugin(ctx, "scene", true, false, true, scene.Version)
	n.addPlugin(ctx, "sensor", true, false, true, sensor.Version)
	n.addPlugin(ctx, "slack", true, false, true, slack.Version)
	n.addPlugin(ctx, "sun", true, false, true, sun.Version)
	n.addPlugin(ctx, "telegram", true, false, true, telegram.Version)
	n.addPlugin(ctx, "triggers", true, true, false, triggers.Version)
	n.addPlugin(ctx, "twilio", false, false, true, twilio.Version)
	n.addPlugin(ctx, "updater", true, false, false, updater.Version)
	n.addPlugin(ctx, "uptime", true, false, false, uptime.Version)
	n.addPlugin(ctx, "weather_met", false, false, true, weather_met.Version)
	n.addPlugin(ctx, "weather_owm", false, false, true, weather_owm.Version)
	n.addPlugin(ctx, "mqtt", true, false, true, mqtt.Version)
	n.addPlugin(ctx, "zigbee2mqtt", false, false, true, zigbee2mqtt.Version)
	n.addPlugin(ctx, "html5_notify", true, false, false, html5_notify.Version)
	n.addPlugin(ctx, "webpush", true, false, false, webpush.Version)
	n.addPlugin(ctx, "onvif", false, false, true, onvif.Version)
	return nil
}
