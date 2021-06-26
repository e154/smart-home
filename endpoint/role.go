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

package endpoint

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/debug"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/validation"
)

// RoleEndpoint ...
type RoleEndpoint struct {
	*CommonEndpoint
}

// NewRoleEndpoint ...
func NewRoleEndpoint(common *CommonEndpoint) *RoleEndpoint {
	return &RoleEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *RoleEndpoint) Add(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Role.GetByName(params.Name)

	return
}

// GetByName ...
func (n *RoleEndpoint) GetByName(name string) (result *m.Role, err error) {

	result, err = n.adaptors.Role.GetByName(name)

	return
}

// Update ...
func (n *RoleEndpoint) Update(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	role, err := n.adaptors.Role.GetByName(params.Name)
	if err != nil {
		return
	}

	if params.Parent.Name == "" {
		role.Parent = nil
	} else {
		role.Parent = &m.Role{
			Name: params.Parent.Name,
		}
	}

	// validation
	_, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Update(role); err != nil {
		return
	}

	role, err = n.adaptors.Role.GetByName(params.Name)

	return
}

// GetList ...
func (n *RoleEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.Role, total int64, err error) {

	if limit == 0 {
		 limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Role.List(limit, offset, order, sortBy)

	return
}

// Delete ...
func (n *RoleEndpoint) Delete(name string) (err error) {

	if name == "admin" {
		err = errors.New("admin is base role")
		return
	}

	if _, err = n.adaptors.Role.GetByName(name); err != nil {
		return
	}

	err = n.adaptors.Role.Delete(name)

	return
}

// Search ...
func (n *RoleEndpoint) Search(query string, limit, offset int) (result []*m.Role, total int64, err error) {

	if limit == 0 {
		limit = int(common.DefaultPageSize)
	}

	result, total, err = n.adaptors.Role.Search(query, limit, offset)

	return
}

// GetAccessList ...
func (n *RoleEndpoint) GetAccessList(roleName string,
	accessListService access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(role.Name)

	return
}

// UpdateAccessList ...
func (n *RoleEndpoint) UpdateAccessList(roleName string, accessListDif map[string]map[string]bool) (err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	tx := n.adaptors.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	addPerms := make([]*m.Permission, 0)
	delPerms := make([]string, 0)
	var exist bool
	for packName, pack := range accessListDif {
		for levelName, dir := range pack {

			exist = false
			for _, lName := range role.AccessList[packName] {
				if levelName == lName {
					exist = true
				}
			}

			fmt.Println(levelName, exist)

			if dir && !exist {
				addPerms = append(addPerms, &m.Permission{
					RoleName:    role.Name,
					PackageName: packName,
					LevelName:   levelName,
				})
			} else if !dir && exist {
				delPerms = append(delPerms, levelName)
			}

			if len(delPerms) > 0 {
				if err = tx.Permission.Delete(roleName, packName, delPerms); err != nil {
					return
				}
				delPerms = []string{}
			}
		}
	}

	if len(addPerms) > 0 {
		for _, perm := range addPerms {
			fmt.Println("----2")
			debug.Println(perm)
			tx.Permission.Add(perm)
		}
	}

	err = tx.Commit()

	return
}
