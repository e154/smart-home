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
	"fmt"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

type MigrationBackup struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationBackup(adaptors *adaptors.Adaptors) *MigrationBackup {
	return &MigrationBackup{
		adaptors: adaptors,
	}
}

func (n *MigrationBackup) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	if _, err := n.adaptors.Variable.GetByName(ctx, "createBackupAt"); err != nil {
		_ = n.adaptors.Variable.Add(context.Background(), m.Variable{
			Name:   "createBackupAt",
			Value:  "0 0 0 * * *",
			System: true,
		})
	}

	if _, err := n.adaptors.Variable.GetByName(ctx, "maximumNumberOfBackups"); err != nil {
		_ = n.adaptors.Variable.Add(context.Background(), m.Variable{
			Name:   "maximumNumberOfBackups",
			Value:  fmt.Sprintf("%d", 60),
			System: true,
		})
	}

	if _, err := n.adaptors.Variable.GetByName(ctx, "sendbackuptoTelegramBot"); err != nil {
		_ = n.adaptors.Variable.Add(context.Background(), m.Variable{
			Name:   "sendbackuptoTelegramBot",
			Value:  "",
			System: true,
		})
	}

	if _, err := n.adaptors.Variable.GetByName(ctx, "sendTheBackupInPartsMb"); err != nil {
		_ = n.adaptors.Variable.Add(context.Background(), m.Variable{
			Name:   "sendTheBackupInPartsMb",
			Value:  fmt.Sprintf("%d", 40),
			System: true,
		})
	}

	return nil
}
