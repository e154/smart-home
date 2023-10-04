// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// EntityStorage ...
type EntityStorage struct{}

// NewEntityStorageDto ...
func NewEntityStorageDto() EntityStorage {
	return EntityStorage{}
}

func (_ EntityStorage) List(list *m.EntityStorageList, total uint64, pagination common.PageParams) (result *api.GetEntityStorageResult) {

	var items = make([]*api.EntityStorage, 0, len(list.Attributes))

	for _, item := range list.Items {
		attributes := list.Attributes[item.EntityId].Copy()
		attributes.Deserialize(item.Attributes)
		items = append(items, &api.EntityStorage{
			Id:         item.Id,
			EntityId:   string(item.EntityId),
			State:      item.State,
			Attributes: AttributeToApi(attributes),
			CreatedAt:  timestamppb.New(item.CreatedAt),
		})
	}

	return &api.GetEntityStorageResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}
