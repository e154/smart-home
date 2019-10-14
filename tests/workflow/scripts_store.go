package workflow

import "github.com/e154/smart-home/system/scripts"

var store func(interface{})

func storeRegisterCallback(scriptService *scripts.ScriptService) {
	scriptService.PushFunctions("store", func(value interface{}) {
		store(value)
	})
}
