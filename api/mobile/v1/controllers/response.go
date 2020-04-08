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

package controllers

import "github.com/gin-gonic/gin"

// simple response
type Response struct {
	StatusCode int
	Payload    interface{}
}

// NewSuccess ...
func NewSuccess() *Response {
	resp := &Response{
		StatusCode: 200,
	}
	resp.Success()
	return resp
}

// Success ...
func (r *Response) Success() *Response {
	r.Payload = map[string]interface{}{}
	return r
}

// Page ...
func (r *Response) Page(limit, offset int, total int64, items interface{}) *Response {
	r.Payload = map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"limit":         limit,
			"offset":        offset,
			"objects_count": total,
		},
	}

	return r
}

// Item ...
func (r *Response) Item(name string, item interface{}) *Response {
	r.Payload = map[string]interface{}{
		name: item,
	}
	return r
}

// SetData ...
func (r *Response) SetData(data interface{}) *Response {
	r.Payload = data
	return r
}

// Send ...
func (r *Response) Send(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, r.Payload)
}
