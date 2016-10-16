package serial_test

import (
	"testing"
	"../serial"
)

const expect = 245

func TestLRC(t *testing.T) {
	lrc := serial.LRC([]byte{2,1,0,0,0,8})
	if lrc != expect {
		t.Fatalf("lrc expected %v, actual %v", expect, lrc)
	}
}