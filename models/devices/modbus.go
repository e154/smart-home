package devices

import (
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
)

const (
	DevTypeModbus = DeviceType("modbus")
)

type DevModBusConfig struct {
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required"`
	Sleep    int `json:"sleep"`
}

func (d DevModBusConfig) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type DevModBusRequest struct {
	Command []byte `json:"command"`
}

type DevModBusResponse struct {
	BaseResponse
	Result []byte `json:"result"`
}
