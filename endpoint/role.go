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
	"context"
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/go-playground/validator/v10"
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
func (n *RoleEndpoint) Add(ctx context.Context, params *m.Role) (result *m.Role, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if err = n.adaptors.Role.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Role.GetByName(params.Name)

	return
}

// GetByName ...
func (n *RoleEndpoint) GetByName(ctx context.Context, name string) (result *m.Role, err error) {

	result, err = n.adaptors.Role.GetByName(name)

	return
}

// Update ...
func (n *RoleEndpoint) Update(ctx context.Context, params *m.Role) (result *m.Role, errs validator.ValidationErrorsTranslations, err error) {

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

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if err = n.adaptors.Role.Update(role); err != nil {
		return
	}

	role, err = n.adaptors.Role.GetByName(params.Name)

	return
}

// GetList ...
func (n *RoleEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// Delete ...
func (n *RoleEndpoint) Delete(ctx context.Context, name string) (err error) {

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
func (n *RoleEndpoint) Search(ctx context.Context, query string, limit, offset int) (result []*m.Role, total int64, err error) {

	if limit == 0 {
		limit = int(common.DefaultPageSize)
	}

	result, total, err = n.adaptors.Role.Search(query, limit, offset)

	return
}

// GetAccessList ...
func (n *RoleEndpoint) GetAccessList(ctx context.Context, roleName string,
	accessListService access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(role.Name)

	return
}

// UpdateAccessList ...
func (n *RoleEndpoint) UpdateAccessList(ctx context.Context, roleName string, accessListDif map[string]map[string]bool) (err error) {

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
			tx.Permission.Add(perm)
		}
	}

	err = tx.Commit()

	return
}
