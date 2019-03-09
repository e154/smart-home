package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerLog struct {
	*ControllerCommon
}

func NewControllerLog(common *ControllerCommon) *ControllerLog {
	return &ControllerLog{ControllerCommon: common}
}

// Log godoc
// @tags log
// @Summary Add new log
// @Description
// @Produce json
// @Accept  json
// @Param log body models.NewLogModel true "log params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /log [post]
// @Security ApiKeyAuth
func (c ControllerLog) AddLog(ctx *gin.Context) {

	log := &models.NewLogModel{}
	if err := ctx.ShouldBindJSON(&log); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	id, errs, err := AddLog(log, c.adaptors)
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

// Log godoc
// @tags log
// @Summary Show log
// @Description Get log by id
// @Produce json
// @Accept  json
// @Param id path int true "Log ID"
// @Success 200 {object} models.ResponseLog
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /log/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerLog) GetLogById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	log, err := GetLogById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("log", log).Send(ctx)
}

// Log godoc
// @tags log
// @Summary Log list
// @Description Get log list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Param query query string false "query"
// @Success 200 {object} models.ResponseLogList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /logs [Get]
// @Security ApiKeyAuth
func (c ControllerLog) GetLogList(ctx *gin.Context) {

	query, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetLogList(int64(limit), int64(offset), order, sortBy, query, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
}

// Log godoc
// @tags log
// @Summary Delete log
// @Description Delete log by id
// @Produce json
// @Accept  json
// @Param  id path int true "Log ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /log/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerLog) DeleteLogById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = DeleteLogById(int64(aid), c.adaptors); err != nil {
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

// Log godoc
// @tags log
// @Summary Search log
// @Description Search log by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.ResponseSearchLog
// @Failure 401 "Unauthorized"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /logs/search [Get]
func (c ControllerLog) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	logs, _, err := SearchLog(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("logs", logs)
	resp.Send(ctx)
}
