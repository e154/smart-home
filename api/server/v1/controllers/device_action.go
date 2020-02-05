package controllers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
)

type ControllerDeviceAction struct {
	*ControllerCommon
}

func NewControllerDeviceAction(common *ControllerCommon) *ControllerDeviceAction {
	return &ControllerDeviceAction{ControllerCommon: common}
}

// swagger:operation POST /device_action deviceActionAdd
// ---
// parameters:
// - description: device action params
//   in: body
//   name: device_action
//   required: true
//   schema:
//     $ref: '#/definitions/NewDeviceAction'
//     type: object
// summary: add new device action
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
// responses:
//   "200":
//     schema:
//       $ref: '#/definitions/DeviceAction'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceAction) Add(ctx *gin.Context) {

	var params models.NewDeviceAction
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	action := &m.DeviceAction{}
	_ = common.Copy(&action, &params, common.JsonEngine)

	if params.Device != nil && params.Device.Id != 0 {
		action.DeviceId = params.Device.Id
	}
	if params.Script != nil && params.Script.Id != 0 {
		action.ScriptId = params.Script.Id
	}

	action, errs, err := c.endpoint.DeviceAction.Add(action)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.DeviceAction{}
	_ = common.Copy(&result, &action)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /device_action/{id} deviceActionGetById
// ---
// parameters:
// - description: DeviceAction ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get device action by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceAction'
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
func (c ControllerDeviceAction) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	action, err := c.endpoint.DeviceAction.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.DeviceAction{}
	_ = common.Copy(&result, &action)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /device_action/{id} deviceActionUpdateById
// ---
// parameters:
// - description: DeviceAction ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update device action params
//   in: body
//   name: device_action
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateDeviceAction'
//     type: object
// summary: update device action by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
// responses:
//   "200":
//     schema:
//       $ref: '#/definitions/DeviceAction'
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
func (c ControllerDeviceAction) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var params models.UpdateDeviceAction
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)
	action := &m.DeviceAction{}
	_ = copier.Copy(&action, &params)

	if params.Device != nil && params.Device.Id != 0 {
		action.DeviceId = params.Device.Id
	}
	if params.Script != nil && params.Script.Id != 0 {
		action.ScriptId = params.Script.Id
	}

	action, errs, err := c.endpoint.DeviceAction.Update(action)
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

	result := &models.DeviceAction{}
	_ = common.Copy(&result, &action)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation DELETE /device_action/{id} deviceActionDeleteById
// ---
// parameters:
// - description: DeviceAction ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete device action by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
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
func (c ControllerDeviceAction) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.DeviceAction.Delete(int64(aid)); err != nil {
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

// swagger:operation GET /device_actions/{id} deviceActionGetListByDeviceId
// ---
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get device actions by device id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
// responses:
//   "200":
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/DeviceAction'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceAction) GetActionList(ctx *gin.Context) {

	id := ctx.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	actions, err := c.endpoint.DeviceAction.GetList(int64(deviceId))
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	//result := make([]*models.DeviceAction, 0)
	//common.Copy(&result, &actions)

	resp := NewSuccess()
	resp.SetData(actions).Send(ctx)
}

// swagger:operation GET /device_action/search deviceActionSearch
// ---
// summary: search device actions
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_action
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
//	   $ref: '#/responses/DeviceActionSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceAction) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	actions, _, err := c.endpoint.DeviceAction.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.DeviceAction, 0)
	_ = common.Copy(&result, &actions)

	resp := NewSuccess()
	resp.Item("actions", actions)
	resp.Send(ctx)
}
