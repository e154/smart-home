package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response MapZoneList
type MapZoneList struct {
	// in:body
	Body struct {
		Items []*models.MapZone `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response MapZoneSearch
type MapZoneSearch struct {
	// in:body
	Body struct {
		MapZones []*models.MapZone `json:"zones"`
	}
}
