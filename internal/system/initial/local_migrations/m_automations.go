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

	"github.com/e154/smart-home/internal/endpoint"
	"github.com/e154/smart-home/pkg/adaptors"
)

type MigrationAutomations struct {
	adaptors *adaptors.Adaptors
	endpoint *endpoint.Endpoint
}

func NewMigrationAutomations(adaptors *adaptors.Adaptors, endpoint *endpoint.Endpoint) *MigrationAutomations {
	return &MigrationAutomations{
		adaptors: adaptors,
		endpoint: endpoint,
	}
}

func (n *MigrationAutomations) Up(ctx context.Context) (err error) {
	//if adaptors != nil {
	//	n.adaptors = adaptors
	//}
	//
	//d := dto.NewDto()
	//for _, raw := range []string{} {
	//	req := &api.Task{}
	//	if err = json.Unmarshal([]byte(raw), req); err != nil {
	//		return
	//	}
	//	task := d.Automation.ImportTask(req)
	//	if _, _, err = n.endpoint.Task.Import(ctx, task); err != nil {
	//		return err
	//	}
	//}

	return
}
