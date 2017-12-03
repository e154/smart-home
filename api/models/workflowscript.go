package models

import (
	"github.com/astaxie/beego"
)

type WorkflowScript struct {
	Id   		int64  		`orm:"pk;auto" json:"id"`
	Workflow	*Workflow	`orm:"rel(fk);column(workflow_id)" json:"workflow"`
	Script		*Script		`orm:"rel(fk);column(script_id)" json:"script"`
}

func (m *WorkflowScript) TableName() string {
	return beego.AppConfig.String("db_workflow_scripts")
}
