package env1

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/access_list"
)

// env1
//
// node1 + node2
// script1 + script2 + script3
// device1
// 		+ child device2
// device3
//
func Init(adaptors *adaptors.Adaptors,
	accessList *access_list.AccessListService,
	scriptService *scripts.ScriptService) {

	// roles
	// ------------------------------------------------
	roles(adaptors, accessList)

	// nodes
	// ------------------------------------------------
	node1, _ := nodes(adaptors)

	// scripts
	// ------------------------------------------------
	script1, script2, script3, script4, script5, script6, script7 := addScripts(adaptors, scriptService)

	// devices
	// ------------------------------------------------
	_, _, _, deviceAction3, deviceAction7 := devices(node1, adaptors, script1, script2, script3, script7)

	// workflow
	// ------------------------------------------------
	addWorkflow(adaptors, deviceAction3, deviceAction7, script4, script5, script6)
}
