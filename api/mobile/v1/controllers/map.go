package controllers

import (
	"github.com/e154/smart-home/api/mobile/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
)

type ControllerMap struct {
	*ControllerCommon
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{ControllerCommon: common}
}

// swagger:operation GET /map/active_elements mapGetActiveElements
// ---
// summary: get active map elements
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - map
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
//     description: OK
//     schema:
//       $ref: '#/responses/MapActiveElementList'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMap) GetActiveElements(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Map.GetActiveElements(sortBy, order, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MapElement, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}
