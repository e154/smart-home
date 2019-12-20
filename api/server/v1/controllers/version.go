package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
)

type ControllerVersion struct {
	*ControllerCommon
}

func NewControllerVersion(common *ControllerCommon) *ControllerVersion {
	return &ControllerVersion{ControllerCommon: common}
}

// swagger:operation GET /version getServerVersion
// ---
// summary: get server version
// description:
// tags:
// - version
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Version'
func (c ControllerVersion) Version(ctx *gin.Context) {

	serverVersion := c.endpoint.Version.ServerVersion()

	result := &models.Version{}
	common.Copy(&result, &serverVersion)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
