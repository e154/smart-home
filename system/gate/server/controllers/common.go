package controllers

import (
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	ApiFullAddress string
	Mode string
}

// NewControllerCommon ...
func NewControllerCommon(apiFullAddress, mode string) *ControllerCommon {
	return &ControllerCommon{
		ApiFullAddress: apiFullAddress,
		Mode: mode,
	}
}
