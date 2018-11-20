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

func (n Workers) Delete(workerId int64) (err error) {
	err = n.Db.Delete(&Worker{Id: workerId}).Error
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
