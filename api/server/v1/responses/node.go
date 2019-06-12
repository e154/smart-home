package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response NodeList
type NodeList struct {
	// in:body
	Body struct {
		Items []*models.Node `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response NodeSearch
type NodeSearch struct {
	// in:body
	Body struct {
		Nodes []*models.Node `json:"nodes"`
	}
}
