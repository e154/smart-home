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
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationAreas struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationAreas(adaptors *adaptors.Adaptors) *MigrationAreas {
	return &MigrationAreas{
		adaptors: adaptors,
	}
}

func (n *MigrationAreas) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	if _, err = n.addArea("living_room", "Гостинная"); err != nil {
		return
	}
	if _, err = n.addArea("bedroom", "Спальня"); err != nil {
		return
	}
	_, err = n.addArea("kitchen", "Кухня")

	return
}

func (n *MigrationAreas) addArea(name, desc string) (area *m.Area, err error) {
	if area, err = n.adaptors.Area.GetByName(context.Background(), name); err == nil {
		return
	}
	area = &m.Area{
		Name:        name,
		Description: desc,
	}
	area.Id, err = n.adaptors.Area.Add(context.Background(), area)
	So(err, ShouldBeNil)
	return
}
