// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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

	// mb_dev1_actions
	// ------------------------------------------------
	script2 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mb_dev1_actions",
		Source:      MbDev1ActionsV1,
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

	scripts["mb_dev1_actions"] = script2

	// mi_pir_sensor
	// ------------------------------------------------
	script10 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mi_pir_sensor",
		Source:      MiPirSensor,
		Description: "mi pir sensor",
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

	scripts["mi_pir_sensor"] = script10

	// mi_pir_sensor
	// ------------------------------------------------
	script12 := &m.Script{
		Lang:        "coffeescript",
		Name:        "mi_door_sensor",
		Source:      MiDoorSensor,
		Description: "mi door sensor",
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

	scripts["mi_door_sensor"] = script12

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
    
    res = Device.ModBus FUNC, ADDRESS, COUNT, COMMAND
    if res.error
        print 'error: ', res.error
        objects.forEach (obj)->
            Map.SetElementState obj.name, 'ERROR'
            return
        temps.forEach (obj)->
            Map.SetElementState obj.name, 'ERROR'
            return
        doors.forEach (obj)->
            Map.SetElementState obj.name, 'ERROR'
            return
        return
    else 
        # print 'ok: ', res.result
        objects.forEach (obj)->
            newStatus = getStatus(res.Result[obj.id])
            Map.SetElementState obj.name, obj.systemName + newStatus
            return
            
        doors.forEach (obj)->
            newStatus = doorStatus(res.Result[obj.id])
            Map.SetElementState obj.name, obj.systemName + newStatus
            return
        
        temps.forEach (obj)->
            Map.SetElementState obj.name, obj.systemName + 'ON'

            element = Map.GetElement obj.name
            temp = if res.Result[obj.id] then res.Result[obj.id] else 0
            element.SetOptions {'text': temp}
            return

        return
    
main =->
    fetchStatus()

main()
`
const MbDev1ActionsV1 = `
writeRegisters =(d, c, r)->
    res = Device.ModBus 'WriteMultipleRegisters', d, c, r
    if res.error
        print 'error: ', res.error

main =->
    switch Action.Name
        when "turn on light1" then writeRegisters(0, 1, [1])
        when "turn off light1" then writeRegisters(0, 1, [0])
        when "turn on light2" then writeRegisters(1, 1, [1])
        when "turn off light2" then writeRegisters(1, 1, [0])
        when "turn on light3" then writeRegisters(2, 1, [1])
        when "turn off light3" then writeRegisters(2, 1, [0])
        when "turn on light4" then writeRegisters(3, 1, [1])
        when "turn off light4" then writeRegisters(3, 1, [0])
        when "turn on fan1" then writeRegisters(4, 1, [1])
        when "turn off fan1" then writeRegisters(4, 1, [0])
        when "turn on all lights" then writeRegisters(0, 4, [1,1,1,1])
        when "turn off all lights" then writeRegisters(0, 4, [0,0,0,0])

main()
`

const MiPirSensor = `
options = 
    text:'' 
stateStatus = 'SILENCE'

main =->
    return if !message.Mqtt
    payload = JSON.parse(message.GetVar('mqtt_payload'))
    
    if payload['occupancy']
        stateStatus = 'OCCUPANCY'
    if payload['battery'] < 50
        stateStatus = 'WARNING'
        options.text = 'battery:'+ payload['battery']

main()
`

const MiDoorSensor = `
options = 
    text:'' 
stateStatus = 'CLOSED'

main =->
    return if !message.Mqtt
    payload = JSON.parse(message.GetVar('mqtt_payload'))
    
    if !payload['contact']
        stateStatus = 'OPENED'
    if payload['battery'] < 50
        stateStatus = 'WARNING'
        options.text = 'battery:'+ payload['battery']

main()
`

const CmdConditionCheckV1 = `
main =->
    
    NAME = './data/scripts/ping.sh'
    ARGS = ['ya.ru']
    result = Device.RunCommand NAME, ARGS
    print result

main()
`

const WflowScenarioWeekdayV1 = `
# variables:

#Workflow
#    .GetName()
#    .GetDescription()
#    .SetVar(string, interface)
#    .GetVar(string)
#    .GetScenario() string
#    .SetScenario(string)

#
# workflow script example
#

on_enter =->
    scenario = 'weekday'
    print WFLOW_VAR1, 'scenario', scenario
    Workflow.SetVar 'scenario', scenario

on_exit =->
    
`

const WflowScenarioWeekendV1 = `
# variables:

#Workflow
#    .GetName()
#    .GetDescription()
#    .SetVar(string, interface)
#    .GetVar(string)
#    .GetScenario() string
#    .SetScenario(string)

#
# workflow script example
#

on_enter =->
    scenario = 'weekend'
    print WFLOW_VAR1, 'scenario', scenario
    Workflow.SetVar 'scenario', scenario

on_exit =->

`

const BaseScript = `

main =->
`

const WflowScriptV1 = `
# variables:

#Workflow
#    .GetName()
#    .GetDescription()
#    .SetVar(string, interface)
#    .GetVar(string)
#    .GetScenario() string
#    .SetScenario(string)

#
# workflow script example
#

# global variable
WFLOW_VAR1 = 'workflow1'
`

