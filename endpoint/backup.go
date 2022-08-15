// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/system/backup"
)

// BackupEndpoint ...
type BackupEndpoint struct {
	*CommonEndpoint
	backup *backup.Backup
}

// NewBackupEndpoint ...
func NewBackupEndpoint(common *CommonEndpoint, backup *backup.Backup) *BackupEndpoint {
	return &BackupEndpoint{
		CommonEndpoint: common,
		backup:         backup,
	}
}

// New ...
func (b *BackupEndpoint) New(ctx context.Context) (err error) {
	err = b.backup.New()
	return
}

// Restore ...
func (b *BackupEndpoint) Restore(ctx context.Context, name string) (err error) {
	err = b.backup.Restore(name)
	return
}

// GetList ...
func (b *BackupEndpoint) GetList(ctx context.Context) (list []string) {
	list = b.backup.List()
	return
}
