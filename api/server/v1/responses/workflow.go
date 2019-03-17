package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response WorkflowList
type WorkflowList struct {
	// in:body
	Body struct {
		Items []*models.Workflow `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response WorkflowSearch
type WorkflowSearch struct {
	// in:body
	Body struct {
		Workflows []*models.Workflow `json:"workflows"`
	}
}

// swagger:response WorkflowScenarioSearch
type WorkflowScenarioSearch struct {
	// in:body
	Body struct {
		Scenarios []*models.WorkflowScenario `json:"scenarios"`
	}
}

// swagger:response WorkflowScenarios
type WorkflowScenarios struct {
	// in:body
	Body struct {
		Scenarios []*models.WorkflowScenario `json:"scenarios"`
	}
}
