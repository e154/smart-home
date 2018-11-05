package controllers

import (
	"net/http"
	"io"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/index"
)

type ControllerIndex struct {
	*ControllerCommon
}

func NewControllerIndex(common *ControllerCommon) *ControllerIndex {
	return &ControllerIndex{ControllerCommon: common}
}

func (c ControllerIndex) ControllerIndex(params index.IndexParams) middleware.Responder {
	return NewResponse(func(rw http.ResponseWriter, producer runtime.Producer) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(rw, "Server API V0.1: OK")
	})
}
