package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response TemplateList
type TemplateList struct {
	// in:body
	Body struct {
		Items []*models.Template `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response TemplateSearch
type TemplateSearch struct {
	// in:body
	Body struct {
		Templates []*models.Template `json:"templates"`
	}
}
