package devices

import (
	. "github.com/e154/smart-home/common"
)

type BaseResponse struct {
	Error     string  `json:"error"`
	Time      float64 `json:"time"`
}

const (
	DevTypeDefault  = DeviceType("default")
)
