package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/notify"
	"github.com/gin-gonic/gin"
)

type ControllerNotifr struct {
	*ControllerCommon
}

func NewControllerNotifr(common *ControllerCommon) *ControllerNotifr {
	return &ControllerNotifr{ControllerCommon: common}
}

// swagger:operation PUT /notifr/config notifyUpdateSettings
// ---
// parameters:
// - description: Update notifr params
//   in: body
//   name: notifr
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateNotifr'
//     type: object
// summary: update notifr settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNotifr) Update(ctx *gin.Context) {

	params := &models.UpdateNotifr{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	settings := &notify.NotifyConfig{}
	_ = common.Copy(&settings, &params, common.JsonEngine)

	err := c.endpoint.Notify.UpdateSettings(settings)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /notifr/config notifyGetSettings
// ---
// parameters:
// summary: get notifr settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notifr
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Notifr'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerNotifr) GetSettings(ctx *gin.Context) {

	settings, err := c.endpoint.Notify.GetSettings()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Notify{}
	common.Copy(&result, &settings, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
