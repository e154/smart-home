package local_migrations

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

type MigrationMqtt struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationMqtt(adaptors *adaptors.Adaptors) *MigrationMqtt {
	return &MigrationMqtt{
		adaptors: adaptors,
	}
}

func (n *MigrationMqtt) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	n.addPlugin("mqtt", true, false, true)
	return nil
}

func (n *MigrationMqtt) addPlugin(name string, enabled, system, actor bool) (node *m.Plugin) {
	_ = n.adaptors.Plugin.CreateOrUpdate(&m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
	return
}
