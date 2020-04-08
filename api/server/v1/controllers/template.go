// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
)

// ControllerTemplate ...
type ControllerTemplate struct {
	*ControllerCommon
}

// NewControllerTemplate ...
func NewControllerTemplate(common *ControllerCommon) *ControllerTemplate {
	return &ControllerTemplate{ControllerCommon: common}
}

// swagger:operation POST /template templateAdd
// ---
// parameters:
// - description: template params
//   in: body
//   name: template
//   required: true
//   schema:
//     $ref: '#/definitions/NewTemplate'
//     type: object
// summary: add new template item
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplate) Add(ctx *gin.Context) {

	params := &models.NewTemplate{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	template := &m.Template{}
	_ = common.Copy(&template, &params, common.JsonEngine)
	template.Type = m.TemplateTypeTemplate

	errs, err := c.endpoint.Template.UpdateOrCreate(template)
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

// swagger:operation GET /template/{name} templateGetByName
// ---
// parameters:
// - description: Template Name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get template by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Template'
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
func (c ControllerTemplate) GetByName(ctx *gin.Context) {

	name := ctx.Param("name")
	if name == "" {
		NewError(400, "bad param name").Send(ctx)
		return
	}

	template, err := c.endpoint.Template.GetByName(name)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err.Error()).Send(ctx)
		return
	}

	result := &models.Template{}
	_ = common.Copy(&result, &template, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /template/{name} templateUpdateByName
// ---
// parameters:
// - description: Template Name
//   in: path
//   name: name
//   required: true
//   type: string
// - description: Update item params
//   in: body
//   name: template
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateTemplate'
//     type: object
// summary: update template by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Template'
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
func (c ControllerTemplate) Update(ctx *gin.Context) {

	name := ctx.Param("name")
	if name == "" {
		NewError(400, "bad param name").Send(ctx)
		return
	}

	params := &models.UpdateTemplate{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	template := &m.Template{}
	_ = common.Copy(&template, &params, common.JsonEngine)

	errs, err := c.endpoint.Template.UpdateOrCreate(template)
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
	resp.SetData(template).Send(ctx)
}

// swagger:operation GET /templates templateList
// ---
// summary: get template list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
// responses:
//   "200":
//	   $ref: '#/responses/TemplateList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplate) GetList(ctx *gin.Context) {

	total, items, err := c.endpoint.Template.GetList()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Template, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(999, 0, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /template/{name} templateDeleteByName
// ---
// parameters:
// - description: Template Name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: delete template by string
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
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
func (c ControllerTemplate) Delete(ctx *gin.Context) {

	name := ctx.Param("name")
	if name == "" {
		NewError(400, "bad param name").Send(ctx)
		return
	}

	if err := c.endpoint.Template.Delete(name); err != nil {
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

// swagger:operation GET /templates/search templateSearch
// ---
// summary: search template
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
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
//	   $ref: '#/responses/TemplateSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplate) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Template.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Template, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("templates", result)
	resp.Send(ctx)
}

// swagger:operation POST /templates/preview templatePreview
// ---
// summary: preview template
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template
// parameters:
// - description: Update item params
//   in: body
//   name: template
//   required: true
//   schema:
//     $ref: '#/definitions/TemplateContent'
//     type: object
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplate) Preview(ctx *gin.Context) {

	params := &models.TemplateContent{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if len(params.Items) == 0 {
		return
	}

	templatePreview := &m.TemplateContent{}
	common.Copy(&templatePreview, &params)

	data, err := c.endpoint.Template.Preview(templatePreview)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	ctx.String(200, data)
}
