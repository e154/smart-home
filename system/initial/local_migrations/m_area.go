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
