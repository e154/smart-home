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

func (n *MigrationAreas) addArea(name, desc string) (node *m.Area, err error) {
	_, err = n.adaptors.Area.Add(&m.Area{
		Name:        name,
		Description: desc,
	})
	So(err, ShouldBeNil)
	return
}
