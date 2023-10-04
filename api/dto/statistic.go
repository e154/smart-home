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
	m "github.com/e154/smart-home/models"
)

func GetStatistic(statistic []*m.Statistic) (result *api.Statistics) {
	result = &api.Statistics{
		Items: make([]*api.Statistic, 0, len(statistic)),
	}
	for _, item := range statistic {
		result.Items = append(result.Items, &api.Statistic{
			Name:        item.Name,
			Description: item.Description,
			Value:       item.Value,
			Diff:        item.Diff,
		})
	}
	return
}
