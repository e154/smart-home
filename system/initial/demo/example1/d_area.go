package example1

import (
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

func (n *AreaManager) addArea(name, desc string) (area *m.Area, err error) {
	if area, err = n.adaptors.Area.GetByName(name); err == nil {
		return
	}
	area = &m.Area{
		Name:        name,
		Description: desc,
	}

	area.Id, err = n.adaptors.Area.Add(area)
	So(err, ShouldBeNil)
	return
}
