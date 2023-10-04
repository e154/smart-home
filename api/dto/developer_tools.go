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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/system/bus"
)

type DeveloperTools struct {
}

func NewDeveloperToolsDto() DeveloperTools {
	return DeveloperTools{}
}

func (DeveloperTools) GetEventBusState(state bus.Stats, total int64) (result *api.EventBusStateListResult) {
	result = &api.EventBusStateListResult{
		Items: make([]*api.BusStateItem, 0, len(state)),
		Meta: &api.Meta{
			Limit: uint64(total),
			Page:  1,
			Total: uint64(total),
		},
	}
	for _, item := range state {
		result.Items = append(result.Items, &api.BusStateItem{
			Topic:       item.Topic,
			Subscribers: int32(item.Subscribers),
		})
	}
	return
}
