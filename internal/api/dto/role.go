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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/internal/system/rbac/access_list"
	m "github.com/e154/smart-home/pkg/models"
)

// Role ...
type Role struct{}

// NewRoleDto ...
func NewRoleDto() Role {
	return Role{}
}

// FromNewRoleRequest ...
func (r Role) FromNewRoleRequest(from *stub.ApiNewRoleRequest) (to *m.Role) {
	to = &m.Role{
		Name:        from.Name,
		Description: from.Description,
	}
	if from.Parent != nil {
		to.Parent = &m.Role{
			Name: *from.Parent,
		}
	}
	return
}

// FromUpdateRoleRequest ...
func (r Role) FromUpdateRoleRequest(from *stub.RoleServiceUpdateRoleByNameJSONBody, name string) (to *m.Role) {
	to = &m.Role{
		Name:        name,
		Description: from.Description,
	}
	if from.Parent != nil {
		to.Parent = &m.Role{
			Name: *from.Parent,
		}
	}
	return
}

// ToSearchResult ...
func (r Role) ToSearchResult(list []*m.Role) *stub.ApiSearchRoleListResult {

	items := make([]stub.ApiRole, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetStubRole(i))
	}

	return &stub.ApiSearchRoleListResult{
		Items: items,
	}
}

// ToListResult ...
func (r Role) ToListResult(list []*m.Role) []stub.ApiRole {

	items := make([]stub.ApiRole, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetStubRole(i))
	}

	return items
}

// GetStubRole ...
func (r Role) GetStubRole(from *m.Role) (to stub.ApiRole) {
	to = stub.ApiRole{
		Name:        from.Name,
		Description: from.Description,
		AccessList:  nil,
		CreatedAt:   from.CreatedAt,
		UpdatedAt:   from.UpdatedAt,
	}
	if from.Parent != nil {
		role := r.GetStubRole(from.Parent)
		to.Parent = &role
	}
	if len(from.Children) > 0 {
		for _, ch := range from.Children {
			to.Children = append(to.Children, r.GetStubRole(ch))
		}
	}
	if from.AccessList != nil {
		to.AccessList = &stub.ApiRoleAccessList{
			Levels: make(map[string]stub.AccessListListOfString),
		}
		for levelName, levels := range from.AccessList {
			if len(levels) > 0 {
				var items []string
				items = append(items, levels...)
				if _, ok := to.AccessList.Levels[levelName]; !ok {
					to.AccessList.Levels[levelName] = stub.AccessListListOfString{
						Items: items,
					}
				}
			}
		}
	}
	return
}

// ToRoleAccessListResult ...
func (r Role) ToRoleAccessListResult(accessList access_list.AccessList) *stub.ApiRoleAccessListResult {
	res := &stub.ApiRoleAccessListResult{
		Levels: make(map[string]stub.ApiAccessLevels),
	}
	for levelName, levels := range accessList {
		for itemName, item := range levels {
			if _, ok := res.Levels[levelName]; !ok {
				res.Levels[levelName] = stub.ApiAccessLevels{
					Items: make(map[string]stub.ApiAccessItem),
				}
			}
			res.Levels[levelName].Items[itemName] = stub.ApiAccessItem{
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
func (r Role) ToAccessListResult(accessList *access_list.AccessList) *stub.ApiAccessList {
	res := &stub.ApiAccessList{
		Levels: make(map[string]stub.ApiAccessLevels),
	}
	for levelName, levels := range *accessList {
		for itemName, item := range levels {
			if _, ok := res.Levels[levelName]; !ok {
				res.Levels[levelName] = stub.ApiAccessLevels{
					Items: make(map[string]stub.ApiAccessItem),
				}
			}
			res.Levels[levelName].Items[itemName] = stub.ApiAccessItem{
				Actions:     item.Actions,
				Method:      item.Method,
				Description: item.Description,
				RoleName:    item.RoleName,
			}
		}
	}
	return res
}

// UpdateRoleAccessList ...
func (r Role) UpdateRoleAccessList(req *stub.RoleServiceUpdateRoleAccessListJSONBody, name string) (accessListDif map[string]map[string]bool) {

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
