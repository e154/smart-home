package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type WorkflowScenarioScripts struct {
	Db *gorm.DB
}

type WorkflowScenarioScript struct {
	Id                 int64 `gorm:"primary_key"`
	ScriptId           int64
	WorkflowScenarioId int64
}

func (d *WorkflowScenarioScript) TableName() string {
	return "workflow_scenario_scripts"
}

func (n WorkflowScenarioScripts) Add(scenario *WorkflowScenarioScript) (id int64, err error) {
	if err = n.Db.Create(&scenario).Error; err != nil {
		return
	}
	id = scenario.Id
	return
}

func (n WorkflowScenarioScripts) Delete(workflowId int64) (err error) {
	err = n.Db.Delete(&WorkflowScenarioScript{Id: workflowId}).Error
	return
}

func (n *WorkflowScenarioScripts) List(limit, offset int64, orderBy, sort string) (list []*WorkflowScenarioScript, total int64, err error) {

	if err = n.Db.Model(WorkflowScenarioScript{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*WorkflowScenarioScript, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
