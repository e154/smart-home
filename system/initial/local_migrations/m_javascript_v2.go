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
	"strings"

	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

type MigrationJavascriptV2 struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationJavascriptV2(adaptors *adaptors.Adaptors) *MigrationJavascriptV2 {
	return &MigrationJavascriptV2{
		adaptors: adaptors,
	}
}

func (n *MigrationJavascriptV2) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	list, _, err := n.adaptors.Script.List(ctx, 999, 0, "desc", "id", nil)
	So(err, ShouldBeNil)

	var engine *scripts.Engine
	for _, script := range list {
		script.Source = strings.ReplaceAll(script.Source, "SetState", "EntitySetState")
		script.Source = strings.ReplaceAll(script.Source, "SetStateName", "EntitySetStateName")
		script.Source = strings.ReplaceAll(script.Source, "GetState", "EntityGetState")
		script.Source = strings.ReplaceAll(script.Source, "SetAttributes", "EntitySetAttributes")
		script.Source = strings.ReplaceAll(script.Source, "GetAttributes", "EntityGetAttributes")
		script.Source = strings.ReplaceAll(script.Source, "GetSettings", "EntityGetSettings")
		script.Source = strings.ReplaceAll(script.Source, "SetMetric", "EntitySetMetric")
		script.Source = strings.ReplaceAll(script.Source, "CallAction", "EntityCallAction")
		script.Source = strings.ReplaceAll(script.Source, "CallScene", "EntityCallScene")

		engine, err = scripts.NewEngine(script, nil, nil)
		So(err, ShouldBeNil)

		err = engine.Compile()
		So(err, ShouldBeNil)

		err = n.adaptors.Script.Update(ctx, script)
		So(err, ShouldBeNil)
	}

	return nil
}
