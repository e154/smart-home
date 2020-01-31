package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/gin-gonic/gin"
)

type ControllerGate struct {
	*ControllerCommon
}

func NewControllerGate(common *ControllerCommon) *ControllerGate {
	return &ControllerGate{ControllerCommon: common}
}

// swagger:operation GET /gate gateGetSettings
// ---
// parameters:
// summary: get gate settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - gate
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/GateSettings'
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
func (c ControllerGate) GetSettings(ctx *gin.Context) {

	settings, err := c.endpoint.Gate.GetSettings()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.GateSettings{}
	_ = common.Copy(&result, &settings)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /gate gateUpdateSettings
// ---
// parameters:
// - description: Update gate params
//   in: body
//   name: user
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateGateSettings'
//     type: object
// summary: update gate settings
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - gate
// responses:
//   "200":
//     $ref: '#/responses/Success'
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
func (c ControllerGate) UpdateSettings(ctx *gin.Context) {

	n := &models.UpdateGateSettings{}
	if err := ctx.ShouldBindJSON(n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	settings := gate_client.Settings{}
	_ = common.Copy(&settings, &n)

	if err := c.endpoint.Gate.UpdateSettings(settings); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /gate/mobiles gateGetMobileList
// ---
// parameters:
// summary: get gate mobile list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - gate
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/GateMobileList'
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
func (c ControllerGate) GetMobileList(ctx *gin.Context) {

	list, err := c.endpoint.Gate.GetMobileList(ctx)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.GateMobileList{}
	_ = common.Copy(&result, &list)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation DELETE /gate/mobile/{token} gateDeleteMobile
// ---
// parameters:
// - description: mobile token
//   in: path
//   name: token
//   required: true
//   type: string
// summary: delete mobile by token
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - gate
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
func (c ControllerGate) DeleteMobile(ctx *gin.Context) {

	token := ctx.Param("token")
	if token == "" {
		NewError(404, "record not found").Send(ctx)
		return
	}

	if _, err := c.endpoint.Gate.DeleteMobile(token, ctx); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation POST /gate/mobile gateAddMobile
// ---
// parameters:
// summary: add new mobile client
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - gate
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerGate) AddMobile(ctx *gin.Context) {

	if _, err := c.endpoint.Gate.AddMobile(ctx); err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}
