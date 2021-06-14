// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

// ControllerTemplateItem ...
type ControllerTemplateItem struct {
	*ControllerCommon
}

// NewControllerTemplateItem ...
func NewControllerTemplateItem(common *ControllerCommon) *ControllerTemplateItem {
	return &ControllerTemplateItem{ControllerCommon: common}
}

// swagger:operation POST /template_item templateAddItem
// ---
// parameters:
// - description: template params
//   in: body
//   name: template
//   required: true
//   schema:
//     $ref: '#/definitions/NewTemplateItem'
//     type: object
// summary: add new template item
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
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
func (c ControllerTemplateItem) Add(ctx *gin.Context) {

	params := &models.NewTemplateItem{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	template := &m.Template{}
	_ = common.Copy(&template, &params, common.JsonEngine)
	template.Type = m.TemplateTypeItem

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

// swagger:operation GET /template_item/{name} templateGetItemByName
// ---
// parameters:
// - description: TemplateItem Name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get template item by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/TemplateItem'
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
func (c ControllerTemplateItem) GetByName(ctx *gin.Context) {

	name := ctx.Param("name")
	if name == "" {
		NewError(400, "bad param name").Send(ctx)
		return
	}

	item, err := c.endpoint.Template.GetItemByName(name)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.TemplateItem{}
	_ = common.Copy(&result, &item, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /template_item/{name} templateUpdateItemByName
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
//     $ref: '#/definitions/UpdateTemplateItem'
//     type: object
// summary: update template by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
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
func (c ControllerTemplateItem) Update(ctx *gin.Context) {

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

// swagger:operation PUT /template_item/status/{name} templateUpdateStatusItemByName
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
//     $ref: '#/definitions/UpdateTemplateItemStatus'
//     type: object
// summary: update template by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
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
func (c ControllerTemplateItem) UpdateStatus(ctx *gin.Context) {

	name := ctx.Param("name")
	if name == "" {
		NewError(400, "bad param name").Send(ctx)
		return
	}

	params := &models.UpdateTemplateItemStatus{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	template := &m.Template{}
	_ = common.Copy(&template, &params, common.JsonEngine)

	errs, err := c.endpoint.Template.UpdateStatus(template)
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

// swagger:operation GET /template_items templateGetItemList
// ---
// summary: get template item list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
// responses:
//   "200":
//	   $ref: '#/responses/TemplateItemSortedList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplateItem) GetList(ctx *gin.Context) {

	total, items, err := c.endpoint.Template.GetItemsSortedList()
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(map[string]interface{}{
		"items": items,
		"total": total,
	}).Send(ctx)
	return
}

// swagger:operation DELETE /template_item/{name} templateDeleteItemByName
// ---
// parameters:
// - description: TemplateItem Name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: delete template item by string
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
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
func (c ControllerTemplateItem) Delete(ctx *gin.Context) {

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

// swagger:operation GET /template_items/tree templateGetItemsTree
// ---
// parameters:
// summary: get template items tree
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/TemplateTree'
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
func (c ControllerTemplateItem) GetTree(ctx *gin.Context) {

	tree, err := c.endpoint.Template.GetItemsTree()
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.TemplateTree{}
	_ = common.Copy(&result, &tree, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /template_items/tree templateUpdateItemsTree
// ---
// parameters:
// - description: Update item params
//   in: body
//   name: template
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateTemplateTree'
//     type: object
// summary: update template by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - template_item
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
func (c ControllerTemplateItem) UpdateTree(ctx *gin.Context) {

	params := make(models.UpdateTemplateTree, 0)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	tree := make([]*m.TemplateTree, 0, len(params))
	_ = common.Copy(&tree, &params, common.JsonEngine)

	err := c.endpoint.Template.UpdateItemsTree(tree)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}
