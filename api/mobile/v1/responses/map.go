package responses

import (
	"github.com/e154/smart-home/api/mobile/v1/models"
)

// swagger:response MapActiveElementList
type MapActiveElementList struct {
	// in:body
	Body struct {
		Items []*models.MapElement `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}
