package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerDevice struct {
	*ControllerCommon
}

func NewControllerDevice(common *ControllerCommon) *ControllerDevice {
	return &ControllerDevice{ControllerCommon: common}
}

// Device godoc
// @tags device
// @Summary Add new device
// @Description
// @Produce json
// @Accept  json
// @Param device body models.NewDevice true "device params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device [post]
// @Security ApiKeyAuth
func (c ControllerDevice) AddDevice(ctx *gin.Context) {

	var params models.NewDevice
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &m.Device{
		Id:          params.Id,
		//NodeId:      params.NodeId,
		//Baud:        params.Baud,
		//Sleep:       params.Sleep,
		Description: params.Description,
		Name:        params.Name,
		Status:      params.Status,
		//StopBite:    params.StopBite,
		//Timeout:     time.Duration(params.Timeout),
		//Tty:         params.Tty,
		IsGroup:     params.IsGroup,
	}

	if params.DeviceId != nil && *params.DeviceId != 0 {
		n.Device = &m.Device{Id: *params.DeviceId}
	}

	//if params.Address != nil && *params.Address != 0 {
	//	n.Address = params.Address
	//}

	_, id, errs, err := AddDevice(n, c.adaptors, c.core)
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
	resp.Item("id", id).Send(ctx)
}

// Device godoc
// @tags device
// @Summary Show device
// @Description Get device by id
// @Produce json
// @Accept  json
// @Param id path int true "Device ID"
// @Success 200 {object} models.ResponseDevice
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerDevice) GetDeviceById(ctx *gin.Context) {

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
	resp.Item("device", device).Send(ctx)
}

// Device godoc
// @tags device
// @Summary Update device
// @Description Update device by id
// @Produce json
// @Accept  json
// @Param  id path int true "Device ID"
// @Param  device body models.UpdateDevice true "Update device"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerDevice) UpdateDevice(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n := &m.Device{}
	if err := ctx.ShouldBindJSON(&n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n.Id = int64(aid)

	_, errs, err := UpdateDevice(n, c.adaptors, c.core)
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
	resp.Send(ctx)
}

// Device godoc
// @tags device
// @Summary Device list
// @Description Get device list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.ResponseDeviceList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /devices [Get]
// @Security ApiKeyAuth
func (c ControllerDevice) GetDeviceList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetDeviceList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, int(total), items).Send(ctx)
	return
}

// Device godoc
// @tags device
// @Summary Delete device
// @Description Delete device by id
// @Produce json
// @Accept  json
// @Param  id path int true "Device ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /device/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerDevice) DeleteDeviceById(ctx *gin.Context) {

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
