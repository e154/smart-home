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

package adaptors

import (
	"encoding/json"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"github.com/jinzhu/gorm"
)

// FlowElement ...
type FlowElement struct {
	table *db.FlowElements
	db    *gorm.DB
}

// GetFlowElementAdaptor ...
func GetFlowElementAdaptor(d *gorm.DB) *FlowElement {
	return &FlowElement{
		table: &db.FlowElements{Db: d},
		db:    d,
	}
}

// Add ...
func (n *FlowElement) Add(element *m.FlowElement) (id uuid.UUID, err error) {
	dbFlowElement := n.toDb(element)
	id, err = n.table.Add(dbFlowElement)
	return
}

// GetAllEnabled ...
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

// GetById ...
func (n *FlowElement) GetById(elementId uuid.UUID) (element *m.FlowElement, err error) {

	var dbFlowElement *db.FlowElement
	if dbFlowElement, err = n.table.GetById(elementId); err != nil {
		return
	}

	element = n.fromDb(dbFlowElement)

	return
}

// Update ...
func (n *FlowElement) Update(element *m.FlowElement) (err error) {
	dbFlowElement := n.toDb(element)
	err = n.table.Update(dbFlowElement)
	return
}

// AddOrUpdateFlowElement ...
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

// Delete ...
func (n *FlowElement) Delete(ids []uuid.UUID) (err error) {
	err = n.table.Delete(ids)
	return
}

// List ...
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
