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
	"strings"
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// LogEndpoint ...
type LogEndpoint struct {
	*CommonEndpoint
}

// NewLogEndpoint ...
func NewLogEndpoint(common *CommonEndpoint) *LogEndpoint {
	return &LogEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (l *LogEndpoint) Add(ctx context.Context, log *m.Log) (result *m.Log, err error) {

	if ok, errs := l.validation.Valid(log); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = l.adaptors.Log.Add(ctx, log); err != nil {
		return
	}

	result, err = l.adaptors.Log.GetById(ctx, id)

	return
}

// GetById ...
func (l *LogEndpoint) GetById(ctx context.Context, id int64) (log *m.Log, err error) {

	log, err = l.adaptors.Log.GetById(ctx, id)

	return
}

// GetList ...
func (l *LogEndpoint) GetList(ctx context.Context, pagination common.PageParams, query *string, startDate, endDate *time.Time) (list []*m.Log, total int64, err error) {

	var queryObj = &m.LogQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Levels:    nil,
	}
	//if startDate != nil {
	//	date, _ := time.Parse("2006-01-02T15:04:05.999Z07", *startDate)
	//	queryObj.StartDate = &date
	//}
	//if endDate != nil {
	//	date, _ := time.Parse("2006-01-02T15:04:05.999Z07", *endDate)
	//	queryObj.EndDate = &date
	//}
	if query != nil {
		queryObj.Levels = strings.Split(strings.Replace(*query, "'", "", -1), ",")
	}

	list, total, err = l.adaptors.Log.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, queryObj)

	return
}

// Search ...
func (l *LogEndpoint) Search(ctx context.Context, query string, limit, offset int) (list []*m.Log, total int64, err error) {

	list, total, err = l.adaptors.Log.Search(ctx, query, limit, offset)

	return
}

// Delete ...
func (l *LogEndpoint) Delete(ctx context.Context, logId int64) (err error) {

	_, err = l.adaptors.Log.GetById(ctx, logId)
	if err != nil {
		return
	}

	err = l.adaptors.Log.Delete(ctx, logId)

	return
}
