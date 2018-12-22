package controllers

import (
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/e154/smart-home/system/access_list"
)

var (
	log = logging.MustGetLogger("controllers")
)

type ControllerCommon struct {
	adaptors   *adaptors.Adaptors
	core       *core.Core
	accessList *access_list.AccessListService
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	core *core.Core,
	accessList *access_list.AccessListService) *ControllerCommon {
	return &ControllerCommon{
		adaptors:   adaptors,
		core:       core,
		accessList: accessList,
	}
}

func (c ControllerCommon) query(ctx *gin.Context, query string) string {
	return ctx.Request.URL.Query().Get(query)
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
func (c ControllerCommon) list(ctx *gin.Context) (query, sortby, order string, limit, offset int) {
	query = ctx.Request.URL.Query().Get("query")
	sortby = ctx.Request.URL.Query().Get("sortby")
	order = ctx.Request.URL.Query().Get("order")
	limit, _ = strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	offset, _ = strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	return
}
