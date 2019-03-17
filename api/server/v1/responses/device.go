package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response DeviceList
type DeviceList struct {
	// in:body
	Body struct {
		Items []*models.Device `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response DeviceSearch
type DeviceSearch struct {
	// in:body
	Body struct {
		Devices []*models.DeviceShort `json:"devices"`
	}
}
