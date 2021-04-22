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

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	log = common.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	endpoint   *endpoint.Endpoint
}

// NewControllerCommon ...
func NewControllerCommon(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	endpoint *endpoint.Endpoint) *ControllerCommon {
	return &ControllerCommon{
		adaptors:   adaptors,
		accessList: accessList,
		endpoint:   endpoint,
	}
}

//query
//limit
//offset
func (c ControllerCommon) select2(ctx *gin.Context) (query string, limit, offset int) {
	query = ctx.Request.URL.Query().Get("query")
	limit, _ = strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	offset, _ = strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	return
}

//query
//sortby
//order
//limit
//offset
func (c ControllerCommon) list(ctx *gin.Context) (query, sortBy, order string, limit, offset int) {

	limit = 15
	offset = 0
	order = "DESC"
	sortBy = "created_at"

	if ctx.Request.URL.Query().Get("query") != "" {
		query = ctx.Request.URL.Query().Get("query")
	}

	if ctx.Request.URL.Query().Get("sortby") != "" {
		sortBy = ctx.Request.URL.Query().Get("sortby")
	}

	if ctx.Request.URL.Query().Get("order") != "" {
		order = ctx.Request.URL.Query().Get("order")
	}

	if ctx.Request.URL.Query().Get("limit") != "" {
		limit, _ = strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	}

	if ctx.Request.URL.Query().Get("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	}
	return
}

func (c ControllerCommon) getUser(ctx *gin.Context) (user *m.User, err error) {

	u, ok := ctx.Get("currentUser")
	if !ok {
		err = fmt.Errorf("bad user object")
		return
	}

	user, ok = u.(*m.User)
	if !ok {
		err = fmt.Errorf("bad user object")
		return
	}

	return
}
