package controllers

import (
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/gin-gonic/gin"
)

var (
	log = logging.MustGetLogger("controllers")
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	core     *core.Core
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	core *core.Core) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		core:     core,
	}
}

func (c ControllerCommon) query(ctx *gin.Context, query string) string {
	return ctx.Request.URL.Query().Get(query)
}