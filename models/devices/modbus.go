package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeModbusRtu = DeviceType("modbus_rtu")
	DevTypeModbusTcp = DeviceType("modbus_tcp")

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

type DevModBusRtuConfig struct {
	Validation
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`   // 1-32
	Baud     int    `json:"baud"`                               // 9600, 19200, ...
	DataBits int    `json:"data_bits" mapstructure:"data_bits"` // 5-9
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"` // 1, 2
	Parity   string `json:"parity"`                             // none, odd, even
	Timeout  int    `json:"timeout"`                            // milliseconds
}

type DevModBusTcpConfig struct {
	Validation
	AddressPort string `json:"address_port" mapstructure:"address_port"`
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
