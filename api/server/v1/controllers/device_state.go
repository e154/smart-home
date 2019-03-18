package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerDeviceState struct {
	*ControllerCommon
}

func NewControllerDeviceState(common *ControllerCommon) *ControllerDeviceState {
	return &ControllerDeviceState{ControllerCommon: common}
}

// swagger:operation POST /device_state deviceStateAdd
// ---
// parameters:
// - description: device state params
//   in: body
//   name: device_state
//   required: true
//   schema:
//     $ref: '#/definitions/NewDeviceState'
//     type: object
// summary: add new device state
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) Add(ctx *gin.Context) {

	var params models.NewDeviceState
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	_, id, errs, err := AddDeviceState(params, c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	state, err := GetDeviceStateById(id, c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(state).Send(ctx)
}

// swagger:operation GET /device_state/{id} deviceStateGetById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
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
func (c ControllerDeviceState) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	state, err := GetDeviceStateById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(state).Send(ctx)
}

// swagger:operation PUT /device_state/{id} deviceStateUpdateById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update device state params
//   in: body
//   name: device_state
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateDeviceState'
//     type: object
// summary: update device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
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
func (c ControllerDeviceState) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var params models.UpdateDeviceState
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	deviceState, errs, err := UpdateDeviceState(params, int64(aid), c.adaptors, c.core)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(deviceState).Send(ctx)
}

// swagger:operation DELETE /device_state/{id} deviceStateDeleteById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
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
func (c ControllerDeviceState) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteDeviceStateById(int64(aid), c.adaptors, c.core); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /device_states/{id} deviceStateList
// ---
// summary: get device state list by device id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// responses:
//   "200":
//     description: OK
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/DeviceState'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) GetStateList(ctx *gin.Context) {

	id := ctx.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	items, err := GetDeviceStateList(int64(deviceId), c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(items).Send(ctx)
	return
}
