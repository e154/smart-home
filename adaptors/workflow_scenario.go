package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type WorkflowScenario struct {
	table *db.WorkflowScenarios
	db    *gorm.DB
}

func GetWorkflowScenarioAdaptor(d *gorm.DB) *WorkflowScenario {
	return &WorkflowScenario{
		table: &db.WorkflowScenarios{Db: d},
		db:    d,
	}
}

func (n *WorkflowScenario) Add(workflow *m.WorkflowScenario) (id int64, err error) {

	dbWorkflowScenario := n.toDb(workflow)
	if id, err = n.table.Add(dbWorkflowScenario); err != nil {
		return
	}

	return
}

func (n *WorkflowScenario) GetById(workflowId int64) (workflow *m.WorkflowScenario, err error) {

	var dbWorkflowScenario *db.WorkflowScenario
	if dbWorkflowScenario, err = n.table.GetById(workflowId); err != nil {
		return
	}

	workflow = n.fromDb(dbWorkflowScenario)

	return
}

func (n *WorkflowScenario) Update(workflow *m.WorkflowScenario) (err error) {
	dbWorkflowScenario := n.toDb(workflow)
	err = n.table.Update(dbWorkflowScenario)
	return
}

func (n *WorkflowScenario) Delete(workflowId int64) (err error) {
	err = n.table.Delete(workflowId)
	return
}

func (n *WorkflowScenario) List(limit, offset int64, orderBy, sort string) (list []*m.WorkflowScenario, total int64, err error) {
	var dbList []*db.WorkflowScenario
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.WorkflowScenario, 0)
	for _, dbWorkflowScenario := range dbList {
		workflow := n.fromDb(dbWorkflowScenario)
		list = append(list, workflow)
	}

	return
}

func (n *WorkflowScenario) AddScript(workflowScenario *m.WorkflowScenario, script *m.Script) (err error) {
	err = n.table.AddScript(workflowScenario.Id, script.Id)
	return
}

func (n *WorkflowScenario) RemoveScript(workflowScenario *m.WorkflowScenario, script *m.Script) (err error) {
	err = n.table.RemoveScript(workflowScenario.Id, script.Id)
	return
}

func (n *WorkflowScenario) fromDb(dbWorkflowScenario *db.WorkflowScenario) (workflow *m.WorkflowScenario) {
	workflow = &m.WorkflowScenario{
		Id:         dbWorkflowScenario.Id,
		Name:       dbWorkflowScenario.Name,
		WorkflowId: dbWorkflowScenario.WorkflowId,
		SystemName: dbWorkflowScenario.SystemName,
		CreatedAt:  dbWorkflowScenario.CreatedAt,
		UpdatedAt:  dbWorkflowScenario.UpdatedAt,
	}

	// scripts
	workflow.Scripts = make([]*m.Script, 0)
	scriptAdaptor := GetScriptAdaptor(n.db)
	for _, dbScript := range dbWorkflowScenario.Scripts {
		script, _ := scriptAdaptor.fromDb(dbScript)
		workflow.Scripts = append(workflow.Scripts, script)
	}

	return
}

func (n *WorkflowScenario) toDb(workflow *m.WorkflowScenario) (dbWorkflowScenario *db.WorkflowScenario) {
	dbWorkflowScenario = &db.WorkflowScenario{
		Id:         workflow.Id,
		Name:       workflow.Name,
		WorkflowId: workflow.WorkflowId,
		SystemName: workflow.SystemName,
		CreatedAt:  workflow.CreatedAt,
		UpdatedAt:  workflow.UpdatedAt,
	}
	return
}
