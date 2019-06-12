package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"encoding/json"
)

type FlowElement struct {
	table *db.FlowElements
	db    *gorm.DB
}

func GetFlowElementAdaptor(d *gorm.DB) *FlowElement {
	return &FlowElement{
		table: &db.FlowElements{Db: d},
		db:    d,
	}
}

func (n *FlowElement) Add(element *m.FlowElement) (id uuid.UUID, err error) {
	dbFlowElement := n.toDb(element)
	id, err = n.table.Add(dbFlowElement)
	return
}

func (n *FlowElement) GetAllEnabled() (list []*m.FlowElement, err error) {

	var dbList []*db.FlowElement
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.FlowElement, 0)
	for _, dbFlowElement := range dbList {
		element := n.fromDb(dbFlowElement)
		list = append(list, element)
	}

	return
}

func (n *FlowElement) GetById(elementId uuid.UUID) (element *m.FlowElement, err error) {

	var dbFlowElement *db.FlowElement
	if dbFlowElement, err = n.table.GetById(elementId); err != nil {
		return
	}

	element = n.fromDb(dbFlowElement)

	return
}

func (n *FlowElement) Update(element *m.FlowElement) (err error) {
	dbFlowElement := n.toDb(element)
	err = n.table.Update(dbFlowElement)
	return
}

func (n *FlowElement) AddOrUpdateFlowElement(element *m.FlowElement) (err error) {

	if element.Uuid.String() == "00000000-0000-0000-0000-000000000000" {
		_, err = n.Add(element)
		return
	}

	if _, err = n.table.GetById(element.Uuid); err != nil {
		_, err = n.Add(element)
		return
	}

	err = n.Update(element)

	return
}

func (n *FlowElement) Delete(ids []uuid.UUID) (err error) {
	err = n.table.Delete(ids)
	return
}

func (n *FlowElement) List(limit, offset int64, orderBy, sort string) (list []*m.FlowElement, total int64, err error) {
	var dbList []*db.FlowElement
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.FlowElement, 0)
	for _, dbFlowElement := range dbList {
		element := n.fromDb(dbFlowElement)
		list = append(list, element)
	}

	return
}

func (n *FlowElement) fromDb(dbFlowElement *db.FlowElement) (element *m.FlowElement) {
	element = &m.FlowElement{
		Uuid:          dbFlowElement.Uuid,
		Name:          dbFlowElement.Name,
		Status:        dbFlowElement.Status,
		Description:   dbFlowElement.Description,
		FlowLink:      dbFlowElement.FlowLink,
		ScriptId:      dbFlowElement.ScriptId,
		FlowId:        dbFlowElement.FlowId,
		PrototypeType: dbFlowElement.PrototypeType,
		CreatedAt:     dbFlowElement.CreatedAt,
		UpdatedAt:     dbFlowElement.UpdatedAt,
	}

	scriptAdaptor := GetScriptAdaptor(n.db)
	if dbFlowElement.Script != nil {
		element.Script, _ = scriptAdaptor.fromDb(dbFlowElement.Script)
	}

	graphSettings, _ := dbFlowElement.GraphSettings.MarshalJSON()
	if err := json.Unmarshal(graphSettings, &element.GraphSettings); err != nil {
		log.Error(err.Error())
	}

	return
}

func (n *FlowElement) toDb(element *m.FlowElement) (dbFlowElement *db.FlowElement) {
	dbFlowElement = &db.FlowElement{
		Uuid:          element.Uuid,
		Name:          element.Name,
		PrototypeType: element.PrototypeType,
		ScriptId:      element.ScriptId,
		FlowLink:      element.FlowLink,
		FlowId:        element.FlowId,
		Status:        element.Status,
		Description:   element.Description,
	}

	graphSettings, _ := json.Marshal(element.GraphSettings)
	dbFlowElement.GraphSettings.UnmarshalJSON(graphSettings)

	return
}
