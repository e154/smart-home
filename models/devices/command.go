package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeCommand = DeviceType("command")
)

type DevCommandConfig struct {
	Validation
}

type DevCommandRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type DevCommandResponse struct {
	BaseResponse
	Result string `json:"result"`
}
