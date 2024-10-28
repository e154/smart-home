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
	"github.com/e154/smart-home/internal/api/stub"
	m "github.com/e154/smart-home/pkg/models"
)

// EntityStorage ...
type EntityStorage struct{}

// NewEntityStorageDto ...
func NewEntityStorageDto() EntityStorage {
	return EntityStorage{}
}

func (_ EntityStorage) ToListResult(list *m.EntityStorageList) []*stub.ApiEntityStorage {

	var items = make([]*stub.ApiEntityStorage, 0, len(list.Attributes))

	for _, item := range list.Items {
		attributes := list.Attributes[item.EntityId].Copy()
		attributes.Deserialize(item.Attributes)
		items = append(items, &stub.ApiEntityStorage{
			Id:                item.Id,
			EntityId:          string(item.EntityId),
			EntityDescription: item.EntityDescription,
			State:             item.State,
			StateDescription:  item.StateDescription,
			Attributes:        AttributeToApi(attributes),
			CreatedAt:         item.CreatedAt,
		})
	}

	return items
}
