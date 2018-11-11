package db

type WorkflowScripts struct {
	Id         int64 `gorm:"primary_key"`
	Workflow   *Workflow
	WorkflowId int64
	Script     *Script
	ScriptId   int64
	Weight     int64
}

func (d *WorkflowScripts) TableName() string {
	return "workflow_scripts"
}
