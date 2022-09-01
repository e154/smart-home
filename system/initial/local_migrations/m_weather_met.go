package local_migrations

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

type MigrationWeatherMet struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationWeatherMet(adaptors *adaptors.Adaptors) *MigrationWeatherMet {
	return &MigrationWeatherMet{
		adaptors: adaptors,
	}
}

func (n *MigrationWeatherMet) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	_ = n.adaptors.Plugin.CreateOrUpdate(m.Plugin{
		Name:    "weather_met",
		Version: "0.0.1",
		Enabled: true,
		System:  true,
		Actor:   true,
	})

	return nil
}
