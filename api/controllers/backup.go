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

package controllers

import (
	"github.com/e154/smart-home/common/apperr"
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/stub"
)

const maxMemory = 128 << 20

// ControllerBackup ...
type ControllerBackup struct {
	*ControllerCommon
}

// NewControllerBackup ...
func NewControllerBackup(common *ControllerCommon) *ControllerBackup {
	return &ControllerBackup{
		ControllerCommon: common,
	}
}

// NewBackup ...
func (c ControllerBackup) BackupServiceNewBackup(ctx echo.Context, _ stub.BackupServiceNewBackupParams) error {

	err := c.endpoint.Backup.New(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// RestoreBackup ...
func (c ControllerBackup) BackupServiceRestoreBackup(ctx echo.Context, name string) error {

	err := c.endpoint.Backup.Restore(ctx.Request().Context(), name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// GetBackupList ...
func (c ControllerBackup) BackupServiceGetBackupList(ctx echo.Context, params stub.BackupServiceGetBackupListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Backup.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Backup.ToBackupListResult(items), total, pagination))
}

// DeleteBackup ...
func (c ControllerBackup) BackupServiceDeleteBackup(ctx echo.Context, name string) error {

	if err := c.endpoint.Backup.Delete(ctx.Request().Context(), name); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// BackupServiceMuxUploadBackup ...
func (c ControllerBackup) BackupServiceUploadBackup(ctx echo.Context, _ stub.BackupServiceUploadBackupParams) error {

	r := ctx.Request()

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Error(err.Error())
	}

	form := r.MultipartForm
	if len(form.File) == 0 {
		return c.ERROR(ctx, apperr.ErrInvalidRequest)
	}

	list, errs := c.endpoint.Backup.Upload(r.Context(), form.File)

	var resultBackups = make([]interface{}, 0)

	for _, file := range list {
		resultBackups = append(resultBackups, map[string]string{
			"name": file.Name,
		})
	}

	return c.HTTP200(ctx, map[string]interface{}{
		"files":  resultBackups,
		"errors": errs,
	})
}

// BackupServiceApplyChanges ...
func (c ControllerBackup) BackupServiceApplyState(ctx echo.Context, _ stub.BackupServiceApplyStateParams) error {
	err := c.endpoint.Backup.ApplyChanges(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// BackupServiceRevertState ...
func (c ControllerBackup) BackupServiceRevertState(ctx echo.Context, _ stub.BackupServiceRevertStateParams) error {
	err := c.endpoint.Backup.RollbackChanges(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}
