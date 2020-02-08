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

package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Image struct {
	Id        int64     `json:"id"`
	Thumb     string    `json:"thumb"`
	Url       string    `json:"url"`
	Image     string    `json:"image"`
	MimeType  string    `json:"mime_type"`
	Title     string    `json:"title"`
	Size      int64     `json:"size"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Image) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

type ImageFilterList struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
