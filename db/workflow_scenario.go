package db

import (
	"time"
	"fmt"
	"github.com/jinzhu/gorm"
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

func (n WorkflowScenarios) GetById(workflowId int64) (scenario *WorkflowScenario, err error) {
	scenario = &WorkflowScenario{Id: workflowId}
	err = n.Db.First(&scenario).Error
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

func (n *WorkflowScenarios) AddScript(workflowScenarioId, scriptId int64) (err error) {
	err = n.Db.Create(&WorkflowScenarioScript{WorkflowScenarioId: workflowScenarioId, ScriptId: scriptId}).Error
	return
}

func (n *WorkflowScenarios) RemoveScript(workflowScenarioId, scriptId int64) (err error) {
	err = n.Db.Delete(&WorkflowScenarioScript{WorkflowScenarioId: workflowScenarioId, ScriptId: scriptId}).Error
	return
}
