package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Workflow struct {
	table *db.Workflows
	db    *gorm.DB
}

func GetWorkflowAdaptor(d *gorm.DB) *Workflow {
	return &Workflow{
		table: &db.Workflows{Db: d},
		db:    d,
	}
}

func (n *Workflow) Add(workflow *m.Workflow) (id int64, err error) {

	dbWorkflow := n.toDb(workflow)
	if id, err = n.table.Add(dbWorkflow); err != nil {
		return
	}

	return
}

func (n *Workflow) GetAllEnabled() (list []*m.Workflow, err error) {

	var dbList []*db.Workflow
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Workflow, 0)
	for _, dbWorkflow := range dbList {
		workflow := n.fromDb(dbWorkflow)
		list = append(list, workflow)
	}

	return
}

func (n *Workflow) GetById(workflowId int64) (workflow *m.Workflow, err error) {

	var dbWorkflow *db.Workflow
	if dbWorkflow, err = n.table.GetById(workflowId); err != nil {
		return
	}

	workflow = n.fromDb(dbWorkflow)

	return
}

func (n *Workflow) Update(workflow *m.Workflow) (err error) {
	dbWorkflow := n.toDb(workflow)
	err = n.table.Update(dbWorkflow)
	return
}

func (n *Workflow) Delete(workflowId int64) (err error) {
	err = n.table.Delete(workflowId)
	return
}

func (n *Workflow) List(limit, offset int64, orderBy, sort string) (list []*m.Workflow, total int64, err error) {
	var dbList []*db.Workflow
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Workflow, 0)
	for _, dbWorkflow := range dbList {
		workflow := n.fromDb(dbWorkflow)
		list = append(list, workflow)
	}

	return
}

func (n *Workflow) DependencyLoading(workflow *m.Workflow) (err error) {
	dbWorkflow := n.toDb(workflow)
	err = n.table.DependencyLoading(dbWorkflow)
	return
}

func (n *Workflow) fromDb(dbWorkflow *db.Workflow) (workflow *m.Workflow) {
	workflow = &m.Workflow{
		Id:          dbWorkflow.Id,
		Name:        dbWorkflow.Name,
		Description: dbWorkflow.Description,
		Status:      dbWorkflow.Status,
		CreatedAt:   dbWorkflow.CreatedAt,
		UpdatedAt:   dbWorkflow.UpdatedAt,
	}

	// scripts
	workflow.Scripts = make([]*m.Script, 0)
	scriptAdaptor := GetScriptAdaptor(n.db)
	for _, dbScript := range dbWorkflow.Scripts {
		script, _ := scriptAdaptor.fromDb(dbScript)
		workflow.Scripts = append(workflow.Scripts, script)
	}

	// scenario
	scenarioAdaptor := GetWorkflowScenarioAdaptor(n.db)
	if dbWorkflow.WorkflowScenarioId != nil {
		for _, dbScenario := range dbWorkflow.Scenarios {
			if dbScenario.Id == *dbWorkflow.WorkflowScenarioId {
				workflow.Scenario = scenarioAdaptor.fromDb(dbScenario)
				break
			}
		}
	}

	// scenarios
	workflow.Scenarios = make([]*m.WorkflowScenario, 0)
	for _, dbScenario := range dbWorkflow.Scenarios {
		scenario := scenarioAdaptor.fromDb(dbScenario)
		workflow.Scenarios = append(workflow.Scenarios, scenario)
	}

	return
}

func (n *Workflow) toDb(workflow *m.Workflow) (dbWorkflow *db.Workflow) {
	dbWorkflow = &db.Workflow{
		Id:          workflow.Id,
		Name:        workflow.Name,
		Description: workflow.Description,
		Status:      workflow.Status,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}
	return
}
