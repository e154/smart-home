package serial

import (
	"io/ioutil"
	"strings"
)

var (
	Devices []string
)

func FindSerials() {
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyS") ||
			strings.Contains(f.Name(), "ttyUSB") {
			Devices = append(Devices, "/dev/" + f.Name())
		}
	}

}

func HI(b byte) (byte) {
	return (b >> 8) & 0xFF
}

func LOW(b byte) (byte) {
	return b & 0x0F
}
