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

package example1

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type AreaManager struct {
	adaptors *adaptors.Adaptors
}

func NewAreaManager(adaptors *adaptors.Adaptors) *AreaManager {
	return &AreaManager{
		adaptors: adaptors,
	}
}

func (n *AreaManager) addArea(ctx context.Context, name, desc string) (area *m.Area, err error) {
	if area, err = n.adaptors.Area.GetByName(context.Background(), name); err == nil {
		return
	}
	area = &m.Area{
		Name:        name,
		Description: desc,
	}

	area.Id, err = n.adaptors.Area.Add(ctx, area)
	So(err, ShouldBeNil)
	return
}
