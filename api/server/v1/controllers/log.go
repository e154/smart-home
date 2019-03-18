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

// swagger:operation POST /log logAdd
// ---
// parameters:
// - description: log params
//   in: body
//   name: log
//   required: true
//   schema:
//     $ref: '#/definitions/NewLog'
//     type: object
// summary: add new log
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - log
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Log'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerLog) Add(ctx *gin.Context) {

	params := &models.NewLog{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	result, errs, err := AddLog(params, c.adaptors)
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
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /log/{id} logGetById
// ---
// parameters:
// - description: Log ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get log by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - log
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Log'
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
func (c ControllerLog) GetById(ctx *gin.Context) {

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
	resp.SetData(log).Send(ctx)
}

// swagger:operation GET /logs logList
// ---
// summary: get log list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - log
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
// - default: id
//   description: query
//   in: query
//   name: query
//   type: string
// responses:
//   "200":
//	   $ref: '#/responses/LogList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerLog) GetList(ctx *gin.Context) {

	query, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetLogList(int64(limit), int64(offset), order, sortBy, query, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
}

// swagger:operation DELETE /log/{id} logDeleteById
// ---
// parameters:
// - description: Log ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete log by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - log
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
func (c ControllerLog) Delete(ctx *gin.Context) {

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

// swagger:operation GET /logs/search logSearch
// ---
// summary: search log
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - log
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
//	   $ref: '#/responses/LogSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
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
