// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Entity ...
type Entity struct{}

// NewEntityDto ...
func NewEntityDto() Entity {
	return Entity{}
}

// AddEntity ...
func (r Entity) AddEntity(obj *api.NewEntityRequest) (entity *m.Entity) {
	entity = &m.Entity{
		Description: obj.Description,
		PluginName:  obj.PluginName,
		Actions:     make([]*m.EntityAction, 0, len(entity.Actions)),
		States:      make([]*m.EntityState, 0, len(entity.States)),
		Scripts:     make([]*m.Script, 0, len(entity.Scripts)),
		Hidden:      obj.Hidden,
		AutoLoad:    obj.AutoLoad,
		Icon:        obj.Icon,
		Attributes:  AttributeFromApi(obj.Attributes),
		Settings:    AttributeFromApi(obj.Settings),
	}

	// image
	if obj.Image != nil && obj.Image.Id != 0 {
		entity.Image = &m.Image{
			Id: obj.Image.Id,
		}
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.EntityAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
			Type:        a.Type,
		}
		if a.Image != nil {
			action.ImageId = common.Int64(a.Image.Id)
		}
		if a.Script != nil {
			action.ScriptId = common.Int64(a.Script.Id)
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
		}
		if s.Image != nil {
			state.ImageId = common.Int64(s.Image.Id)
		}

		entity.States = append(entity.States, state)
	}
	// area
	if obj.Area != nil && obj.Area.Id != 0 {
		entity.AreaId = &obj.Area.Id
	}
	// scripts
	for _, s := range obj.Scripts {
		script := &m.Script{
			Id: s.Id,
		}
		entity.Scripts = append(entity.Scripts, script)
	}
	// parent
	if obj.Parent != nil {
		parentId := common.EntityId(obj.Parent.Id)
		entity.ParentId = &parentId
	}

	return
}

// UpdateEntity ...
func (r Entity) UpdateEntity(obj *api.UpdateEntityRequest) (entity *m.Entity) {
	entity = &m.Entity{
		Id:          common.EntityId(obj.Id),
		Description: obj.Description,
		PluginName:  obj.PluginName,
		Actions:     make([]*m.EntityAction, 0, len(entity.Actions)),
		States:      make([]*m.EntityState, 0, len(entity.States)),
		Scripts:     make([]*m.Script, 0, len(entity.Scripts)),
		Hidden:      obj.Hidden,
		AutoLoad:    obj.AutoLoad,
		Icon:        obj.Icon,
		Attributes:  AttributeFromApi(obj.Attributes),
		Settings:    AttributeFromApi(obj.Settings),
	}

	// image
	if obj.Image != nil && obj.Image.Id != 0 {
		entity.Image = &m.Image{
			Id: obj.Image.Id,
		}
	}
	// actions
	for _, a := range obj.Actions {
		action := &m.EntityAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
			Type:        a.Type,
		}
		if a.Image != nil {
			action.ImageId = common.Int64(a.Image.Id)
		}
		if a.Script != nil {
			action.ScriptId = common.Int64(a.Script.Id)
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
		}
		if s.Image != nil {
			state.ImageId = common.Int64(s.Image.Id)
		}

		entity.States = append(entity.States, state)
	}
	// area
	if obj.Area != nil && obj.Area.Id != 0 {
		entity.AreaId = &obj.Area.Id
	}
	// scripts
	for _, s := range obj.Scripts {
		script := &m.Script{
			Id: s.Id,
		}
		entity.Scripts = append(entity.Scripts, script)
	}
	// parent
	if obj.Parent != nil {
		parentId := common.EntityId(obj.Parent.Id)
		entity.ParentId = &parentId
	}
	return
}

// ToSearchResult ...
func (r Entity) ToSearchResult(list []*m.Entity) *api.SearchEntityResult {

	items := make([]*api.Entity, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToEntity(i))
	}

	return &api.SearchEntityResult{
		Items: items,
	}
}

// ToListResult ...
func (r Entity) ToListResult(list []*m.Entity, total, limit, offset uint64) *api.GetEntityListResult {

	items := make([]*api.Entity, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToEntity(i))
	}

	return &api.GetEntityListResult{
		Items: items,
		Meta: &api.GetEntityListResult_Meta{
			Limit:        limit,
			ObjectsCount: total,
			Offset:       offset,
		},
	}
}

// ToEntity ...
func (r Entity) ToEntity(entity *m.Entity) (obj *api.Entity) {
	imageDto := NewImageDto()
	scriptDto := NewScriptDto()
	obj = &api.Entity{
		Id:          entity.Id.String(),
		PluginName:  entity.PluginName,
		Description: entity.Description,
		Icon:        entity.Icon,
		Hidden:      entity.Hidden,
		AutoLoad:    entity.AutoLoad,
		Actions:     make([]*api.EntityAction, 0, len(entity.Actions)),
		States:      make([]*api.EntityState, 0, len(entity.States)),
		Scripts:     make([]*api.Script, 0, len(entity.Scripts)),
		CreatedAt:   timestamppb.New(entity.CreatedAt),
		UpdatedAt:   timestamppb.New(entity.UpdatedAt),
	}
	// area
	if entity.Area != nil {
		obj.Area = &api.Area{
			Id:          entity.Area.Id,
			Name:        entity.Area.Name,
			Description: entity.Area.Description,
		}
	}
	// image
	if entity.Image != nil {
		obj.Image = imageDto.ToImage(entity.Image)
	}
	// parent
	if entity.ParentId != nil {
		obj.Parent = &api.EntityParent{
			Id: entity.ParentId.String(),
		}
	}
	// actions
	for _, a := range entity.Actions {
		action := &api.EntityAction{
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
			action.Script = scriptDto.ToGScript(a.Script)
		}
		obj.Actions = append(obj.Actions, action)
	}
	// states
	for _, s := range entity.States {
		state := &api.EntityState{
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
		script := scriptDto.ToGScript(s)
		obj.Scripts = append(obj.Scripts, script)
	}
	return
}
