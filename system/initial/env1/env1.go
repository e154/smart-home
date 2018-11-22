package env1

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

func Init(adaptors *adaptors.Adaptors, core *core.Core, scriptService *scripts.ScriptService) {

	// nodes
	// ------------------------------------------------
	node1, _ := nodes(adaptors)

	// scripts
	// ------------------------------------------------
	script1, script2, script3 := addScripts(adaptors, scriptService)

	// devices
	// ------------------------------------------------
	devices(node1, adaptors, core, script1, script2, script3)

}
