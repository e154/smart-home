package responses

import (
	m "github.com/e154/smart-home/api/server/v1/models"
)

// Error response
// swagger:response Error
type Error struct {
	// in:body
	Body struct{
		m.Error
	}
}

