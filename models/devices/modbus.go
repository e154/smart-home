package devices

import (
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
)

const (
	DevTypeModbus = DeviceType("modbus")

	// bit access
	ReadCoils          = "ReadCoils"
	ReadDiscreteInputs = "ReadDiscreteInputs"
	WriteSingleCoil    = "WriteSingleCoil"
	WriteMultipleCoils = "WriteMultipleCoils"

	// 16-bit access
	ReadInputRegisters         = "ReadInputRegisters"
	ReadHoldingRegisters       = "ReadHoldingRegisters"
	ReadWriteMultipleRegisters = "ReadWriteMultipleRegisters"
	WriteSingleRegister        = "WriteSingleRegister"
	WriteMultipleRegisters     = "WriteMultipleRegisters"
)

type DevModBusConfig struct {
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`   // 1-32
	Baud     int    `json:"baud"`                               // 9600, 19200, ...
	DataBits int    `json:"data_bits" mapstructure:"data_bits"` // 5-9
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"` // 1, 2
	Parity   string `json:"parity"`                             // none, odd, even
	Timeout  int    `json:"timeout"`                            // milliseconds
}

func (d DevModBusConfig) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type DevModBusRequest struct {
	Function string   `json:"function"`
	Address  uint16   `json:"address"`
	Count    uint16   `json:"count"`
	Command  []uint16 `json:"command"`
}

// params:
// result
// error
// time
type DevModBusResponse struct {
	BaseResponse
	Result []uint16 `json:"result"`
}
