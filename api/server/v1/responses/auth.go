package responses

import (
	m "github.com/e154/smart-home/api/server/v1/models"
)

// AuthSignInResponse response
// swagger:response AuthSignInResponse
type AuthSignInResponse struct {
	// in:body
	Body struct {
		m.AuthSignInResponse
	}
}

// AccessList response
// swagger:response AccessList
type AccessList struct {
	// in:body
	Body struct {
		AccessList m.AccessList `json:"access_list"`
	}
}
