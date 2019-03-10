package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerDeviceAction struct {
	*ControllerCommon
}

func NewControllerDeviceAction(common *ControllerCommon) *ControllerDeviceAction {
	return &ControllerDeviceAction{ControllerCommon: common}
}

// DeviceAction godoc
// @tags device_action
// @Summary Add new device action
// @Description
// @Produce json
// @Accept  json
// @Param device_action body models.NewDeviceAction true "device action params"
// @Success 200 {object} models.DeviceAction
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_action [post]
// @Security ApiKeyAuth
func (c ControllerDeviceAction) Add(ctx *gin.Context) {

	var params models.NewDeviceAction
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	_, id, errs, err := AddDeviceAction(params, c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	action, err := GetDeviceActionById(id, c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(action).Send(ctx)
}

// DeviceAction godoc
// @tags device_action
// @Summary Show device action
// @Description Get device action by id
// @Produce json
// @Accept  json
// @Param id path int true "DeviceAction ID"
// @Success 200 {object} models.DeviceAction
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_action/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerDeviceAction) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	action, err := GetDeviceActionById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(action).Send(ctx)
}

// DeviceAction godoc
// @tags device_action
// @Summary Update device action
// @Description Update device action by id
// @Produce json
// @Accept  json
// @Param  id path int true "DeviceAction ID"
// @Param  device action body models.UpdateDeviceAction true "Update device action"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_action/{id} [Put]
// @Security ApiKeyAuth
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

	_, errs, err := UpdateDeviceAction(params, int64(aid), c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	deviceAction, err := GetDeviceActionById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(deviceAction).Send(ctx)
}

// DeviceAction godoc
// @tags device_action
// @Summary Delete device action
// @Description Delete device action by id
// @Produce json
// @Accept  json
// @Param  id path int true "DeviceAction ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_action/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerDeviceAction) DeleteById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteDeviceActionById(int64(aid), c.adaptors, c.core); err != nil {
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

// DeviceAction godoc
// @tags device_action
// @Summary DeviceAction list
// @Description Get device list
// @Produce json
// @Accept  json
// @Param  id path int true "Device ID"
// @Success 200 {array} models.DeviceAction
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_actions/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerDeviceAction) GetActionList(ctx *gin.Context) {

	id := ctx.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	items, err := GetDeviceActionList(int64(deviceId), c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(items).Send(ctx)
	return
}

// RoleModel godoc
// @tags device_action
// @Summary Search device_action
// @Description Search device_action by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.SearchDeviceActionResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /device_action1/search [Get]
func (c ControllerDeviceAction) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	actions, _, err := SearchDeviceAction(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("actions", actions)
	resp.Send(ctx)
}