package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeSmartBus = DeviceType("smartbus")
)

type DevSmartBusConfig struct {
	Validation
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required" mapstructure:"stop_bits"`
	Sleep    int `json:"sleep"`
}

type DevSmartBusRequest struct {
	Command []byte `json:"command"`
}
type DevSmartBusResponse struct {
	BaseResponse
	Result []byte `json:"result"`
}
