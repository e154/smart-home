package env1

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

func addScripts(adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (scripts map[string]*m.Script) {

	scripts = make(map[string]*m.Script)

	// device status
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

	// light1
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

	// light2
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

	// light3
	// ------------------------------------------------
	script10 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_on_light3_v1",
		Source:      MbDev1TurnOnLight3V1,
		Description: "turn on light3",
	}
	ok, _ = script10.Valid()
	So(ok, ShouldEqual, true)

	engine10, err := scriptService.NewEngine(script10)
	So(err, ShouldBeNil)
	err = engine10.Compile()
	So(err, ShouldBeNil)
	script10Id, err := adaptors.Script.Add(script10)
	So(err, ShouldBeNil)
	script10, err = adaptors.Script.GetById(script10Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_on_light3_v1"] = script10

	script11 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_off_light3_v1",
		Source:      MbDev1TurnOffLight3V1,
		Description: "turn off light3",
	}
	ok, _ = script11.Valid()
	So(ok, ShouldEqual, true)

	engine11, err := scriptService.NewEngine(script11)
	So(err, ShouldBeNil)
	err = engine11.Compile()
	So(err, ShouldBeNil)
	script11Id, err := adaptors.Script.Add(script11)
	So(err, ShouldBeNil)
	script11, err = adaptors.Script.GetById(script11Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_off_light3_v1"] = script11

	// light4
	// ------------------------------------------------
	script12 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_on_light4_v1",
		Source:      MbDev1TurnOnLight4V1,
		Description: "turn on light4",
	}
	ok, _ = script12.Valid()
	So(ok, ShouldEqual, true)

	engine12, err := scriptService.NewEngine(script12)
	So(err, ShouldBeNil)
	err = engine12.Compile()
	So(err, ShouldBeNil)
	script12Id, err := adaptors.Script.Add(script12)
	So(err, ShouldBeNil)
	script12, err = adaptors.Script.GetById(script12Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_on_light4_v1"] = script12

	script13 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_off_light4_v1",
		Source:      MbDev1TurnOffLight4V1,
		Description: "turn off light4",
	}
	ok, _ = script13.Valid()
	So(ok, ShouldEqual, true)

	engine13, err := scriptService.NewEngine(script13)
	So(err, ShouldBeNil)
	err = engine13.Compile()
	So(err, ShouldBeNil)
	script13Id, err := adaptors.Script.Add(script13)
	So(err, ShouldBeNil)
	script13, err = adaptors.Script.GetById(script13Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_off_light4_v1"] = script13

	// fan1
	// ------------------------------------------------
	script14 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_on_fan1_v1",
		Source:      MbDev1TurnOnFan1V1,
		Description: "turn on fan1",
	}
	ok, _ = script14.Valid()
	So(ok, ShouldEqual, true)

	engine14, err := scriptService.NewEngine(script14)
	So(err, ShouldBeNil)
	err = engine14.Compile()
	So(err, ShouldBeNil)
	script14Id, err := adaptors.Script.Add(script14)
	So(err, ShouldBeNil)
	script14, err = adaptors.Script.GetById(script14Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_on_fan1_v1"] = script14

	script15 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_turn_off_fan1_v1",
		Source:      MbDev1TurnOffFan1V1,
		Description: "turn off fan1",
	}
	ok, _ = script15.Valid()
	So(ok, ShouldEqual, true)

	engine15, err := scriptService.NewEngine(script15)
	So(err, ShouldBeNil)
	err = engine15.Compile()
	So(err, ShouldBeNil)
	script15Id, err := adaptors.Script.Add(script15)
	So(err, ShouldBeNil)
	script15, err = adaptors.Script.GetById(script15Id)
	So(err, ShouldBeNil)

	scripts["mb_dev1_turn_off_fan1_v1"] = script15

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

	// workflow1
	// ------------------------------------------------
	script16 := &m.Script{
		Lang:        "coffeescript",
		Name:        "wflow_script_v1",
		Source:      WflowScriptV1,
		Description: "workflow script",
	}
	ok, _ = script16.Valid()
	So(ok, ShouldEqual, true)

	engine16, err := scriptService.NewEngine(script16)
	So(err, ShouldBeNil)
	err = engine16.Compile()
	So(err, ShouldBeNil)
	script16Id, err := adaptors.Script.Add(script16)
	So(err, ShouldBeNil)
	script16, err = adaptors.Script.GetById(script16Id)
	So(err, ShouldBeNil)

	scripts["wflow_script_v1"] = script16

	return
}

const MbDev1ConditionCheckV1 = `

objects = [
    {name:'dev1_light1', id:0,systemName:'LIGHT_1_'}
    {name:'dev1_light3', id:2,systemName:'LIGHT_3_'}
    {name:'dev1_light2', id:1,systemName:'LIGHT_2_'}
    {name:'dev1_light4', id:3,systemName:'LIGHT_4_'}
    {name:'dev1_fan1', id:4,systemName:'FAN_1_'}
]

temps = [
    {name:'dev1_temp1', id:5,systemName:'TEMP_1_'}
    {name:'dev1_temp2', id:6,systemName:'TEMP_2_'}
]

doors = [
    {name:'dev1_door_main', id:7,systemName:'DOOR_MAIN_'}
    {name:'dev1_door_second', id:8,systemName:'DOOR_SECOND_'}
]

getStatus =(status)->
	if status == 1
		return 'ON'
	else
		return 'OFF'

doorStatus =(status)->
	if status == 1
		return 'OPENED'
	else
		return 'CLOSED'

fetchStatus =->

    COMMAND = []
    FUNC = 'ReadHoldingRegisters'
    ADDRESS = 0
    COUNT = 16
    
    res = device.modBus FUNC, ADDRESS, COUNT, COMMAND
    if res.error
        print 'error: ', res.error
        objects.forEach (obj)->
            IC.Map.setElementState device.getModel(), obj.name, 'ERROR'
            return
        temps.forEach (obj)->
            IC.Map.setElementState device.getModel(), obj.name, 'ERROR'
            return
        doors.forEach (obj)->
            IC.Map.setElementState device.getModel(), obj.name, 'ERROR'
            return
        return
    else 
        # print 'ok: ', res.result
        objects.forEach (obj)->
            newStatus = getStatus(res.result[obj.id])
            IC.Map.setElementState device.getModel(), obj.name, obj.systemName + newStatus
            return
            
        doors.forEach (obj)->
            newStatus = doorStatus(res.result[obj.id])
            IC.Map.setElementState device.getModel(), obj.name, obj.systemName + newStatus
            return
        
        temps.forEach (obj)->
            IC.Map.setElementState device.getModel(), obj.name, obj.systemName + 'ON'

            element = IC.Map.getElement device.getModel(), obj.name
            temp = if res.result[obj.id] then res.result[obj.id] else 0
            element.setOptions {'text': temp}
            return

        return
    
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

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOnLight3V1 = `
# turn on first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 2, 1, [1]
    if res.error
        print 'error: ', res.error

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOffLight3V1 = `
# turn off first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 2, 1, [0]
    if res.error
        print 'error: ', res.error

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOnLight4V1 = `
# turn on first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 3, 1, [1]
    if res.error
        print 'error: ', res.error

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOffLight4V1 = `
# turn off first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 3, 1, [0]
    if res.error
        print 'error: ', res.error

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`

const MbDev1TurnOnFan1V1 = `
# turn on first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 4, 1, [1]
    if res.error
        print 'error: ', res.error

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus()

main()
`
const MbDev1TurnOffFan1V1 = `
# turn off first light
fetchStatus =->
    
    res = device.modBus 'WriteMultipleRegisters', 4, 1, [0]
    if res.error
        print 'error: ', res.error

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
    scenario = 'weekday'
    print WFLOW_VAR1, 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario

on_exit =->
    
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
    scenario = 'weekend'
    print WFLOW_VAR1, 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario

on_exit =->

`

const BaseScript = `

main =->
`

const WflowScriptV1 = `
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

# global variable
WFLOW_VAR1 = 'workflow1'
`
