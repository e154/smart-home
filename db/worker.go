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
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Workers struct {
	Db *gorm.DB
}

type Worker struct {
	Id             int64 `gorm:"primary_key"`
	Workflow       *Workflow
	WorkflowId     int64
	DeviceAction   *DeviceAction
	DeviceActionId int64
	Flow           *Flow
	FlowId         int64
	Status         string
	Name           string
	Time           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *Worker) TableName() string {
	return "workers"
}

func (n Workers) Add(worker *Worker) (id int64, err error) {
	if err = n.Db.Create(&worker).Error; err != nil {
		return
	}
	id = worker.Id
	return
}

func (n Workers) GetAllEnabled() (list []*Worker, err error) {
	list = make([]*Worker, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	return
}

func (n Workers) GetById(workerId int64) (worker *Worker, err error) {
	worker = &Worker{Id: workerId}
	err = n.Db.First(&worker).Error
	return
}

func (n Workers) Update(m *Worker) (err error) {
	err = n.Db.Model(&Worker{Id: m.Id}).Updates(map[string]interface{}{
		"name":             m.Name,
		"status":           m.Status,
		"workflow_id":      m.WorkflowId,
		"flow_id":          m.FlowId,
		"time":             m.Time,
		"device_action_id": m.DeviceActionId,
	}).Error
	return
}

func (n Workers) Delete(ids []int64) (err error) {
	err = n.Db.Delete(&Worker{}, "id in (?)", ids).Error
	return
}

func (n *Workers) List(limit, offset int64, orderBy, sort string) (list []*Worker, total int64, err error) {

	if err = n.Db.Model(Worker{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Worker, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
