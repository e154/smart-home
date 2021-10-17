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

package example1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

// AreaManager ...
type AreaManager struct {
	adaptors *adaptors.Adaptors
}

// NewAreaManager ...
func NewAreaManager(adaptors *adaptors.Adaptors) *AreaManager {
	return &AreaManager{
		adaptors: adaptors,
	}
}

func (a *AreaManager) create() []*m.Area {
	area1 := a.add("zone51")
	return []*m.Area{area1}
}

func (a *AreaManager) add(name string) *m.Area {
	area := &m.Area{
		Name:        name,
	}
	area.Id, _ = a.adaptors.Area.Add(area)
	return area
}
