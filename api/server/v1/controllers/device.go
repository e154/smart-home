package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
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

	var properties []byte
	var err error
	if properties, err = json.Marshal(params.Properties); err != nil {
		return
	}

	device := &m.Device{
		Name:        params.Name,
		Description: params.Description,
		Properties:  properties,
		Status:      params.Status,
		Type:        common.DeviceType(params.Type),
	}

	if params.Device != nil && params.Device.Id != 0 {
		device.Device = &m.Device{Id: params.Device.Id}
	}

	if params.Node != nil && params.Node.Id != 0 {
		device.Node = &m.Node{Id: params.Node.Id}
	}

	device, errs, err := c.endpoint.Device.Add(device)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Device{}
	_ = common.Copy(&result, &device)

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

	device, err := c.endpoint.Device.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Device{}
	_ = common.Copy(&result, &device, common.JsonEngine)

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

	var properties []byte
	if properties, err = json.Marshal(params.Properties); err != nil {
		return
	}

	device := &m.Device{
		Id:          int64(aid),
		Name:        params.Name,
		Description: params.Description,
		Properties:  properties,
		Status:      params.Status,
		Type:        common.DeviceType(params.Type),
	}
	device, errs, err := c.endpoint.Device.Update(device)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Device{}
	_ = common.Copy(&result, &device, common.JsonEngine)

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
	devices, total, err := c.endpoint.Device.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Device, 0)
	_ = common.Copy(&result, &devices, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
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

	if err := c.endpoint.Device.Delete(int64(aid)); err != nil {
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
	devices, _, err := c.endpoint.Device.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.DeviceShort, 0)
	_ = common.Copy(&result, &devices, common.JsonEngine)

	resp := NewSuccess()
	resp.Item("devices", result)
	resp.Send(ctx)
}
