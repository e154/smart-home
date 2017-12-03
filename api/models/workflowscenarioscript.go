package models

import (
	"github.com/astaxie/beego"
)

type WorkflowScenarioScript struct {
	Id   		int64  			`orm:"pk;auto" json:"id"`
	Scenario	*WorkflowScenario       `orm:"rel(fk);column(workflow_scenario_id)" json:"scenario"`
	Script		*Script			`orm:"rel(fk);column(script_id)" json:"script"`
}

func (m *WorkflowScenarioScript) TableName() string {
	return beego.AppConfig.String("db_workflow_scenario_scripts")
}
