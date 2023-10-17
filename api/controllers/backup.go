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
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/api/stub"
)

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
func (c ControllerBackup) BackupServiceRestoreBackup(ctx echo.Context, _ stub.BackupServiceRestoreBackupParams) error {

	obj := &stub.ApiRestoreBackupRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	err := c.endpoint.Backup.Restore(ctx.Request().Context(), obj.Name)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// GetBackupList ...
func (c ControllerBackup) BackupServiceGetBackupList(ctx echo.Context) error {

	result := c.endpoint.Backup.GetList(ctx.Request().Context())

	return c.HTTP200(ctx, ResponseWithObj(ctx, &stub.ApiGetBackupListResult{
		Items: result,
	}))
}
