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

package controllers

import (
	"context"
	"github.com/e154/smart-home/api/stub/api"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerBackup ...
type ControllerBackup struct {
	*ControllerCommon
}

// NewControllerBackup ...
func NewControllerBackup(common *ControllerCommon) ControllerBackup {
	return ControllerBackup{
		ControllerCommon: common,
	}
}

// NewBackup ...
func (c ControllerBackup) NewBackup(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	err := c.endpoint.Backup.New(ctx)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}


// RestoreBackup ...
func (c ControllerBackup) RestoreBackup(ctx context.Context, req *api.RestoreBackupRequest) (*emptypb.Empty, error) {

	err := c.endpoint.Backup.Restore(ctx, req.Name)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// GetBackupList ...
func (c ControllerBackup) GetBackupList(ctx context.Context, _ *emptypb.Empty) (*api.GetBackupListResult, error) {

	result := c.endpoint.Backup.GetList(ctx)

	return &api.GetBackupListResult{
		Items: result,
	}, nil
}
