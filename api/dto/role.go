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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Role ...
type Role struct{}

// NewRoleDto ...
func NewRoleDto() Role {
	return Role{}
}

// FromNewRoleRequest ...
func (r Role) FromNewRoleRequest(from *api.NewRoleRequest) (to *m.Role) {
	to = &m.Role{
		Name:        from.Name,
		Description: from.Description,
	}
	if from.Parent != "" {
		to.Parent = &m.Role{
			Name: from.Parent,
		}
	}
	return
}

// FromUpdateRoleRequest ...
func (r Role) FromUpdateRoleRequest(from *api.UpdateRoleRequest) (to *m.Role) {
	to = &m.Role{
		Name:        from.Name,
		Description: from.Description,
	}
	if from.Parent != "" {
		to.Parent = &m.Role{
			Name: from.Parent,
		}
	}
	return
}

// ToSearchResult ...
func (r Role) ToSearchResult(list []*m.Role) *api.SearchRoleListResult {

	items := make([]*api.Role, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToGRole(i))
	}

	return &api.SearchRoleListResult{
		Items: items,
	}
}

// ToListResult ...
func (r Role) ToListResult(list []*m.Role, total uint64, pagination common.PageParams) *api.GetRoleListResult {

	items := make([]*api.Role, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToGRole(i))
	}

	return &api.GetRoleListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToGRole ...
func (r Role) ToGRole(from *m.Role) (to *api.Role) {
	to = &api.Role{
		Name:        from.Name,
		Description: from.Description,
		AccessList:  nil,
		CreatedAt:   timestamppb.New(from.CreatedAt),
		UpdatedAt:   timestamppb.New(from.UpdatedAt),
	}
	if from.Parent != nil {
		to.Parent = r.ToGRole(from.Parent)
	}
	if len(from.Children) > 0 {
		for _, ch := range from.Children {
			to.Children = append(to.Children, r.ToGRole(ch))
		}
	}
	if from.AccessList != nil {
		to.AccessList = &api.Role_AccessList{
			Levels: make(map[string]*api.Role_AccessList_ListOfString),
		}
		for levelName, levels := range from.AccessList {
			for _, item := range levels {
				if _, ok := to.AccessList.Levels[levelName]; !ok {
					to.AccessList.Levels[levelName] = &api.Role_AccessList_ListOfString{}
				}
				to.AccessList.Levels[levelName].Items = append(to.AccessList.Levels[levelName].Items, item)
			}
		}
	}
	return
}

// ToRoleAccessListResult ...
func (r Role) ToRoleAccessListResult(accessList access_list.AccessList) *api.RoleAccessListResult {
	res := &api.RoleAccessListResult{
		Levels: make(map[string]*api.AccessLevels),
	}
	for levelName, levels := range accessList {
		for itemName, item := range levels {
			if _, ok := res.Levels[levelName]; !ok {
				res.Levels[levelName] = &api.AccessLevels{
					Items: make(map[string]*api.AccessItem),
				}
			}
			res.Levels[levelName].Items[itemName] = &api.AccessItem{
				Actions:     item.Actions,
				Method:      item.Method,
				Description: item.Description,
				RoleName:    item.RoleName,
			}
		}
	}
	return res
}

// ToAccessListResult ...
func (r Role) ToAccessListResult(accessList access_list.AccessList) *api.AccessList {
	res := &api.AccessList{
		Levels: make(map[string]*api.AccessLevels),
	}
	for levelName, levels := range accessList {
		for itemName, item := range levels {
			if _, ok := res.Levels[levelName]; !ok {
				res.Levels[levelName] = &api.AccessLevels{
					Items: make(map[string]*api.AccessItem),
				}
			}
			res.Levels[levelName].Items[itemName] = &api.AccessItem{
				Actions:     item.Actions,
				Method:      item.Method,
				Description: item.Description,
				RoleName:    item.RoleName,
			}
		}
	}
	return res
}

// FromUpdateRoleAccessListRequest ...
func (r Role) FromUpdateRoleAccessListRequest(req *api.UpdateRoleAccessListRequest) (accessListDif map[string]map[string]bool) {

	accessListDif = make(map[string]map[string]bool)

	if req == nil {
		return
	}

	for levelName, levels := range req.AccessListDiff {
		for itemName, tr := range levels.Items {
			if _, ok := accessListDif[levelName]; !ok {
				accessListDif[levelName] = make(map[string]bool)
			}
			accessListDif[levelName][itemName] = tr
		}
	}

	return
}
