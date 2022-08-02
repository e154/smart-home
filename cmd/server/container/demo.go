package container

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/initial/demo"
	"github.com/e154/smart-home/system/initial/demo/example1"
	"github.com/e154/smart-home/system/scripts"
)

func NewDemo(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService) (d *demo.Demos) {
	list := make(map[string]demo.Demo)
	list["example1"] = example1.NewExample1(adaptors, scriptService)
	d = demo.NewDemos(list)
	return
}
