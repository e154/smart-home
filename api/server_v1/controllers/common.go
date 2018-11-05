package controllers

import (
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
)

var (
	log = logging.MustGetLogger("controllers")
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
}

func NewControllerCommon(adaptors *adaptors.Adaptors) *ControllerCommon {
	return &ControllerCommon{adaptors: adaptors}
}
