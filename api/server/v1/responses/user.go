package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response UserList
type UserList struct {
	// in:body
	Body struct {
		Items []*models.UserShot `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response UserSearch
type UserSearch struct {
	// in:body
	Body struct {
		Users []*models.UserShot `json:"users"`
	}
}
