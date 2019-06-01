package controllers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
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

// swagger:operation POST /script scriptAdd
// ---
// parameters:
// - description: script params
//   in: body
//   name: script
//   required: true
//   schema:
//     $ref: '#/definitions/NewScript'
//     type: object
// summary: add new script
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Script'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerScript) Add(ctx *gin.Context) {

	var params models.NewScript
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	script := &m.Script{
		Lang:        common.ScriptLang(params.Lang),
		Name:        params.Name,
		Source:      params.Source,
		Description: params.Description,
	}

	script, errs, err := c.endpoint.Script.Add(script)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Script{}
	common.Copy(&result, &script, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /script/{id} scriptGetById
// ---
// parameters:
// - description: Script ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get script by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Script'
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
func (c ControllerScript) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	script, err := c.endpoint.Script.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Script{}
	common.Copy(&result, &script, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /script/{id} scriptUpdateById
// ---
// parameters:
// - description: Script ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update script params
//   in: body
//   name: script
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateScript'
//     type: object
// summary: update script by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Script'
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
func (c ControllerScript) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateScript{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	script := &m.Script{}
	common.Copy(&script, &params, common.JsonEngine)

	script, errs, err := c.endpoint.Script.Update(script)
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

	result := &models.Script{}
	common.Copy(&result, &script, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /scripts scriptList
// ---
// summary: get script list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
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
//	   $ref: '#/responses/ScriptList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerScript) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Script.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Script, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /script/{id} scriptDeleteById
// ---
// parameters:
// - description: Script ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete script by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
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
func (c ControllerScript) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Script.DeleteScriptById(int64(aid)); err != nil {
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

// swagger:operation POST /script/{id}/exec scriptExecById
// ---
// parameters:
// - description: Script ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: Execute script by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
// responses:
//   "200":
//	   $ref: '#/responses/ScriptExec'
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
func (c ControllerScript) Exec(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, err := c.endpoint.Script.Execute(int64(aid))
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

// swagger:operation POST /script/{id}/exec_src scriptExecSrc
// ---
// parameters:
// - description: source script
//   in: body
//   name: script
//   required: true
//   schema:
//     $ref: '#/definitions/ExecScript'
//     type: object
// summary: Exec script from request params
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
// responses:
//   "200":
//	   $ref: '#/responses/ScriptExec'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerScript) ExecSrc(ctx *gin.Context) {

	script := &m.Script{}
	if err := ctx.ShouldBindJSON(script); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	result, err := c.endpoint.Script.ExecuteSource(script)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("result", result).Send(ctx)
}

// swagger:operation GET /scripts/search scriptSearch
// ---
// summary: search script
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - script
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
//	   $ref: '#/responses/ScriptSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerScript) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Script.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Script, 0)
	common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Item("scripts", result)
	resp.Send(ctx)
}
