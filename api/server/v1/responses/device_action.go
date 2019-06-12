package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response DeviceActionList
type DeviceActionList struct {
	// in:body
	Body struct {
		Items []*models.DeviceAction `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response DeviceActionSearch
type DeviceActionSearch struct {
	// in:body
	Body struct {
		DeviceActions []*models.DeviceAction `json:"device_actions"`
	}
}
