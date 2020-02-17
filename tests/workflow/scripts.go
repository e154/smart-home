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

package workflow

var coffeeScripts = map[string]string{
	"coffeeScript1": coffeeScript1,
	"coffeeScript2": coffeeScript2,
	"coffeeScript3": coffeeScript3,
	"coffeeScript4": coffeeScript4,
	"coffeeScript5": coffeeScript5,
	"coffeeScript6": coffeeScript6,
	"coffeeScript7": coffeeScript7,
	"coffeeScript8": coffeeScript8,
	"coffeeScript9": coffeeScript9,
	"coffeeScript10": coffeeScript10,
	"coffeeScript11": coffeeScript11,
	"coffeeScript12": coffeeScript12,
	"coffeeScript13": coffeeScript13,
	"coffeeScript14": coffeeScript14,
	//"coffeeScript15": coffeeScript15,
	"coffeeScript16": coffeeScript16,
	"coffeeScript17": coffeeScript17,
	"coffeeScript18": coffeeScript18,
	"coffeeScript19": coffeeScript19,
	"coffeeScript20": coffeeScript20,
	"coffeeScript21": coffeeScript21,
	"coffeeScript22": coffeeScript22,
	"coffeeScript23": coffeeScript23,
	"coffeeScript24": coffeeScript24,
	"coffeeScript25": coffeeScript25,
	"coffeeScript26": coffeeScript26,
	"coffeeScript27": coffeeScript27,
}

// test1, test2
// ------------------------------------------------
const coffeeScript1 = `
main =->
    wf = Workflow
    #print wf.GetName() + " script1"
`

const coffeeScript2 = `
main =->
    wf = Workflow
    #print wf.GetName() + " script2"
`

const coffeeScript3 = `
main =->
    #print "workflow scenario script 3"

on_enter =->
    wf = Workflow
    #print "enter to " + wf.GetName() + " " + wf.GetScenarioName()
    wf.SetVar("var1", 123)
    wf.SetScenario("wf_scenario_2")

on_exit =->
    wf = Workflow

    #print "exit from " + wf.GetName() + " " + wf.GetScenarioName()
    #print "description " + wf.getDescription()
`

const coffeeScript4 = `
main =->
    #print "workflow scenario script 4"

on_enter =->
    wf = Workflow

    #print "enter to " + wf.GetName() + " " + wf.GetScenarioName()
    var1 = wf.GetVar("var1")
    #print "var: " + var1

on_exit =->
    wf = Workflow

    #print "exit from " + wf.GetName() + " " + wf.GetScenarioName()
    #print "description " + wf.getDescription()
`

const coffeeScript5 = `
# test 2

#print "run message emitter script (script 5)"
#print "message:", message
#print "ENV", ENV.a

#print "message.ENV:", message.GetVar "ENV"

store ENV.a
`

const coffeeScript6 = `
# test 2

ENV = {"a": "b"}
#print "run message handler script (script 6)"

message.SetVar("ENV", ENV) 
`

const coffeeScript7 = `
# test 2
# workflow script

#print "run workflow script (script 7)"

on_enter =->
	#print 'on_enter'
on_exit =->
`

// test3
// ------------------------------------------------
const coffeeScript8 = `
#print "run workflow script (script 8)"
`

const coffeeScript9 = `
#print "run workflow script (script 9)"
`

// test 4
// ------------------------------------------------
const coffeeScript10 = `
#print "run workflow script (script 10)"
`

const coffeeScript11 = `
#print "run workflow script (script 11)"
`

// test6, test7
// ------------------------------------------------
const coffeeScript12 = `
#print "run workflow script (script 12)"
store('script12')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript13 = `
#print "run workflow script (script 13)"
store('script13')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript14 = `
#print "run workflow script (script 14)"
store('script14')
#print "message:", message.GetVar('val')
`

// test8, test9, test10
// ------------------------------------------------
const coffeeScript16 = `
#print "run workflow script (script 16)"
store('script16')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript17 = `
#print "run workflow script (script 17)"

store('script17')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript18 = `
#print "run workflow script (script 18)"

store('script18')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript19 = `
#print "run workflow script (script 19)"

store('script19')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript20 = `
#print "run workflow script (script 20)"

store('script20')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript21 = `
#print "run workflow script (script 21)"

store('script21')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

// test11
// ------------------------------------------------
const coffeeScript22 = `
#print "run workflow script (script 22)"
store('script22')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript23 = `
#print "run workflow script (script 23)"
store('script23')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript24 = `
#print "run workflow script (script 24)"
store('script24')
#print "message:", message.GetVar('val')
c = message.GetVar('val') + 1
message.SetVar('val', c)
store2(c)
`

const coffeeScript25 = `
#print "run workflow script (script 25)"
`

const coffeeScript26 = `
#print "run workflow script (script 26)"
`

const coffeeScript27 = `
#print "run workflow script (script 27)"
`

// test8...
// ------------------------------------------------
