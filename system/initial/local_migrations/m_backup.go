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
	. "github.com/e154/smart-home/system/initial/assertions"
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

	err := AddVariableIfNotExist(n.adaptors, ctx, "createBackupAt", "0 0 0 * * *")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "maximumNumberOfBackups", fmt.Sprintf("%d", 60))
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "sendbackuptoTelegramBot", "")
	So(err, ShouldBeNil)
	err = AddVariableIfNotExist(n.adaptors, ctx, "createBackupAt", fmt.Sprintf("%d", 40))
	So(err, ShouldBeNil)

	return nil
}
