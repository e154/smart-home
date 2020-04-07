// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type AlexaApplicationEndpoint struct {
	*CommonEndpoint
}

func NewAlexaApplicationEndpoint(common *CommonEndpoint) *AlexaApplicationEndpoint {
	return &AlexaApplicationEndpoint{
		CommonEndpoint: common,
	}
}

func (n *AlexaApplicationEndpoint) Add(params *m.AlexaApplication) (result *m.AlexaApplication, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.AlexaApplication.Add(params); err != nil {
		return
	}

	if result, err = n.adaptors.AlexaApplication.GetById(id); err != nil {
		return
	}

	n.alexa.Add(result)

	return
}

func (n *AlexaApplicationEndpoint) GetById(appId int64) (result *m.AlexaApplication, err error) {

	result, err = n.adaptors.AlexaApplication.GetById(appId)

	return
}

func (n *AlexaApplicationEndpoint) Update(params *m.AlexaApplication) (app *m.AlexaApplication, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.AlexaApplication.Update(params); err != nil {
		return
	}

	if app, err = n.adaptors.AlexaApplication.GetById(params.Id); err != nil {
		return
	}

	n.alexa.Update(app)

	return
}

func (n *AlexaApplicationEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.AlexaApplication, total int64, err error) {

	result, total, err = n.adaptors.AlexaApplication.List(limit, offset, order, sortBy)

	return
}

func (n *AlexaApplicationEndpoint) Delete(appId int64) (err error) {

	if appId == 0 {
		err = errors.New("app id is null")
		return
	}

	var app *m.AlexaApplication
	if app, err = n.adaptors.AlexaApplication.GetById(appId); err != nil {
		return
	}

	if err = n.adaptors.AlexaApplication.Delete(app.Id); err != nil {
		return
	}

	n.alexa.Delete(app)

	return
}
