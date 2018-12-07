package devices

import (
	"github.com/e154/smart-home/system/validation"
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeCommand  = DeviceType("command")
)

type DevCommandConfig struct {
}

func (d DevCommandConfig) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type DevCommandRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}


type DevCommandResponse struct {
	BaseResponse
	Result string `json:"result"`
}
