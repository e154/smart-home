package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response ScriptList
type ScriptList struct {
	// in:body
	Body struct {
		Items []*models.Script `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response ScriptSearch
type ScriptSearch struct {
	// in:body
	Body struct {
		Scripts []*models.Script `json:"scripts"`
	}
}

// swagger:response ScriptExec
type ScriptExec struct {
	// in:body
	Body struct {
		Result string `json:"result"`
	}
}
