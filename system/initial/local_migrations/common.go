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

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

func AddVariableIfNotExist(adaptors *adaptors.Adaptors, ctx context.Context, name, value string) (err error) {

	if _, err = adaptors.Variable.GetByName(ctx, name); err == nil {
		return
	}

	err = adaptors.Variable.Add(context.Background(), m.Variable{
		Name:   name,
		Value:  value,
		System: true,
	})

	return
}

type Common struct {
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
}

func (c *Common) addPlugin(ctx context.Context, name string, enabled, system, actor bool, version string) error {
	return c.adaptors.Plugin.CreateOrUpdate(ctx, &m.Plugin{
		Name:    name,
		Version: version,
		Enabled: enabled,
		System:  system,
		Actor:   actor,
	})
}

func (n *Common) removePlugin(ctx context.Context, name string) error {
	return n.adaptors.Plugin.Delete(ctx, name)
}

func (s *Common) addScript(ctx context.Context, name, source, desc string) (script *m.Script, err error) {

	if script, err = s.adaptors.Script.GetByName(ctx, name); err == nil {
		return
	}

	script = &m.Script{
		Lang:        common.ScriptLangCoffee,
		Name:        name,
		Source:      source,
		Description: desc,
	}

	engineScript, err := s.scriptService.NewEngine(script)

	err = engineScript.Compile()
	So(err, ShouldBeNil)

	script.Id, err = s.adaptors.Script.Add(ctx, script)
	So(err, ShouldBeNil)
	return
}
