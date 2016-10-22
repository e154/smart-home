package serial

import (
	"fmt"
	"time"
)

func Run() {

	usb1 := &Serial{Dev:"/dev/ttyUSB1", Baud:19200, ReadTimeout: time.Second * 2, StopBits:2}
	//usb1 := &Serial2{config: serial.Config{Address: "/dev/ttyUSB1", BaudRate: 9600, StopBits: 2, Timeout: time.Second * 2, Parity: "N", DataBits: 8}}
	//err := usb1.Connect()
	//if(err != nil) {
	//	fmt.Println(err.Error())
	//}

	modbus := &Modbus{serial: usb1}
	res, err := modbus.Send([]byte{0x01,0x03,0x00,0x00,0x00,0x05})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%d\r\n", res)

	//var max_temp int = 1089
	//hi := (max_temp >> 4) & 0xFF
	//low := max_temp & 0xFF
	//
	//fmt.Printf("%d - %d\r\n", hi, low)
	//
	//temp := hi << 4
	//temp |= low & 0xFF
	//
	//fmt.Printf("%d\r\n", temp)
}

