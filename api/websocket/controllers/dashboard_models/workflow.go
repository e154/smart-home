package dashboard_models

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/telemetry"
)

type Workflow struct {
	Total    int64                               `json:"total"`
	Status   map[int64]m.DashboardWorkflowStatus `json:"status"`
	adaptors *adaptors.Adaptors                  `json:"-"`
	core     *core.Core                          `json:"-"`
}

func NewWorkflow(adaptors *adaptors.Adaptors,
	core *core.Core) *Workflow {
	return &Workflow{
		adaptors: adaptors,
		core:     core,
	}
}

func (w *Workflow) Update() {

}

func (w *Workflow) Broadcast() (map[string]interface{}, bool) {

	w.Update()

	return map[string]interface{}{
		"workflows": w,
	}, true
}

func (w *Workflow) BroadcastOne(params telemetry.WorkflowScenario) (map[string]interface{}, bool) {
	log.Warning("realize workflow scenario broadcast")

	return map[string]interface{}{
		"workflow":  map[string]interface{}{"id": params.WorkflowId, "scenario_id": params.ScenarioId},
		"workflows": w,
	}, true
}
