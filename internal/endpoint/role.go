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
	"errors"
	"fmt"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/internal/system/rbac/access_list"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/models"
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
func (n *RoleEndpoint) Add(ctx context.Context, params *models.Role) (result *models.Role, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Role.Add(ctx, params); err != nil {
		return
	}

	if result, err = n.adaptors.Role.GetByName(ctx, params.Name); err != nil {
		return
	}

	log.Infof("added role %s", params.Name)

	return
}

// GetByName ...
func (n *RoleEndpoint) GetByName(ctx context.Context, name string) (result *models.Role, err error) {

	result, err = n.adaptors.Role.GetByName(ctx, name)

	return
}

// Update ...
func (n *RoleEndpoint) Update(ctx context.Context, params *models.Role) (result *models.Role, err error) {

	var role *models.Role
	role, err = n.adaptors.Role.GetByName(ctx, params.Name)
	if err != nil {
		return
	}

	if role.Name == "admin" && n.checkSuperUser(ctx) {
		err = apperr.ErrRoleUpdateForbidden
		return
	}

	if params.Parent.Name == "" {
		role.Parent = nil
	} else {
		role.Parent = &models.Role{
			Name: params.Parent.Name,
		}
	}

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrValidation
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

	log.Infof("updated role %s", params.Name)

	return
}

// GetList ...
func (n *RoleEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*models.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	return
}

// Delete ...
func (n *RoleEndpoint) Delete(ctx context.Context, name string) (err error) {

	if name == "admin" {
		err = apperr.ErrRoleDeleteForbidden
		return
	}

	_, err = n.adaptors.Role.GetByName(ctx, name)
	if err != nil {
		return
	}

	if err = n.adaptors.Role.Delete(ctx, name); err != nil {
		return
	}

	log.Infof("role %s was deleted", name)

	return
}

// Search ...
func (n *RoleEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*models.Role, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Role.Search(ctx, query, limit, offset)

	return
}

// GetAccessList ...
func (n *RoleEndpoint) GetAccessList(ctx context.Context, roleName string,
	accessListService access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *models.Role
	role, err = n.adaptors.Role.GetByName(ctx, roleName)
	if err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(ctx, role.Name)

	return
}

// UpdateAccessList ...
func (n *RoleEndpoint) UpdateAccessList(ctx context.Context, roleName string, accessListDif map[string]map[string]bool) (err error) {

	var role *models.Role
	role, err = n.adaptors.Role.GetByName(ctx, roleName)
	if err != nil {
		return
	}

	err = n.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		addPerms := make([]*models.Permission, 0)
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
					addPerms = append(addPerms, &models.Permission{
						RoleName:    role.Name,
						PackageName: packName,
						LevelName:   levelName,
					})
				} else if !dir && exist {
					delPerms = append(delPerms, levelName)
				}

				if len(delPerms) > 0 {
					if err = n.adaptors.Permission.Delete(ctx, roleName, packName, delPerms); err != nil {
						return err
					}
					delPerms = []string{}
				}
			}
		}

		if len(addPerms) > 0 {
			for _, perm := range addPerms {
				if _, err = n.adaptors.Permission.Add(ctx, perm); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrInternal)
	}

	return
}
