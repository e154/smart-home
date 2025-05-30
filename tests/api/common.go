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

package api

import (
	"context"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"
)

// AddPlugin ...
func AddPlugin(adaptors *adaptors.Adaptors, name string, opts ...models.AttributeValue) (err error) {
	plugin := &models.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: true,
		System:  true,
	}
	if len(opts) > 0 {
		plugin.Settings = opts[0]
	}
	err = adaptors.Plugin.CreateOrUpdate(context.Background(), plugin)
	return
}

// AddArea ...
func AddArea(adaptors *adaptors.Adaptors, name string, _ ...models.Attributes) (area *models.Area, err error) {
	area = &models.Area{
		Name:        name,
		Description: "description " + name,
	}

	area.Id, err = adaptors.Area.Add(context.Background(), area)
	return
}

// AddScript ...
func AddScript(name, src string, adaptors *adaptors.Adaptors, scriptService scripts.ScriptService) (script *models.Script, err error) {

	script = &models.Script{
		Lang:        common.ScriptLangCoffee,
		Name:        name,
		Source:      src,
		Description: "description " + name,
	}

	var engine scripts.Engine
	if engine, err = scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	script.Id, err = adaptors.Script.Add(context.Background(), script)

	return
}
