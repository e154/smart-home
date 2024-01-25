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

package access_list

import (
	"context"
	"embed"
	"encoding/json"

	"github.com/e154/smart-home/common/logger"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

var (
	log = logger.MustGetLogger("access_list")
)

//go:embed data.json
var DATA embed.FS

// AccessListService ...
type AccessListService interface {
	ReadConfig(ctx context.Context) (err error)
	GetFullAccessList(ctx context.Context, roleName string) (accessList AccessList, err error)
	GetShotAccessList(ctx context.Context, role *m.Role) (err error)
	List(ctx context.Context) *AccessList
}

// accessListService ...
type accessListService struct {
	list     *AccessList
	adaptors *adaptors.Adaptors
}

// NewAccessListService ...
func NewAccessListService(adaptors *adaptors.Adaptors) AccessListService {
	accessList := &accessListService{
		adaptors: adaptors,
	}
	_ = accessList.ReadConfig(context.Background())
	return accessList
}

// ReadConfig ...
func (a *accessListService) ReadConfig(ctx context.Context) (err error) {

	//var file []byte
	//file, err = ioutil.ReadFile(path)
	//if err != nil {
	//	log.Fatal("Error reading config file")
	//	return
	//}

	a.list = &AccessList{}
	data, _ := DATA.ReadFile("data.json")
	err = json.Unmarshal(data, a.list)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	return
}

// GetFullAccessList ...
func (a *accessListService) GetFullAccessList(ctx context.Context, roleName string) (accessList AccessList, err error) {

	var permissions []*m.Permission
	if permissions, err = a.adaptors.Permission.GetAllPermissions(ctx, roleName); err != nil {
		return
	}

	accessList = make(AccessList)
	var item AccessItem
	var levels AccessLevels
	var ok bool
	list := *a.list
	for _, perm := range permissions {

		if levels, ok = list[perm.PackageName]; !ok {
			continue
		}

		if accessList[perm.PackageName] == nil {
			accessList[perm.PackageName] = NewAccessLevels()
		}

		if item, ok = levels[perm.LevelName]; !ok {
			continue
		}

		item.RoleName = perm.RoleName
		accessList[perm.PackageName][perm.LevelName] = item
	}

	return
}

// GetShotAccessList ...
func (a *accessListService) GetShotAccessList(ctx context.Context, role *m.Role) (err error) {

	err = a.adaptors.Role.GetAccessList(ctx, role)
	return
}

// List ...
func (a *accessListService) List(ctx context.Context) *AccessList {
	return a.list
}
