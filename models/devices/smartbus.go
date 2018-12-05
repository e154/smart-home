package devices

import "github.com/e154/smart-home/system/validation"

type DevSmartBusConfig struct {
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required"`
	Sleep    int `json:"sleep"`
}

func (d DevSmartBusConfig) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type DevSmartBusRequest struct {
	Command []byte `json:"command"`
}
type DevSmartBusResponse struct {
	BaseResponse
	Result []byte `json:"result"`
}
