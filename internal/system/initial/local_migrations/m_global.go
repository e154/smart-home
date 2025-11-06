// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023-2025, Filippov Alex
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
	"io"

	scriptsAssets "github.com/e154/smart-home/data/scripts"
	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

type MigrationGlobalScripts struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationGlobalScripts(adaptors *adaptors.Adaptors) *MigrationGlobalScripts {
	return &MigrationGlobalScripts{
		adaptors: adaptors,
	}
}

func (s *MigrationGlobalScripts) Up(ctx context.Context) error {

	script, err := s.adaptors.Script.GetByName(ctx, "global.d")
	if script != nil || err == nil {
		return nil
	}

	file, err := scriptsAssets.F.Open("global.d.ts")
	So(err, ShouldBeNil)
	defer file.Close()

	data, err := io.ReadAll(file)
	So(err, ShouldBeNil)

	script = &m.Script{
		Lang:        common.ScriptLangTs,
		Name:        "global.d",
		Source:      string(data),
		Description: "global.d.ts",
	}
	_, err = s.adaptors.Script.Add(ctx, script)
	So(err, ShouldBeNil)

	return nil
}
