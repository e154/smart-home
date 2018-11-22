package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

func addScripts(adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (script1, script2, script3 *m.Script) {

	// add script
	// ------------------------------------------------
	script1 = &m.Script{
		Lang:        "coffeescript",
		Name:        "script1",
		Source:      coffeescript1,
		Description: "test1",
	}
	ok, _ := script1.Valid()
	So(ok, ShouldEqual, true)

	engine1, err := scriptService.NewEngine(script1)
	So(err, ShouldBeNil)
	err = engine1.Compile()
	So(err, ShouldBeNil)
	script1Id, err := adaptors.Script.Add(script1)
	So(err, ShouldBeNil)
	script1, err = adaptors.Script.GetById(script1Id)
	So(err, ShouldBeNil)

	script2 = &m.Script{
		Lang:        "coffeescript",
		Name:        "script2",
		Source:      coffeescript2,
		Description: "script2",
	}
	ok, _ = script2.Valid()
	So(ok, ShouldEqual, true)

	engine2, err := scriptService.NewEngine(script2)
	So(err, ShouldBeNil)
	err = engine2.Compile()
	So(err, ShouldBeNil)
	script2Id, err := adaptors.Script.Add(script2)
	So(err, ShouldBeNil)
	script2, err = adaptors.Script.GetById(script2Id)
	So(err, ShouldBeNil)

	script3 = &m.Script{
		Lang:        "coffeescript",
		Name:        "script3",
		Source:      coffeescript3,
		Description: "script3",
	}
	ok, _ = script3.Valid()
	So(ok, ShouldEqual, true)

	engine3, err := scriptService.NewEngine(script3)
	So(err, ShouldBeNil)
	err = engine3.Compile()
	So(err, ShouldBeNil)
	script3Id, err := adaptors.Script.Add(script3)
	So(err, ShouldBeNil)
	script3, err = adaptors.Script.GetById(script3Id)
	So(err, ShouldBeNil)

	return
}

const coffeescript1 = `

`

const coffeescript2 = `

`

const coffeescript3 = `

`