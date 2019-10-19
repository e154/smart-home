package workflow

import "github.com/e154/smart-home/system/scripts"

var store = func(interface{}) {}
var store2 = func(interface{}) {}

func storeRegisterCallback(scriptService *scripts.ScriptService) {
	scriptService.PushFunctions("store", func(value interface{}) {
		store(value)
	})
	scriptService.PushFunctions("store2", func(value interface{}) {
		store2(value)
	})
}
