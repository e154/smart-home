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
	"encoding/json"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
	"time"
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
func (l *LogEndpoint) Add(log *m.Log) (result *m.Log, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = l.validation.Valid(log); !ok {
		return
	}

	var id int64
	if id, err = l.adaptors.Log.Add(log); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if log, err = l.adaptors.Log.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetById ...
func (l *LogEndpoint) GetById(id int64) (log *m.Log, err error) {

	if log, err = l.adaptors.Log.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetList ...
func (l *LogEndpoint) GetList(limit, offset int64, order, sortBy, query string) (list []*m.Log, total int64, err error) {

	var queryObj *m.LogQuery
	if query != "" {
		queryObj = &m.LogQuery{}
		d := make(map[string]string, 0)
		if err = json.Unmarshal([]byte(query), &d); err != nil {
			err = errors.Wrap(common.ErrInternal, err.Error())
			return
		}

		if startDate, ok := d["start_date"]; ok {
			date, _ := time.Parse("2006-01-02", startDate)
			queryObj.StartDate = &date
		}
		if endDate, ok := d["end_date"]; ok {
			date, _ := time.Parse("2006-01-02", endDate)
			queryObj.EndDate = &date
		}
		if levels, ok := d["levels"]; ok {
			queryObj.Levels = strings.Split(strings.Replace(levels, "'", "", -1), ",")
		}
	}

	list, total, err = l.adaptors.Log.List(limit, offset, order, sortBy, queryObj)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// Search ...
func (l *LogEndpoint) Search(query string, limit, offset int) (list []*m.Log, total int64, err error) {

	list, total, err = l.adaptors.Log.Search(query, limit, offset)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// Delete ...
func (l *LogEndpoint) Delete(logId int64) (err error) {

	_, err = l.adaptors.Log.GetById(logId)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = l.adaptors.Log.Delete(logId)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}
