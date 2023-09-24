package local_migrations

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
)

type MigrationEntity struct {
	adaptors *adaptors.Adaptors
	endpoint *endpoint.Endpoint
}

func NewMigrationEntity(adaptors *adaptors.Adaptors, endpoint *endpoint.Endpoint) *MigrationEntity {
	return &MigrationEntity{
		adaptors: adaptors,
		endpoint: endpoint,
	}
}

func (n *MigrationEntity) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	return nil
}
