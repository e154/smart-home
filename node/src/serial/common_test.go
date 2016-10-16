package serial_test

import (
	"testing"
	"../serial"
)

func TestFindSerials(t *testing.T) {

	old_val := len(serial.Devices)

	serial.FindSerials()

	if old_val == len(serial.Devices) {
		t.Errorf("Device list no updated")
	}
}