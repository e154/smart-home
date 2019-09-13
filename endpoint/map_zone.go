package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type MapZoneEndpoint struct {
	*CommonEndpoint
}

func NewMapZoneEndpoint(common *CommonEndpoint) *MapZoneEndpoint {
	return &MapZoneEndpoint{
		CommonEndpoint: common,
	}
}

func (n *MapZoneEndpoint) Add(zone *m.MapZone) (result *m.MapZone, errs []*validation.Error, err error) {

	_, errs = zone.Valid()
	if len(errs) > 0 {
		return
	}

	if zone.Id, err = n.adaptors.MapZone.Add(zone); err != nil {
		return
	}

	return
}

func (n *MapZoneEndpoint) Delete(zoneName string) (err error) {

	if _, err = n.adaptors.MapZone.GetByName(zoneName); err != nil {
		return
	}

	err = n.adaptors.MapZone.Delete(zoneName)

	return
}

func (n *MapZoneEndpoint) Search(query string, limit, offset int) (result []*m.MapZone, total int64, err error) {

	result, total, err = n.adaptors.MapZone.Search(query, limit, offset)

	return
}
