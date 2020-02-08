// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Workflows struct {
	Db *gorm.DB
}

type Workflow struct {
	Id                 int64 `gorm:"primary_key"`
	Name               string
	Description        string
	Status             string
	WorkflowScenarioId *int64
	WorkflowScenario   *WorkflowScenario
	Scenarios          []*WorkflowScenario
	Scripts            []*Script
	CreatedAt          time.Time
	UpdatedAt          time.Time
	//Flows		[]*Flow			`orm:"-" json:"flows"`
}

func (d *Workflow) TableName() string {
	return "workflows"
}

func (n Workflows) Add(workflow *Workflow) (id int64, err error) {
	if err = n.Db.Create(&workflow).Error; err != nil {
		return
	}
	id = workflow.Id
	return
}

func (n Workflows) GetAllEnabled() (list []*Workflow, err error) {
	list = make([]*Workflow, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	if err != nil {
		return
	}

	for _, item := range list {
		if err = n.DependencyLoading(item); err != nil {
			return
		}
	}

	return
}

func (n Workflows) GetById(workflowId int64) (workflow *Workflow, err error) {
	workflow = &Workflow{Id: workflowId}
	if err = n.Db.First(&workflow).Error; err != nil {
		return
	}
	err = n.DependencyLoading(workflow)
	return
}

func (n Workflows) GetByWorkflowScenarioId(workflowScenarioId int64) (workflow *Workflow, err error) {
	result := make([]*Workflow, 0)
	if err = n.Db.Raw(`select w.*
		from workflow_scenarios ws
	left join workflows w on ws.workflow_id = w.id
	where ws.id = ?`, workflowScenarioId).Scan(&result).Error; err != nil {
		return
	}
	if len(result) == 0 {
		err = errors.New("not found")
		return
	}
	workflow = result[0]
	err = n.DependencyLoading(workflow)
	return
}

func (n Workflows) Update(m *Workflow) (err error) {
	updateParams := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
	}
	if m.WorkflowScenarioId != nil {
		updateParams["workflow_scenario_id"] = m.WorkflowScenarioId
	}
	err = n.Db.Model(&Workflow{Id: m.Id}).Updates(updateParams).Error
	return
}

func (n Workflows) Delete(workflowId int64) (err error) {
	err = n.Db.Delete(&Workflow{Id: workflowId}).Error
	return
}

func (n *Workflows) List(limit, offset int64, orderBy, sort string, onlyEnabled bool) (list []*Workflow, total int64, err error) {

	if err = n.Db.Model(Workflow{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Workflow, 0)
	q := n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy))

	if onlyEnabled {
		q = q.Where("status = ?", "enabled")
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		return
	}

	for _, item := range list {
		if err = n.DependencyLoading(item); err != nil {
			return
		}
	}

	return
}

func (n *Workflows) DependencyLoading(workflow *Workflow) (err error) {

	workflow.Scenarios = make([]*WorkflowScenario, 0)
	workflow.Scripts = make([]*Script, 0)

	//fmt.Println("----")
	//debug.Println(workflow)
	//fmt.Println("----")

	// load scripts
	var rows1, rows3 *sql.Rows
	rows1, err = n.Db.Model(&WorkflowScripts{}).
		Where("workflow_scripts.workflow_id = ?", workflow.Id).
		Joins("left join scripts s on workflow_scripts.script_id = s.id").
		Select("s.id, s.lang, s.name, s.source, s.description, s.compiled, s.created_at, s.updated_at").
		Rows()
	if err != nil {
		return
	}
	defer rows1.Close()

	for rows1.Next() {
		s := &Script{}
		rows1.Scan(&s.Id, &s.Lang, &s.Name, &s.Source, &s.Description, &s.Compiled, &s.CreatedAt, &s.UpdatedAt)
		workflow.Scripts = append(workflow.Scripts, s)
	}

	// load scenarios
	err = n.Db.Model(&WorkflowScenario{}).
		Where("workflow_scenarios.workflow_id = ?", workflow.Id).
		Find(&workflow.Scenarios).Error
	if err != nil {
		return
	}

	// load scenarios scripts
	for _, scenario := range workflow.Scenarios {
		scenario.Scripts = make([]*Script, 0)
		rows3, err = n.Db.Model(&WorkflowScenarioScript{}).
			Where("workflow_scenario_scripts.workflow_scenario_id = ?", scenario.Id).
			Joins("left join scripts s on workflow_scenario_scripts.script_id = s.id").
			Select("s.id, s.lang, s.name, s.source, s.description, s.compiled, s.created_at, s.updated_at").
			Rows()
		if err != nil {
			return
		}
		defer rows3.Close()

		for rows3.Next() {
			s := &Script{}
			rows3.Scan(&s.Id, &s.Lang, &s.Name, &s.Source, &s.Description, &s.Compiled, &s.CreatedAt, &s.UpdatedAt)
			scenario.Scripts = append(scenario.Scripts, s)
		}
	}

	return
}

func (n *Workflows) AddScript(workflowId, scriptId int64) (err error) {
	err = n.Db.Create(&WorkflowScripts{
		WorkflowId: workflowId,
		ScriptId:   scriptId,
		Weight:     0,
	}).Error
	return
}

func (n *Workflows) RemoveScript(workflowId, scriptId int64) (err error) {
	err = n.Db.Delete(&WorkflowScripts{}, "workflow_id = ? and script_id = ?", workflowId, scriptId).Error
	return
}

func (n *Workflows) SetScenario(workflowId int64, scenarioId *int64) (err error) {
	err = n.Db.Model(&Workflow{Id: workflowId}).Updates(map[string]interface{}{
		"workflow_scenario_id": scenarioId,
	}).Error
	return
}

func (n *Workflows) Search(query string, limit, offset int) (list []*Workflow, total int64, err error) {

	q := n.Db.Model(&Workflow{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Workflow, 0)
	err = q.Find(&list).Error

	return
}
