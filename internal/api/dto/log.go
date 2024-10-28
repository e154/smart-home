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
	m "github.com/e154/smart-home/pkg/models"
)

// Log ...
type Log struct{}

// NewLogDto ...
func NewLogDto() Log {
	return Log{}
}

// ToListResult ...
func (r Log) ToListResult(list []*m.Log) []*stub.ApiLog {

	items := make([]*stub.ApiLog, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToLog(i))
	}

	return items
}

// ToLog ...
func (r Log) ToLog(log *m.Log) (obj *stub.ApiLog) {
	obj = ToLog(log)
	return
}

// ToLog ...
func ToLog(log *m.Log) (obj *stub.ApiLog) {
	if log == nil {
		return
	}
	obj = &stub.ApiLog{
		Id:        log.Id,
		Body:      log.Body,
		Level:     string(log.Level),
		Owner:     log.Owner,
		CreatedAt: log.CreatedAt,
	}
	return
}
