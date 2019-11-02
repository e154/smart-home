package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	notify2 "github.com/e154/smart-home/system/notify"
	"github.com/gin-gonic/gin"
)

type ControllerNotify struct {
	*ControllerCommon
}

func NewControllerNotify(common *ControllerCommon) *ControllerNotify {
	return &ControllerNotify{ControllerCommon: common}
}

// swagger:operation PUT /notify/config notifyUpdateSettings
// ---
// parameters:
// - description: Update notify params
//   in: body
//   name: notify
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateNotify'
//     type: object
// summary: update notify settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notify
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
func (c ControllerNotify) Update(ctx *gin.Context) {

	params := &models.UpdateNotify{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	config := &notify2.NotifyConfig{}
	_ = common.Copy(&config, &params, common.JsonEngine)

	err := c.endpoint.Notify.UpdateSettings(config)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /notify/config notifyGetSettings
// ---
// parameters:
// summary: get notify settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - notify
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Notify'
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
func (c ControllerNotify) GetSettings(ctx *gin.Context) {

	notify, err := c.endpoint.Notify.GetSettings()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Notify{}
	common.Copy(&result, &notify, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
