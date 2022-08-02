package local_migrations

import (
	"context"

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
	n.addPlugin("alexa", false, false, false)
	n.addPlugin("cgminer", true, false, true)
	n.addPlugin("cpuspeed", false, false, false)
	n.addPlugin("email", true, false, false)
	n.addPlugin("messagebird", false, false, false)
	n.addPlugin("modbus_rtu", false, false, true)
	n.addPlugin("modbus_tcp", false, false, true)
	n.addPlugin("moon", false, false, true)
	n.addPlugin("memory", true, false, true)
	n.addPlugin("memory_app", true, false, true)
	n.addPlugin("hdd", true, false, true)
	n.addPlugin("logs", true, false, true)
	n.addPlugin("version", true, false, true)
	n.addPlugin("node", true, true, true)
	n.addPlugin("notify", true, true, false)
	n.addPlugin("scene", true, false, true)
	n.addPlugin("script", true, false, true)
	n.addPlugin("sensor", true, false, true)
	n.addPlugin("slack", false, false, false)
	n.addPlugin("sun", false, false, true)
	n.addPlugin("telegram", true, false, true)
	n.addPlugin("triggers", true, true, false)
	n.addPlugin("twilio", false, false, false)
	n.addPlugin("updater", false, false, false)
	n.addPlugin("uptime", false, false, false)
	n.addPlugin("weather", false, false, false)
	n.addPlugin("weather_met", false, false, true)
	n.addPlugin("weather_owm", false, false, true)
	n.addPlugin("zigbee2mqtt", false, false, true)
	n.addPlugin("zone", false, false, true)
	return nil
}

func (n *MigrationPlugins) addPlugin(name string, enabled, system, actor bool) (node *m.Plugin) {
	_ = n.adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
	return
}
