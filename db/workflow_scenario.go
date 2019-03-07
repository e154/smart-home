package db

import (
	"time"
	"fmt"
	"github.com/jinzhu/gorm"
	"database/sql"
)

type WorkflowScenarios struct {
	Db *gorm.DB
}

type WorkflowScenario struct {
	Id         int64 `gorm:"primary_key"`
	Name       string
	SystemName string
	WorkflowId int64
	Scripts    []*Script
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (d *WorkflowScenario) TableName() string {
	return "workflow_scenarios"
}

func (n WorkflowScenarios) Add(scenario *WorkflowScenario) (id int64, err error) {
	if err = n.Db.Create(&scenario).Error; err != nil {
		return
	}
	id = scenario.Id
	return
}

func (n WorkflowScenarios) GetById(scenarioId int64) (scenario *WorkflowScenario, err error) {
	scenario = &WorkflowScenario{Id: scenarioId}
	if err = n.Db.First(&scenario).Error; err != nil {
		return
	}

	err = n.DependencyLoading(scenario)
	return
}

func (n WorkflowScenarios) Update(m *WorkflowScenario) (err error) {
	err = n.Db.Model(&WorkflowScenario{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"system_name": m.SystemName,
	}).Error
	return
}

func (n WorkflowScenarios) Delete(workflowId int64) (err error) {
	err = n.Db.Delete(&WorkflowScenario{Id: workflowId}).Error
	return
}

func (n *WorkflowScenarios) List(limit, offset int64, orderBy, sort string) (list []*WorkflowScenario, total int64, err error) {

	if err = n.Db.Model(WorkflowScenario{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*WorkflowScenario, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *WorkflowScenarios) ListByWorkflow(workflowId int64) (list []*WorkflowScenario, total int64, err error) {

	if err = n.Db.Model(WorkflowScenario{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*WorkflowScenario, 0)
	err = n.Db.
		Where("workflow_id = ?", workflowId).
		Find(&list).
		Error

	return
}

func (n *WorkflowScenarios) AddScript(workflowScenarioId, scriptId int64) (err error) {
	err = n.Db.Create(&WorkflowScenarioScript{WorkflowScenarioId: workflowScenarioId, ScriptId: scriptId}).Error
	return
}

func (n *WorkflowScenarios) RemoveScript(workflowScenarioId, scriptId int64) (err error) {
	err = n.Db.Delete(&WorkflowScenarioScript{}, "workflow_scenario_id = ? and script_id = ?", workflowScenarioId, scriptId).Error
	return
}

func (n *WorkflowScenarios) Search(query string, limit, offset int) (list []*WorkflowScenario, total int64, err error) {

	q := n.Db.Model(&WorkflowScenario{}).
		Where("name ILIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*WorkflowScenario, 0)
	err = q.Find(&list).Error

	return
}

func (n *WorkflowScenarios) DependencyLoading(scenario *WorkflowScenario) (err error) {

	var rows *sql.Rows

	scenario.Scripts = make([]*Script, 0)
	rows, err = n.Db.Model(&WorkflowScenarioScript{}).
		Where("workflow_scenario_scripts.workflow_scenario_id = ?", scenario.Id).
		Joins("left join scripts s on workflow_scenario_scripts.script_id = s.id").
		Select("s.id, s.lang, s.name, s.source, s.description, s.compiled, s.created_at, s.updated_at").
		Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		s := &Script{}
		rows.Scan(&s.Id, &s.Lang, &s.Name, &s.Source, &s.Description, &s.Compiled, &s.CreatedAt, &s.UpdatedAt)
		scenario.Scripts = append(scenario.Scripts, s)
	}

	return
}