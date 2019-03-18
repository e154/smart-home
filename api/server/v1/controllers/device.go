package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerDevice struct {
	*ControllerCommon
}

func NewControllerDevice(common *ControllerCommon) *ControllerDevice {
	return &ControllerDevice{ControllerCommon: common}
}

// swagger:operation POST /device deviceAdd
// ---
// parameters:
// - description: device params
//   in: body
//   name: device
//   required: true
//   schema:
//     $ref: '#/definitions/NewDevice'
//     type: object
// summary: add new device
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Device'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDevice) Add(ctx *gin.Context) {

	var params models.NewDevice
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, errs, err := AddDevice(params, c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /device/{id} deviceGetById
// ---
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get device by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Device'
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
func (c ControllerDevice) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	device, err := GetDeviceById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(device).Send(ctx)
}

// swagger:operation PUT /device/{id} deviceUpdateById
// ---
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update device params
//   in: body
//   name: device
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateDevice'
//     type: object
// summary: update device by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Device'
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
func (c ControllerDevice) UpdateDevice(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var params models.UpdateDevice
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, errs, err := UpdateDevice(params, int64(aid), c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /devices deviceList
// ---
// summary: get device list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// - default: DESC
//   description: order
//   in: query
//   name: order
//   type: string
// - default: id
//   description: sort_by
//   in: query
//   name: sort_by
//   type: string
// responses:
//   "200":
//	   $ref: '#/responses/DeviceList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDevice) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetDeviceList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// swagger:operation DELETE /device/{id} deviceDeleteById
// ---
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete device by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
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
func (c ControllerDevice) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteDeviceById(int64(aid), c.adaptors, c.core); err != nil {
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

// swagger:operation GET /devices/search deviceSearch
// ---
// summary: search device
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device
// parameters:
// - description: query
//   in: query
//   name: query
//   type: string
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/DeviceSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDevice) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	devices, _, err := SearchDevice(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("devices", devices)
	resp.Send(ctx)
}
