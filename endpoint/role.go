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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/pkg/errors"
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
func (n *RoleEndpoint) Add(ctx context.Context, params *m.Role) (result *m.Role, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Role.Add(ctx, params); err != nil {
		return
	}

	result, err = n.adaptors.Role.GetByName(ctx, params.Name)

	return
}

// GetByName ...
func (n *RoleEndpoint) GetByName(ctx context.Context, name string) (result *m.Role, err error) {

	result, err = n.adaptors.Role.GetByName(ctx, name)

	return
}

// Update ...
func (n *RoleEndpoint) Update(ctx context.Context, params *m.Role) (result *m.Role, err error) {

	var role *m.Role
	role, err = n.adaptors.Role.GetByName(ctx, params.Name)
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

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Role.Update(ctx, role); err != nil {
		return
	}

	role, err = n.adaptors.Role.GetByName(ctx, params.Name)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}
	return
}

// GetList ...
func (n *RoleEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	return
}

// Delete ...
func (n *RoleEndpoint) Delete(ctx context.Context, name string) (err error) {

	if name == "admin" {
		err = apperr.ErrBadRequestParams
		return
	}

	_, err = n.adaptors.Role.GetByName(ctx, name)
	if err != nil {
		return
	}

	err = n.adaptors.Role.Delete(ctx, name)

	return
}

// Search ...
func (n *RoleEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Role, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Role.Search(ctx, query, limit, offset)

	return
}

// GetAccessList ...
func (n *RoleEndpoint) GetAccessList(ctx context.Context, roleName string,
	accessListService access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	role, err = n.adaptors.Role.GetByName(ctx, roleName)
	if err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(ctx, role.Name)

	return
}

// UpdateAccessList ...
func (n *RoleEndpoint) UpdateAccessList(ctx context.Context, roleName string, accessListDif map[string]map[string]bool) (err error) {

	var role *m.Role
	role, err = n.adaptors.Role.GetByName(ctx, roleName)
	if err != nil {
		return
	}

	tx := n.adaptors.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
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
				if err = tx.Permission.Delete(ctx, roleName, packName, delPerms); err != nil {
					return
				}
				delPerms = []string{}
			}
		}
	}

	if len(addPerms) > 0 {
		for _, perm := range addPerms {
			_, _ = tx.Permission.Add(ctx, perm)
		}
	}

	if err = tx.Commit(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}

	return
}
