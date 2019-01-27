package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
	"fmt"
)

func addScripts(adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (scripts map[string]*m.Script) {

	scripts = make(map[string]*m.Script)

	// mb_condition_check_v1
	// ------------------------------------------------
	script1 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_condition_check_v1",
		Source:      MbDev1ConditionCheckV1,
		Description: "condition check",
	}
	ok, _ := script1.Valid()
	So(ok, ShouldEqual, true)

	engine1, err := scriptService.NewEngine(script1)
	So(err, ShouldBeNil)
	err = engine1.Compile()
	if err != nil {
		fmt.Println(err.Error())
	}
	So(err, ShouldBeNil)
	script1Id, err := adaptors.Script.Add(script1)
	So(err, ShouldBeNil)
	script1, err = adaptors.Script.GetById(script1Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_condition_check_v1"] = script1

	// mb_turn_on_light1_v1
	// ------------------------------------------------
	script2 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_on_light1_v1",
		Source:      MbDev1TurnOnLight1V1,
		Description: "turn on light1",
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

	scripts["mb_dev1_turn_on_light1_v1"] = script2

	// mb_turn_off_light1_v1
	// ------------------------------------------------
	script3 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_off_light1_v1",
		Source:      MbDev1TurnOffLight1V1,
		Description: "turn off light1",
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

	scripts["mb_dev1_turn_off_light1_v1"] = script3

	// mb_turn_on_light1_v1
	// ------------------------------------------------
	script8 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_on_light2_v1",
		Source:      MbDev1TurnOnLight2V1,
		Description: "turn on light2",
	}
	ok, _ = script8.Valid()
	So(ok, ShouldEqual, true)

	engine8, err := scriptService.NewEngine(script8)
	So(err, ShouldBeNil)
	err = engine8.Compile()
	So(err, ShouldBeNil)
	script8Id, err := adaptors.Script.Add(script8)
	So(err, ShouldBeNil)
	script8, err = adaptors.Script.GetById(script8Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_on_light2_v1"] = script8

	// mb_turn_off_light1_v1
	// ------------------------------------------------
	script9 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_off_light2_v1",
		Source:      MbDev1TurnOffLight2V1,
		Description: "turn off light2",
	}
	ok, _ = script9.Valid()
	So(ok, ShouldEqual, true)

	engine9, err := scriptService.NewEngine(script9)
	So(err, ShouldBeNil)
	err = engine9.Compile()
	So(err, ShouldBeNil)
	script9Id, err := adaptors.Script.Add(script9)
	So(err, ShouldBeNil)
	script9, err = adaptors.Script.GetById(script9Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_off_light2_v1"] = script9

	// cmd_condition_check_v1
	// ------------------------------------------------
	script4 := &m.Script{
		Lang:        "coffeescript",
		Name:        "cmd_condition_check_v1",
		Source:      CmdConditionCheckV1,
		Description: "condition check",
	}
	ok, _ = script4.Valid()
	So(ok, ShouldEqual, true)

	engine4, err := scriptService.NewEngine(script4)
	So(err, ShouldBeNil)
	err = engine4.Compile()
	So(err, ShouldBeNil)
	script4Id, err := adaptors.Script.Add(script4)
	So(err, ShouldBeNil)
	script4, err = adaptors.Script.GetById(script4Id)
	So(err, ShouldBeNil)

	scripts["cmd_condition_check_v1"] = script4

	// wflow_scenario_weekday_v1
	// ------------------------------------------------
	script5 := &m.Script{
		Lang:        "coffeescript",
		Name:        "wflow_scenario_weekday_v1",
		Source:      WflowScenarioWeekdayV1,
		Description: "weekday scenario",
	}
	ok, _ = script5.Valid()
	So(ok, ShouldEqual, true)

	engine5, err := scriptService.NewEngine(script5)
	So(err, ShouldBeNil)
	err = engine5.Compile()
	So(err, ShouldBeNil)
	script5Id, err := adaptors.Script.Add(script5)
	So(err, ShouldBeNil)
	script5, err = adaptors.Script.GetById(script5Id)
	So(err, ShouldBeNil)

	scripts["wflow_scenario_weekday_v1"] = script5

	// wflow_scenario_weekend_v1
	// ------------------------------------------------
	script6 := &m.Script{
		Lang:        "coffeescript",
		Name:        "wflow_scenario_weekend_v1",
		Source:      WflowScenarioWeekendV1,
		Description: "weekend scenario",
	}
	ok, _ = script6.Valid()
	So(ok, ShouldEqual, true)

	engine6, err := scriptService.NewEngine(script6)
	So(err, ShouldBeNil)
	err = engine6.Compile()
	So(err, ShouldBeNil)
	script6Id, err := adaptors.Script.Add(script6)
	So(err, ShouldBeNil)
	script6, err = adaptors.Script.GetById(script6Id)
	So(err, ShouldBeNil)

	scripts["wflow_scenario_weekend_v1"] = script6

	// base_script
	// ------------------------------------------------
	script7 := &m.Script{
		Lang:        "coffeescript",
		Name:        "base_script",
		Source:      BaseScript,
		Description: "weekend scenario",
	}
	ok, _ = script7.Valid()
	So(ok, ShouldEqual, true)

	engine7, err := scriptService.NewEngine(script7)
	So(err, ShouldBeNil)
	err = engine7.Compile()
	So(err, ShouldBeNil)
	script7Id, err := adaptors.Script.Add(script7)
	So(err, ShouldBeNil)
	script7, err = adaptors.Script.GetById(script7Id)
	So(err, ShouldBeNil)

	scripts["base_script"] = script6

	return
}

const MbDev1ConditionCheckV1 = `
# get device status
fetchStatus =->

    COMMAND = []
    FUNC = 'ReadHoldingRegisters'
    ADDRESS = 0
    COUNT = 16
    
    res = device.modBus FUNC, ADDRESS, COUNT, COMMAND
    if res.error
        print 'error: ', res.error
        return
    else
        print 'ok: ', res.result

    #if res.result[0] == 1
    #    device.modBus 'WriteMultipleRegisters', 0, 1, [0]
    #else
    #    device.modBus 'WriteMultipleRegisters', 0, 1, [1]

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOnLight1V1 = `
# turn on first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 0, 1, [1]
    if res.error
        print 'error: ', res.error
    else
        print 'ok: ', res.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOffLight1V1 = `
# turn off first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 0, 1, [0]
    if res.error
        print 'error: ', res.error
    else
        print 'ok: ', res.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOnLight2V1 = `
# turn on first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 1, 1, [1]
    if res.error
        print 'error: ', res.error
    else
        print 'ok: ', res.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOffLight2V1 = `
# turn off first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 1, 1, [0]
    if res.error
        print 'error: ', res.error
    else
        print 'ok: ', res.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`

const CmdConditionCheckV1 = `
main =->
    
    NAME = './data/scripts/ping.sh'
    ARGS = ['ya.ru']
    result = device.runCommand NAME, ARGS
    print result

main()
`

const WflowScenarioWeekdayV1 = `
# variables:

#IC.Workflow()
#    .getName()
#    .getDescription()
#    .setVar(string, interface)
#    .getVar(string)
#    .getScenario() string
#    .setScenario(string)

#
# workflow script example
#

on_enter =->
    print 'on enter'

on_exit =->
    print 'on exit'
    
main =->
    scenario = 'weekday'
    print 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario
    
`

const WflowScenarioWeekendV1 = `
# variables:

#IC.Workflow()
#    .getName()
#    .getDescription()
#    .setVar(string, interface)
#    .getVar(string)
#    .getScenario() string
#    .setScenario(string)

#
# workflow script example
#

on_enter =->
    print 'on enter'

on_exit =->
    print 'on exit'
    
main =->
    scenario = 'weekend'
    print 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario
`

const BaseScript = `

main =->
`