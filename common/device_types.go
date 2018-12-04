package common

import "github.com/e154/smart-home/system/validation"

type DevConfSmartBus struct {
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required"`
	Sleep    int `json:"sleep"`
}

func (d DevConfSmartBus) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type DevConfModBus struct {
}

type DevConfCommand struct {
}