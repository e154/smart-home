package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/e154/smart-home/system/uuid"
	"encoding/json"
	. "github.com/e154/smart-home/common"
)

type FlowElements struct {
	Db *gorm.DB
}

type FlowElement struct {
	Uuid          uuid.UUID `gorm:"primary_key"`
	Name          string
	Description   string
	Flow          *Flow
	FlowId        int64
	Script        *Script
	ScriptId      *int64
	Status        StatusType
	FlowLink      *int64
	PrototypeType FlowElementsPrototypeType
	GraphSettings json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *FlowElement) TableName() string {
	return "flow_elements"
}

func (n FlowElements) Add(flow *FlowElement) (id uuid.UUID, err error) {
	if err = n.Db.Create(&flow).Error; err != nil {
		return
	}
	id = flow.Uuid
	return
}

func (n FlowElements) GetAllEnabled() (list []*FlowElement, err error) {
	list = make([]*FlowElement, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	return
}

func (n FlowElements) GetById(id uuid.UUID) (flow *FlowElement, err error) {
	flow = &FlowElement{Uuid: id}
	err = n.Db.First(&flow).Error
	return
}

func (n FlowElements) Update(m *FlowElement) (err error) {
	err = n.Db.Model(&FlowElement{Uuid: m.Uuid}).Updates(map[string]interface{}{
		"name":           m.Name,
		"description":    m.Description,
		"status":         m.Status,
		"flow_id":        m.FlowId,
		"script_id":      m.ScriptId,
		"flow_link":      m.FlowLink,
		"prototype_type": m.PrototypeType,
		"graph_settings": m.GraphSettings,
	}).Error

	return
}

func (n FlowElements) Delete(ids []uuid.UUID) (err error) {
	err = n.Db.Delete(&FlowElement{}, "uuid in ?", ids).Error
	return
}

func (n *FlowElements) List(limit, offset int64, orderBy, sort string) (list []*FlowElement, total int64, err error) {

	if err = n.Db.Model(FlowElement{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*FlowElement, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
