package env1

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/system/scripts"
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
	scriptService *ScriptService) {

	// images
	// ------------------------------------------------
	imageList := images(adaptors)

	// roles
	// ------------------------------------------------
	roles(adaptors, accessList)

	// nodes
	// ------------------------------------------------
	node1, _ := nodes(adaptors)

	// scripts
	// ------------------------------------------------
	scripts := addScripts(adaptors, scriptService)

	// devices
	// ------------------------------------------------
	devList, deviceActions, deviceStates := devices(node1, adaptors, scripts)

	// workflow
	// ------------------------------------------------
	addWorkflow(adaptors, deviceActions, scripts)

	// maps
	// ------------------------------------------------
	addMaps(adaptors, scripts, devList, imageList, deviceActions, deviceStates)
}
