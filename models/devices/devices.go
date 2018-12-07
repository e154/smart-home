package devices

import (
	. "github.com/e154/smart-home/common"
)

type BaseResponse struct {
	Error     string  `json:"error"`
	ErrorCode string  `json:"error_code"`
	Time      float64 `json:"time"`
}

const (
	DevTypeDefault  = DeviceType("default")
)
