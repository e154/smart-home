package env1

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

// env1
//
// node1 + node2
// script1 + script2 + script3
// device1
// 		+ child device2
// device3
//
func Init(adaptors *adaptors.Adaptors, core *core.Core, scriptService *scripts.ScriptService) {

	// nodes
	// ------------------------------------------------
	node1, _ := nodes(adaptors)

	// scripts
	// ------------------------------------------------
	script1, script2, script3, script4, script5, script6 := addScripts(adaptors, scriptService)

	// devices
	// ------------------------------------------------
	_, _, _, deviceAction1 := devices(node1, adaptors, script1, script2, script3)

	// workflow
	// ------------------------------------------------
	addWorkflow(adaptors, deviceAction1, script4, script5, script6)
}
