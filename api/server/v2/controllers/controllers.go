package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
)

type ControllersV2 struct {
	Index   *ControllerIndex
	Swagger *ControllerSwagger
}

func NewControllersV2(adaptors *adaptors.Adaptors,
	core *core.Core) *ControllersV2 {
	common := NewControllerCommon(adaptors, core)
	return &ControllersV2{
		Index:   NewControllerIndex(common),
		Swagger: NewControllerSwagger(common),
	}
}
