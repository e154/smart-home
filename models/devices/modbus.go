package devices

import (
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
)

const (
	DevTypeModbus = DeviceType("modbus")
)

type DevModBusConfig struct {
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`
	Baud     int    `json:"baud"`
	DataBits int    `json:"data_bits" mapstructure:"data_bits"`
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"`
	Parity   string `json:"parity"` // none, odd, even
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
