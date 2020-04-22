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
	"strconv"
	"time"
)

// ControllerMetric ...
type ControllerMetric struct {
	*ControllerCommon
}

// NewControllerMetric ...
func NewControllerMetric(common *ControllerCommon) *ControllerMetric {
	return &ControllerMetric{ControllerCommon: common}
}

// swagger:operation POST /metric metricAdd
// ---
// parameters:
// - description: metric params
//   in: body
//   name: metric
//   required: true
//   schema:
//     $ref: '#/definitions/NewMetric'
//     type: object
// summary: add new metric
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Metric'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMetric) Add(ctx *gin.Context) {

	params := &models.NewMetric{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	metric := m.Metric{}
	common.Copy(&metric, &params, common.JsonEngine)

	metric, errs, err := c.endpoint.Metric.Add(metric)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Metric{}
	if err = common.Copy(&result, &metric, common.JsonEngine); err != nil {
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /metric/{id} metricGetById
// ---
// summary: get metric
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
// parameters:
// - description: Metric ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - name: range
//   description: 1h, 6h, 12h, 24h, 7d, 30d
//   in: query
//   type: string
// - name: from
//   in: query
//   description: from (2020-04-16T00:00:00Z)
//   type: string
//   format: date-time
// - name: "to"
//   in: query
//   description: from (2020-04-16T00:00:00Z)
//   type: string
//   format: date-time
// responses:
//   "200":
//	   $ref: '#/definitions/Metric'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMetric) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var metricRange *string
	var from, to *time.Time

	if ctx.Request.URL.Query().Get("from") != "" {
		if _from, err := time.Parse(time.RFC3339, ctx.Request.URL.Query().Get("from")); err != nil {
			NewError(403, err).Send(ctx)
			return
		} else {
			from = common.Time(_from)
		}
	}

	if ctx.Request.URL.Query().Get("to") != "" {
		if _to, err := time.Parse(time.RFC3339, ctx.Request.URL.Query().Get("to")); err != nil {
			NewError(403, err).Send(ctx)
			return
		} else {
			to = common.Time(_to)
		}
	}

	if ctx.Request.URL.Query().Get("range") != "" {
		metricRange = common.String(ctx.Request.URL.Query().Get("range"))
	}

	metric, err := c.endpoint.Metric.GetById(int64(aid), from, to, metricRange)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Metric{}
	common.Copy(&result, &metric, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /metric/{id} metricUpdateById
// ---
// parameters:
// - description: Metric ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update metric params
//   in: body
//   name: metric
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateMetric'
//     type: object
// summary: update metric by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Metric'
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
func (c ControllerMetric) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := models.UpdateMetric{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	metric := m.Metric{}
	common.Copy(&metric, &params, common.JsonEngine)

	metric, errs, err := c.endpoint.Metric.Update(metric)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := models.Metric{}
	common.Copy(&result, &metric, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /metrics metricList
// ---
// summary: get metric list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
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
//	   $ref: '#/responses/MetricList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMetric) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Metric.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Metric, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
	return
}

// swagger:operation DELETE /metric/{id} metricDeleteById
// ---
// parameters:
// - description: Metric ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete metric by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
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
func (c ControllerMetric) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Metric.Delete(int64(aid)); err != nil {
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

// swagger:operation GET /metrics/search metricSearch
// ---
// summary: search metric
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
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
//	   $ref: '#/responses/MetricSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMetric) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	items, _, err := c.endpoint.Metric.Search(query, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Metric, 0)
	common.Copy(&result, &items)

	resp := NewSuccess()
	resp.Item("metrics", result)
	resp.Send(ctx)
}

// swagger:operation POST /metric/data metricAddData
// ---
// parameters:
// - description: metric data
//   in: body
//   name: metric
//   required: true
//   schema:
//     $ref: '#/definitions/NewMetricDataItem'
//     type: object
// summary: add new metric
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - metric
// responses:
//   "200":
//     description: OK
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMetric) AddData(ctx *gin.Context) {

}
