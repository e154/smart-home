package bind

import (
	"github.com/e154/smart-home/api/scripts"
)

func init() {
	scripts.PushStruct("Notifr", &NotifrBind{})
	scripts.PushStruct("Log", &LogBind{})
	scripts.PushFunctions("Execute", Execute)
}
