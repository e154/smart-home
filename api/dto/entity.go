// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package dto

import (
	"fmt"

	"github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Entity ...
type Entity struct{}

// NewEntityDto ...
func NewEntityDto() Entity {
	return Entity{}
}

// AddEntity ...
func (r Entity) AddEntity(obj *stub.ApiNewEntityRequest) (entity *m.Entity) {

	entity = &m.Entity{
		Id:          common.EntityId(fmt.Sprintf("%s.%s", obj.PluginName, obj.Name)),
		Description: obj.Description,
		PluginName:  obj.PluginName,
		AutoLoad:    obj.AutoLoad,
		Icon:        obj.Icon,
		Attributes:  AttributeFromApi(obj.Attributes),
		Settings:    AttributeFromApi(obj.Settings),
		Metrics:     AddMetric(obj.Metrics),
		ImageId:     obj.ImageId,
		AreaId:      obj.AreaId,
	}

	// actions
	for _, a := range obj.Actions {
		action := &m.EntityAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
			Type:        a.Type,
			ImageId:     a.ImageId,
			ScriptId:    a.ScriptId,
		}

		entity.Actions = append(entity.Actions, action)
	}
	// states
	for _, s := range obj.States {
		state := &m.EntityState{
			Name:        s.Name,
			Description: s.Description,
			Icon:        s.Icon,
			Style:       s.Style,
			ImageId:     s.ImageId,
		}

		entity.States = append(entity.States, state)
	}

	// scripts
	for _, id := range obj.ScriptIds {
		script := &m.Script{
			Id: id,
		}
		entity.Scripts = append(entity.Scripts, script)
	}
	// parent
	if obj.ParentId != nil {
		parentId := common.EntityId(*obj.ParentId)
		entity.ParentId = &parentId
	}

	return
}

// UpdateEntity ...
func (r Entity) UpdateEntity(obj *stub.EntityServiceUpdateEntityJSONBody) (entity *m.Entity) {
	entity = &m.Entity{
		Id:          common.EntityId(obj.Id),
		Description: obj.Description,
		PluginName:  obj.PluginName,
		Actions:     make([]*m.EntityAction, 0),
		States:      make([]*m.EntityState, 0),
		Scripts:     make([]*m.Script, 0),
		AutoLoad:    obj.AutoLoad,
		Icon:        obj.Icon,
		Attributes:  AttributeFromApi(obj.Attributes),
		Settings:    AttributeFromApi(obj.Settings),
		ParentId:    nil,
		Metrics:     AddMetric(obj.Metrics),
		ImageId:     obj.ImageId,
		AreaId:      obj.AreaId,
	}

	// actions
	for _, a := range obj.Actions {
		action := &m.EntityAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
			Type:        a.Type,
			ImageId:     a.ImageId,
			ScriptId:    a.ScriptId,
		}

		entity.Actions = append(entity.Actions, action)
	}
	// states
	for _, s := range obj.States {
		state := &m.EntityState{
			Name:        s.Name,
			Description: s.Description,
			Icon:        s.Icon,
			Style:       s.Style,
			ImageId:     s.ImageId,
		}

		entity.States = append(entity.States, state)
	}

	// scripts
	for _, id := range obj.ScriptIds {
		script := &m.Script{
			Id: id,
		}
		entity.Scripts = append(entity.Scripts, script)
	}
	// parent
	if obj.ParentId != nil {
		parentId := common.EntityId(*obj.ParentId)
		entity.ParentId = &parentId
	}
	return
}

// ToSearchResult ...
func (r Entity) ToSearchResult(list []*m.Entity) *stub.ApiSearchEntityResult {

	items := make([]stub.ApiEntityShort, 0, len(list))

	for _, i := range list {
		items = append(items, *r.ToEntityShort(i))
	}

	return &stub.ApiSearchEntityResult{
		Items: items,
	}
}

// ToListResult ...
func (r Entity) ToListResult(list []*m.Entity) []*stub.ApiEntityShort {

	items := make([]*stub.ApiEntityShort, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToEntityShort(i))
	}

	return items
}

// ToEntityShort ...
func (r Entity) ToEntityShort(entity *m.Entity) (obj *stub.ApiEntityShort) {
	obj = &stub.ApiEntityShort{
		Id:          entity.Id.String(),
		PluginName:  entity.PluginName,
		Description: entity.Description,
		Icon:        entity.Icon,
		AutoLoad:    entity.AutoLoad,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		ParentId:    entity.ParentId.StringPtr(),
		IsLoaded:    common.Bool(entity.IsLoaded),
	}
	// area
	if entity.Area != nil {
		obj.Area = &stub.ApiArea{
			Id:          entity.Area.Id,
			Name:        entity.Area.Name,
			Description: entity.Area.Description,
		}
	}

	return
}

// ToEntity ...
func ToEntity(entity *m.Entity) (obj *stub.ApiEntity) {
	if entity == nil {
		return
	}
	imageDto := NewImageDto()
	scriptDto := NewScriptDto()
	obj = &stub.ApiEntity{
		Id:          entity.Id.String(),
		PluginName:  entity.PluginName,
		Description: entity.Description,
		Icon:        entity.Icon,
		AutoLoad:    entity.AutoLoad,
		IsLoaded:    common.Bool(entity.IsLoaded),
		Actions:     make([]stub.ApiEntityAction, 0, len(entity.Actions)),
		States:      make([]stub.ApiEntityState, 0, len(entity.States)),
		Scripts:     make([]stub.ApiScript, 0, len(entity.Scripts)),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Attributes:  AttributeToApi(entity.Attributes),
		Settings:    AttributeToApi(entity.Settings),
		Metrics:     Metrics(entity.Metrics),
	}
	// area
	if entity.Area != nil {
		obj.Area = &stub.ApiArea{
			Id:          entity.Area.Id,
			Name:        entity.Area.Name,
			Description: entity.Area.Description,
		}
	}
	// image
	if entity.Image != nil {
		obj.Image = imageDto.ToImageShort(entity.Image)
	}
	// parent
	if entity.ParentId != nil {
		obj.Parent = &stub.ApiEntityParent{
			Id: entity.ParentId.String(),
		}
	}
	// actions
	for _, a := range entity.Actions {
		action := stub.ApiEntityAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
			Type:        a.Type,
		}
		// images
		if a.Image != nil {
			action.Image = imageDto.ToImage(a.Image)
		}
		// script
		if a.Script != nil {
			action.Script = scriptDto.GetStubScriptShort(a.Script)
		}
		obj.Actions = append(obj.Actions, action)
	}
	// states
	for _, s := range entity.States {
		state := stub.ApiEntityState{
			Name:        s.Name,
			Description: s.Description,
			Icon:        s.Icon,
			Image:       nil,
			Style:       s.Style,
		}
		// images
		if s.Image != nil {
			state.Image = imageDto.ToImage(s.Image)
		}
		obj.States = append(obj.States, state)
	}
	// scripts
	for _, s := range entity.Scripts {
		script := scriptDto.GetStubScriptShort(s)
		obj.Scripts = append(obj.Scripts, *script)
		obj.ScriptIds = append(obj.ScriptIds, s.Id)
	}
	return
}

func (r Entity) ImportEntity(from *stub.ApiEntity) (to *m.Entity) {

	areaId, area := ImportArea(from.Area)
	var parentId *common.EntityId
	if from.Parent != nil {
		parentId = common.NewEntityId(from.Parent.Id)
	}

	_, image := ImportImage(from.Image)
	to = &m.Entity{
		Id:          common.EntityId(from.Id),
		Description: from.Description,
		PluginName:  from.PluginName,
		Icon:        from.Icon,
		Image:       image,
		Actions:     make([]*m.EntityAction, 0, len(from.Actions)),
		States:      make([]*m.EntityState, 0, len(from.States)),
		Area:        area,
		AreaId:      areaId,
		Metrics:     AddMetric(from.Metrics),
		Scripts:     make([]*m.Script, 0, len(from.Scripts)),
		Attributes:  AttributeFromApi(from.Attributes),
		Settings:    AttributeFromApi(from.Settings),
		AutoLoad:    from.AutoLoad,
		ParentId:    parentId,
	}

	// ACTIONS
	for _, action := range from.Actions {
		imageId, img := ImportImage(action.Image)
		scriptId, script := ImportScript(action.Script)
		to.Actions = append(to.Actions, &m.EntityAction{
			Name:        action.Name,
			Description: action.Description,
			Icon:        action.Icon,
			Image:       img,
			ImageId:     imageId,
			Script:      script,
			ScriptId:    scriptId,
			Type:        action.Type,
		})
	}

	// STATES
	for _, state := range from.States {
		imageId, img := ImportImage(state.Image)
		to.States = append(to.States, &m.EntityState{
			Name:        state.Name,
			Description: state.Description,
			Icon:        state.Icon,
			Image:       img,
			ImageId:     imageId,
			Style:       state.Style,
		})
	}

	// SCRIPTS
	for _, script := range from.Scripts {
		_, _script := ImportScript(&script)
		to.Scripts = append(to.Scripts, _script)
	}

	return
}
