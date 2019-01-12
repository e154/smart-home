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

// DeviceState godoc
// @tags device_state
// @Summary Add new device state
// @Description
// @Produce json
// @Accept  json
// @Param device_state body models.NewDeviceState true "device state params"
// @Success 200 {object} models.DeviceState
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_state [post]
// @Security ApiKeyAuth
func (c ControllerDeviceState) AddState(ctx *gin.Context) {

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

// DeviceState godoc
// @tags device_state
// @Summary Show device state
// @Description Get device state by id
// @Produce json
// @Accept  json
// @Param id path int true "DeviceState ID"
// @Success 200 {object} models.DeviceState
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_state/{id} [Get]
// @Security ApiKeyAuth
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

// DeviceState godoc
// @tags device_state
// @Summary Update device state
// @Description Update device state by id
// @Produce json
// @Accept  json
// @Param  id path int true "DeviceState ID"
// @Param  device state body models.UpdateDeviceState true "Update device state"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_state/{id} [Put]
// @Security ApiKeyAuth
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

	_, errs, err := UpdateDeviceState(params, int64(aid), c.adaptors, c.core)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	deviceState, err := GetDeviceStateById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(deviceState).Send(ctx)
}

// DeviceState godoc
// @tags device_state
// @Summary Delete device state
// @Description Delete device state by id
// @Produce json
// @Accept  json
// @Param  id path int true "DeviceState ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_state/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerDeviceState) DeleteById(ctx *gin.Context) {

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

// DeviceState godoc
// @tags device_state
// @Summary DeviceState list
// @Description Get device list
// @Produce json
// @Accept  json
// @Param  id path int true "Device ID"
// @Success 200 {array} models.DeviceState
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device_states/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerDeviceState) GetDeviceStateList(ctx *gin.Context) {

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
