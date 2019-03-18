package controllers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"github.com/e154/smart-home/api/server/v1/models"
)

type ControllerFlow struct {
	*ControllerCommon
}

func NewControllerFlow(common *ControllerCommon) *ControllerFlow {
	return &ControllerFlow{ControllerCommon: common}
}

// swagger:operation POST /flow flowAdd
// ---
// parameters:
// - description: flow params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/NewFlow'
//     type: object
// summary: add new flow
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//	   $ref: '#/responses/NewObjectSuccess'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) Add(ctx *gin.Context) {

	flow := &models.NewFlow{}
	if err := ctx.ShouldBindJSON(&flow); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	id, errs, err := AddFlow(flow, c.adaptors, c.core)
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

// swagger:operation GET /flow/{id} flowGetById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Flow'
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
func (c ControllerFlow) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	flow, err := GetFlowById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(flow).Send(ctx)
}

// swagger:operation GET /flow/{id}/redactor flowGetRedactor
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get flow redactor data by flow id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/RedactorFlow'
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
func (c ControllerFlow) GetRedactor(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	flow, err := GetFlowRedactor(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(flow).Send(ctx)
}

// swagger:operation PUT /flow/{id}/redactor flowUpdateRedactor
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: flow redactor params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/RedactorFlow'
//     type: object
// summary: update flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/RedactorFlow'
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
func (c ControllerFlow) UpdateRedactor(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	redactor := &models.RedactorFlow{}
	if err := ctx.ShouldBindJSON(&redactor); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	redactor.Id = int64(aid)

	result, errs, err := UpdateFlowRedactor(redactor, c.adaptors, c.core)
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

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /flow/{id} flowUpdateById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update flow params
//   in: body
//   name: flow
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateFlow'
//     type: object
// summary: update flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Flow'
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
func (c ControllerFlow) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	flow := &models.UpdateFlow{}
	if err := ctx.ShouldBindJSON(&flow); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	flow.Id = int64(aid)

	result, errs, err := UpdateFlow(flow, c.adaptors, c.core)
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

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /flows flowList
// ---
// summary: get flow list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
//	   $ref: '#/responses/FlowList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetFlowList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
}

// swagger:operation DELETE /flow/{id} flowDeleteById
// ---
// parameters:
// - description: Flow ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete flow by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
func (c ControllerFlow) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err = DeleteFlowById(int64(aid), c.adaptors, c.core); err != nil {
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

// swagger:operation GET /flows/search flowSearch
// ---
// summary: search flow
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - flow
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
//	   $ref: '#/responses/FlowSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerFlow) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	flows, _, err := SearchFlow(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("flows", flows)
	resp.Send(ctx)
}
