package local_migrations

import (
	"context"

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
	"github.com/e154/smart-home/plugins/script"
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
	"github.com/e154/smart-home/plugins/zone"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

type MigrationPlugins struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationPlugins(adaptors *adaptors.Adaptors) *MigrationPlugins {
	return &MigrationPlugins{
		adaptors: adaptors,
	}
}

func (n *MigrationPlugins) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	n.addPlugin("alexa", false, false, true, alexa.Version)
	n.addPlugin("cgminer", false, false, true, cgminer.Version)
	n.addPlugin("cpuspeed", true, false, false, cpuspeed.Version)
	n.addPlugin("email", true, false, false, email.Version)
	n.addPlugin("messagebird", false, false, false, messagebird.Version)
	n.addPlugin("modbus_rtu", false, false, true, modbus_rtu.Version)
	n.addPlugin("modbus_tcp", false, false, true, modbus_tcp.Version)
	n.addPlugin("moon", true, false, true, moon.Version)
	n.addPlugin("memory", true, false, true, memory.Version)
	n.addPlugin("memory_app", true, false, false, memory_app.Version)
	n.addPlugin("hdd", true, false, true, hdd.Version)
	n.addPlugin("logs", true, false, false, logs.Version)
	n.addPlugin("version", true, false, false, version.Version)
	n.addPlugin("node", true, true, true, node.Version)
	n.addPlugin("notify", true, true, false, notify.Version)
	n.addPlugin("scene", true, false, true, scene.Version)
	n.addPlugin("script", true, false, true, script.Version)
	n.addPlugin("sensor", true, false, true, sensor.Version)
	n.addPlugin("slack", true, false, false, slack.Version)
	n.addPlugin("sun", true, false, true, sun.Version)
	n.addPlugin("telegram", true, false, true, telegram.Version)
	n.addPlugin("triggers", true, true, false, triggers.Version)
	n.addPlugin("twilio", false, false, false, twilio.Version)
	n.addPlugin("updater", true, false, false, updater.Version)
	n.addPlugin("uptime", true, false, false, uptime.Version)
	n.addPlugin("weather_met", false, false, true, weather_met.Version)
	n.addPlugin("weather_owm", false, false, true, weather_owm.Version)
	n.addPlugin("mqtt", true, false, true, mqtt.Version)
	n.addPlugin("zigbee2mqtt", false, false, true, zigbee2mqtt.Version)
	n.addPlugin("zone", true, false, true, zone.Version)
	n.addPlugin("html5_notify", true, false, false, html5_notify.Version)
	n.addPlugin("webpush", true, false, false, webpush.Version)
	return nil
}

func (n *MigrationPlugins) addPlugin(name string, enabled, system, actor bool, version string) (node *m.Plugin) {
	_ = n.adaptors.Plugin.CreateOrUpdate(&m.Plugin{
		Name:    name,
		Version: version,
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
	return
}
