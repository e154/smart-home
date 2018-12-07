package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeModbus  = DeviceType("modbus")
)

type DevModBusConfig struct {
}
type DevModBusRequest struct {
}
type DevModBusResponse struct {
	BaseResponse
}
