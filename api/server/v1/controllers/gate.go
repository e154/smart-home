package controllers

import (
	"github.com/gin-gonic/gin"
)

type ControllerGate struct {
	*ControllerCommon
}

func NewControllerGate(common *ControllerCommon) *ControllerGate {
	return &ControllerGate{ControllerCommon: common}
}

// swagger:operation GET /gate/{id} gateGetList
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
//       $ref: '#/definitions/Gate'
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
func (c ControllerGate) GetSettings(ctx *gin.Context) {


}

func (c ControllerGate) UpdateSettings(ctx *gin.Context) {


}