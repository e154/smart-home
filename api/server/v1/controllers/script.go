package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/scripts"
)

type ControllerScript struct {
	*ControllerCommon
	scriptService *scripts.ScriptService
}

func NewControllerScript(common *ControllerCommon,
	scriptService *scripts.ScriptService) *ControllerScript {
	return &ControllerScript{
		ControllerCommon: common,
		scriptService:    scriptService,
	}
}

// Script godoc
// @tags script
// @Summary Add new script
// @Description
// @Produce json
// @Accept  json
// @Param script body models.NewScript true "script params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /script [post]
func (c ControllerScript) AddScript(ctx *gin.Context) {

	var params models.NewScript
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	s := &m.Script{
		Lang:        common.ScriptLang(params.Lang),
		Name:        params.Name,
		Source:      params.Source,
		Description: params.Description,
	}
	_, id, errs, err := AddScript(s, c.adaptors, c.core, c.scriptService)
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

// Script godoc
// @tags script
// @Summary Show script
// @Description Get script by id
// @Produce json
// @Accept  json
// @Param id path int true "Script ID"
// @Success 200 {object} models.ResponseScript
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /script/{id} [Get]
func (c ControllerScript) GetScriptById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	script, err := GetScriptById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("script", script).Send(ctx)
}

// Script godoc
// @tags script
// @Summary Update script
// @Description Update script by id
// @Produce json
// @Accept  json
// @Param  id path int true "Script ID"
// @Param  script body models.UpdateScript true "Update script"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /script/{id} [Put]
func (c ControllerScript) UpdateScript(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	script := &m.Script{}
	if err := ctx.ShouldBindJSON(script); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	script.Id = int64(aid)

	_, errs, err := UpdateScript(script, c.adaptors, c.core, c.scriptService)
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

// Script godoc
// @tags script
// @Summary Script list
// @Description Get script list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.ResponseScriptList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /scripts [Get]
func (c ControllerScript) GetScriptList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetScriptList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// Script godoc
// @tags script
// @Summary Delete script
// @Description Delete script by id
// @Produce json
// @Accept  json
// @Param  id path int true "Script ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /script/{id} [Delete]
func (c ControllerScript) DeleteScriptById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteScriptById(int64(aid), c.adaptors, c.core); err != nil {
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

// Script godoc
// @tags script
// @Summary Execute script
// @Description Execute script by id
// @Produce json
// @Accept  json
// @Param  id path int true "Script ID"
// @Success 200 {object} models.ResponseScriptExec
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /script/{id}/exec [Post]
func (c ControllerScript) Exec(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, err := ExecuteScript(int64(aid), c.adaptors, c.core, c.scriptService)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("result", result).Send(ctx)
}

// Script godoc
// @tags device
// @Summary Search device
// @Description Search device by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.SearchScriptResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /scripts/search [Get]
func (c ControllerScript) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	scripts, _, err := SearchScript(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("scripts", scripts)
	resp.Send(ctx)
}
