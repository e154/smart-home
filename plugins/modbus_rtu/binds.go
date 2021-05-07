package modbus_rtu

import (
	"fmt"
	"github.com/e154/smart-home/system/event_bus"
)

type modbusRtu func(f string, address, count uint16, command []uint16) (result ModBusResponse)

func NewModbusRtu(eventBus event_bus.EventBus) (modbus modbusRtu) {

	return func(f string, address, count uint16, command []uint16) (result ModBusResponse) {
		fmt.Printf("send message^ func(%s), %d, %d, %b \n", f, address, count, command)
		return
	}
}
