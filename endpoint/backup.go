// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"bufio"
	"context"
	"mime/multipart"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
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

	if b.appConfig.Mode == common.DemoMode {
		err = apperr.ErrBackupCreateNewForbidden
		return
	}

	go b.backup.New(false)
	return
}

// Restore ...
func (b *BackupEndpoint) Restore(ctx context.Context, name string) (err error) {

	if b.appConfig.Mode == common.DemoMode {
		err = apperr.ErrBackupRestoreForbidden
		return
	}

	err = b.backup.Restore(name)
	return
}

// GetList ...
func (b *BackupEndpoint) GetList(ctx context.Context, pagination common.PageParams) (items []*m.Backup, total int64, err error) {
	items, total, err = b.backup.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	return
}

// Upload ...
func (b *BackupEndpoint) Upload(ctx context.Context, files map[string][]*multipart.FileHeader) (fileList []*m.Backup, errs []error) {

	fileList = make([]*m.Backup, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, err := fileHeader[0].Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		reader := bufio.NewReader(file)
		var newbackup *m.Backup
		newbackup, err = b.backup.UploadBackup(ctx, reader, fileHeader[0].Filename)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fileList = append(fileList, newbackup)

		file.Close()
	}

	return
}

func (b *BackupEndpoint) Delete(ctx context.Context, name string) (err error) {

	var list []*m.Backup
	if list, _, err = b.backup.List(ctx, 999, 0, "", ""); err != nil {
		return
	}
	for _, file := range list {
		if name == file.Name {
			err = b.backup.Delete(file.Name)
			return
		}
	}

	err = apperr.ErrBackupNotFound

	return
}

func (b *BackupEndpoint) ApplyChanges(ctx context.Context) (err error) {

	if b.appConfig.Mode == common.DemoMode {
		err = apperr.ErrBackupApplyForbidden
		return
	}

	err = b.backup.ApplyChanges()
	return
}

func (b *BackupEndpoint) RollbackChanges(ctx context.Context) (err error) {

	if b.appConfig.Mode == common.DemoMode {
		err = apperr.ErrBackupRollbackForbidden
		return
	}

	err = b.backup.RollbackChanges()
	return
}
