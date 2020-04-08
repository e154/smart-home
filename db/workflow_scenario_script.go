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
	"fmt"
	"github.com/jinzhu/gorm"
)

// WorkflowScenarioScripts ...
type WorkflowScenarioScripts struct {
	Db *gorm.DB
}

// WorkflowScenarioScript ...
type WorkflowScenarioScript struct {
	Id                 int64 `gorm:"primary_key"`
	ScriptId           int64
	WorkflowScenarioId int64
}

// TableName ...
func (d *WorkflowScenarioScript) TableName() string {
	return "workflow_scenario_scripts"
}

// Add ...
func (n WorkflowScenarioScripts) Add(scenario *WorkflowScenarioScript) (id int64, err error) {
	if err = n.Db.Create(&scenario).Error; err != nil {
		return
	}
	id = scenario.Id
	return
}

// Delete ...
func (n WorkflowScenarioScripts) Delete(workflowId int64) (err error) {
	err = n.Db.Delete(&WorkflowScenarioScript{Id: workflowId}).Error
	return
}

// List ...
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
