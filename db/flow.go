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
	. "github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

type Flows struct {
	Db *gorm.DB
}

type Flow struct {
	Id                 int64 `gorm:"primary_key"`
	Name               string
	Description        string
	Status             StatusType
	Workflow           *Workflow
	WorkflowId         int64
	WorkflowScenarioId int64
	Connections        []*Connection
	FlowElements       []*FlowElement
	Workers            []*Worker
	Subscriptions      []*FlowSubscription
	Zigbee2mqttDevices []*Zigbee2mqttDevice
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (d *Flow) TableName() string {
	return "flows"
}

func (n Flows) Add(flow *Flow) (id int64, err error) {
	if err = n.Db.Create(&flow).Error; err != nil {
		return
	}
	id = flow.Id

	err = n.DependencyLoading(flow)
	return
}

func (n Flows) GetAllEnabled() (list []*Flow, err error) {
	list = make([]*Flow, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	if err != nil {
		return
	}

	for _, flow := range list {
		if err = n.DependencyLoading(flow); err != nil {
			return
		}
	}
	return
}

func (n Flows) GetAllEnabledByWorkflow(workflowId int64) (list []*Flow, err error) {
	list = make([]*Flow, 0)
	err = n.Db.
		Joins("left join workflows w on w.id = ?", workflowId).
		Where("flows.status = 'enabled' and workflow_id = ?", workflowId).
		Where("flows.workflow_scenario_id = w.workflow_scenario_id").
		Find(&list).Error
	if err != nil {
		return
	}

	for _, flow := range list {
		if err = n.DependencyLoading(flow); err != nil {
			return
		}
	}
	return
}

func (n Flows) GetById(flowId int64) (flow *Flow, err error) {
	flow = &Flow{Id: flowId}
	if err = n.Db.First(&flow).Error; err != nil {
		return
	}

	err = n.DependencyLoading(flow)

	return
}

func (n Flows) Update(m *Flow) (err error) {
	err = n.Db.Model(&Flow{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"workflow_id": m.WorkflowId,
		"scenario_id": m.WorkflowScenarioId,
	}).Error
	return
}

func (n Flows) Delete(flowId int64) (err error) {
	err = n.Db.Delete(&Flow{Id: flowId}).Error
	return
}

func (n *Flows) List(limit, offset int64, orderBy, sort string) (list []*Flow, total int64, err error) {

	if err = n.Db.Model(Flow{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Flow, 0)
	q := n.Db.
		Limit(limit).
		Offset(offset)

	if orderBy != "" && sort != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.Find(&list).Error
	if err != nil {
		return
	}

	for _, flow := range list {
		if err = n.DependencyLoading(flow); err != nil {
			return
		}
	}
	return
}

func (n *Flows) Search(query string, limit, offset int) (list []*Flow, total int64, err error) {

	q := n.Db.Model(&Flow{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Flow, 0)
	err = q.Find(&list).Error

	return
}

func (n *Flows) DependencyLoading(flow *Flow) (err error) {
	flow.Connections = make([]*Connection, 0)
	flow.FlowElements = make([]*FlowElement, 0)
	flow.Workers = make([]*Worker, 0)
	flow.Subscriptions = make([]*FlowSubscription, 0)
	flow.Zigbee2mqttDevices = make([]*Zigbee2mqttDevice, 0)
	flow.Workflow = &Workflow{}

	n.Db.Model(flow).
		Related(&flow.Connections).
		Related(&flow.FlowElements).
		Related(&flow.Workflow).
		Related(&flow.Subscriptions)

	if flow.Workflow.WorkflowScenarioId != nil {
		flow.Workflow.WorkflowScenario = &WorkflowScenario{}
		n.Db.Model(flow).
			Related(flow.Workflow.WorkflowScenario)
	}

	// scripts
	var scriptIds []int64
	for _, element := range flow.FlowElements {
		if element.ScriptId != nil {
			scriptIds = append(scriptIds, *element.ScriptId)
		}
	}

	scripts := make([]*Script, 0)
	err = n.Db.Model(&Script{}).
		Where("id in (?)", scriptIds).
		Find(&scripts).
		Error
	if err != nil {
		return
	}

	for _, element := range flow.FlowElements {
		if element.ScriptId != nil {
			for _, script := range scripts {
				if *element.ScriptId == script.Id {
					element.Script = script
				}
			}
		}
	}

	// workers
	err = n.Db.Model(&Worker{}).
		Where("flow_id = ?", flow.Id).
		Preload("DeviceAction").
		Preload("DeviceAction.Script").
		Preload("DeviceAction.Device").
		Preload("DeviceAction.Device.Devices").
		Preload("DeviceAction.Device.Node").
		Preload("DeviceAction.Device.States").
		Preload("DeviceAction.Device.Actions").
		Preload("DeviceAction.Device.Actions.Script").
		Find(&flow.Workers).
		Error

	if err != nil {
		return
	}

	err = n.Db.Raw(`select *
from zigbee2mqtt_devices
where id in (
    select zigbee2mqtt_device_id
    from flow_zigbee2mqtt_devices
    where flow_id = ?
    )`, flow.Id).Scan(&flow.Zigbee2mqttDevices).
		Error

	return
}
