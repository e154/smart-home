package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response LogList
type LogList struct {
	// in:body
	Body struct {
		Items []*models.Log `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response LogSearch
type LogSearch struct {
	// in:body
	Body struct {
		Logs []*models.Log `json:"logs"`
	}
}
