package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerIndex struct {
	*ControllerCommon
}

func NewControllerIndex(common *ControllerCommon) *ControllerIndex {
	return &ControllerIndex{ControllerCommon: common}
}

// swagger:operation GET / index
// ---
// summary: index page
// description:
// consumes:
// - text/plain
// produces:
// - text/plain
// tags:
// - index
// responses:
//   "200":
//	   description: Success response
func (i ControllerIndex) Index(c *gin.Context) {
	apiVersion := "Server API V1: OK"
	c.String(http.StatusOK, apiVersion)
	return
}
