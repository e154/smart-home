package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
)

type ControllerIndex struct {
	*ControllerCommon
}

func NewControllerIndex(common *ControllerCommon) *ControllerIndex {
	return &ControllerIndex{ControllerCommon: common}
}

// Index godoc
// @tags index
// @Summary index page
// @Description
// @Produce text/plain
// @Accept  text/plain
// @Success 200
// @Router / [get]
func (i ControllerIndex) Index(c *gin.Context) {
	apiVersion := Index()
	c.String(http.StatusOK, apiVersion)
	return
}
