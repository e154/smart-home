package local_migrations

import (
	"context"

	m "github.com/e154/smart-home/models"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/location"
	. "github.com/e154/smart-home/plugins/zone"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationZones struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationZones(adaptors *adaptors.Adaptors) *MigrationZones {
	return &MigrationZones{
		adaptors: adaptors,
	}
}

func (n *MigrationZones) Up(_ context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	err = n.addZone("home", "base geo position")
	return
}

func (n *MigrationZones) addZone(name, desc string) (err error) {

	id := common.EntityId("zone." + name)
	if _, err = n.adaptors.Entity.GetById(id); err == nil {
		return
	}

	loc, err := location.GetRegionInfo()
	So(err, ShouldBeNil)

	var latitude = loc.Latitude
	var longitude = loc.Longitude

	if latitude == 0 || longitude == 0 {
		loc, err := location.GeoLocationFromIP("")
		So(err, ShouldBeNil)
		latitude = loc.Lat
		longitude = loc.Lon
	}

	attributes := m.Attributes{
		AttrLat: {
			Name:  AttrLat,
			Type:  common.AttributeFloat,
			Value: latitude,
		},
		AttrLon: {
			Name:  AttrLon,
			Type:  common.AttributeFloat,
			Value: longitude,
		},
		AttrElevation: {
			Name:  AttrElevation,
			Type:  common.AttributeFloat,
			Value: 150,
		},
		AttrTimezone: {
			Name:  AttrTimezone,
			Type:  common.AttributeInt,
			Value: 7,
		},
	}

	err = n.adaptors.Entity.Add(&m.Entity{
		Id:          id,
		Description: desc,
		PluginName:  "zone",
		Attributes:  attributes,
		AutoLoad:    true,
	})
	So(err, ShouldBeNil)

	_, err = n.adaptors.EntityStorage.Add(&m.EntityStorage{
		EntityId:   common.EntityId("zone." + name),
		Attributes: attributes.Serialize(),
	})
	So(err, ShouldBeNil)

	return
}
