package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
)

type ControllerTemplate struct {
	*ControllerCommon
}

func NewControllerTemplate(common *ControllerCommon) *ControllerTemplate {
	return &ControllerTemplate{ControllerCommon: common}
}

// swagger:operation POST /template templateAddItem
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

// swagger:operation GET /template/{name} templateGetItemByName
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

	item, err := c.endpoint.Template.GetByName(name)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Template{}
	_ = common.Copy(&result, &item, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /template/{name} templateUpdateItemByName
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
//	   $ref: '#/responses/TemplateList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerTemplate) GetList(ctx *gin.Context) {

	//_, sortBy, order, limit, offset := c.list(ctx)
	//items, total, err := c.endpoint.Template.GetList(int64(limit), int64(offset), order, sortBy)
	//if err != nil {
	//	NewError(500, err).Send(ctx)
	//	return
	//}
	//
	//result := make([]*models.Template, 0)
	//common.Copy(&result, &items)
	//
	//resp := NewSuccess()
	//resp.Page(limit, offset, total, result).Send(ctx)
	//return
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

	//query, limit, offset := c.select2(ctx)
	//items, _, err := c.endpoint.Template.Search(query, limit, offset)
	//if err != nil {
	//	NewError(500, err).Send(ctx)
	//	return
	//}
	//
	//result := make([]*models.Template, 0)
	//common.Copy(&result, &items)
	//
	//resp := NewSuccess()
	//resp.Item("templates", result)
	//resp.Send(ctx)
}
