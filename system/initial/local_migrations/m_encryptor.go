package local_migrations

import (
	"context"
	"encoding/hex"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/encryptor"
	m "github.com/e154/smart-home/models"
)

type MigrationEncryptor struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationEncryptor(adaptors *adaptors.Adaptors) *MigrationEncryptor {
	return &MigrationEncryptor{
		adaptors: adaptors,
	}
}

func (n *MigrationEncryptor) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	_, err := n.adaptors.Variable.GetByName(ctx, "encryptor")
	if err != nil {
		err = n.adaptors.Variable.Add(ctx, m.Variable{
			Name:   "encryptor",
			Value:  hex.EncodeToString(encryptor.GenKey()),
			System: true,
		})
		return err
	}

	return nil
}
