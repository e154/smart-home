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

package container

import (
	demo2 "github.com/e154/smart-home/internal/system/initial/demo"
	"github.com/e154/smart-home/internal/system/initial/demo/example1"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/scripts"
)

func NewDemo(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService) (d *demo2.Demos) {
	list := make(map[string]demo2.Demo)
	list["example1"] = example1.NewExample1(adaptors, scriptService)
	d = demo2.NewDemos(list)
	return
}
