// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package scripts

var coffeeScripts = map[string]string{
	"coffeeScript1":  coffeeScript1,
	"coffeeScript2":  coffeeScript2,
	"coffeeScript3":  coffeeScript3,
	"coffeeScript4":  coffeeScript4,
	"coffeeScript5":  coffeeScript5,
	"coffeeScript6":  coffeeScript6,
	"coffeeScript7":  coffeeScript7,
	"coffeeScript8":  coffeeScript8,
	"coffeeScript9":  coffeeScript9,
	"coffeeScript10": coffeeScript10,
	"coffeeScript11": coffeeScript11,
	"coffeeScript12": coffeeScript12,
	"coffeeScript13": coffeeScript13,
	"coffeeScript14": coffeeScript14,
	"coffeeScript15": coffeeScript15,
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

// test1
// ------------------------------------------------
const coffeeScript1 = `
"use strict";

r = ExecuteSync "data/scripts/ping.sh", "google.com"
if r.out == 'ok'
    print "site is available ^^"

store(r.out)
`

// test2
// ------------------------------------------------
const coffeeScript2 = `
"use strict";

require('data/scripts/util.js')
require('data/scripts/vars.js')

# var foo
#print 'foo value:', foo

# var bar
#print 'bar value:', bar

# func from external lib
#print 'run external function: ', area(4)

store(foo + '-' + bar + '-' + months[0])
`

// test3
// ------------------------------------------------
const coffeeScript3 = `
on_enter =->
	store('on_enter')	

on_exit =->
	store('on_exit')
`

// test5,6
// ------------------------------------------------
const coffeeScript4 = `
#print "run workflow script (script 4)"
wf = Workflow

# name
name = wf.GetName()
store(name)

store(wf.GetDescription())

wf.SetVar('var', 'foo')
store(wf.GetVar('var'))
store(wf.GetScenario())
store(wf.GetScenarioName())
`

const coffeeScript5 = `
#print "run workflow script (script 5)"

wf = Workflow
store(wf.GetVar('var'))
wf.SetVar('var', 'bar')

`

const coffeeScript6 = `
#print "workflow scenario script 6"

on_enter =->
	wf = Workflow
	store(wf.GetName())
	store(wf.GetDescription())
	store(wf.GetVar('var'))
	wf.SetScenario('wf_scenario_2')
	wf.SetVar('var2', 'fee')

on_exit =->

`

const coffeeScript7 = `
#print "workflow scenario script 7"

on_enter =->

on_exit =->

`

const coffeeScript8 = `
#print "run flow script (script 8)"

`

const coffeeScript17 = `
#print "run flow script (script 17)"

`

// test7
// ------------------------------------------------
const coffeeScript9 = `
#print "run workflow script (script 9)"

store('enter script9')

Workflow.SetVar('workflow_var', 'foo')
`
const coffeeScript10 = `
#print "run workflow script (script 10)"

on_enter =->
	store('enter script10')
	So(Workflow.GetVar('workflow_var'), 'ShouldEqual', 'foo')
	Workflow.SetVar('workflow_var', 'bar')

on_exit =->
	store('exit script10')

`

const coffeeScript11 = `
#print "run flow script (script 11)"

store('enter script11')

wf = Workflow
flow = Flow

So(wf.GetScenario(), 'ShouldEqual', 'wf_scenario_1')
So(wf.GetVar('workflow_var'), 'ShouldEqual', 'bar')
So(flow.GetName(), 'ShouldEqual', 'flow1')
So(flow.GetDescription(), 'ShouldEqual', 'flow1 desc')
So(CurrentNode(), 'ShouldEqual', null)
So(CurrentDevice(), 'ShouldEqual', null)
So(flow.Node(), 'ShouldEqual', null)

store(flow.GetVar('var1'))
flow.SetVar('var1', 'foo')
store(flow.GetVar('var1'))

So(message.GetVar('flow_var'), 'ShouldEqual', null)
message.SetVar('flow_var', 'foo2')
So(message.GetVar('flow_var'), 'ShouldEqual', 'foo2')
So(message.GetVar('val'), 'ShouldEqual', 1)
`
const coffeeScript12 = `
#print "run flow script (script 12)"

store('enter script12')

So(message.GetVar('flow_var'), 'ShouldEqual', 'foo2')
message.SetVar('flow_var', 'foo3')

store(flow.GetVar('var1'))
flow.SetVar('var1', 'bar')
`
const coffeeScript13 = `
#print "run flow script (script 13)"

store('enter script13')

So(message.GetVar('flow_var'), 'ShouldEqual', 'foo3')

f = Flow
store(f.GetVar('var1'))
store(flow.GetVar('var1'))
`

const coffeeScript14 = `
#print "run flow script (script 14)"

store('enter script14')

wf = Workflow
flow = Flow

So(wf.GetScenario(), 'ShouldEqual', 'wf_scenario_1')
So(wf.GetVar('workflow_var'), 'ShouldEqual', 'bar')
So(flow.GetName(), 'ShouldEqual', 'flow2')
So(flow.GetDescription(), 'ShouldEqual', 'flow2 desc')
So(CurrentNode(), 'ShouldEqual', null)
So(CurrentDevice(), 'ShouldEqual', null)
So(flow.Node(), 'ShouldEqual', null)
So(flow.GetVar('var1'), 'ShouldEqual', null)
So(message.GetVar('flow_var'), 'ShouldEqual', null)

message.SetVar('flow_var', 123)
So(message.GetVar('val'), 'ShouldEqual', 2)
`

const coffeeScript15 = `
#print "run flow script (script 15)"

store('enter script15')

So(message.GetVar('flow_var'), 'ShouldEqual', 123)
`

const coffeeScript16 = `
#print "run flow script (script 16)"

store('enter script16')

So(message.GetVar('flow_var'), 'ShouldEqual', 123)
message.SetVar('val', 123)
`

// test8
// ------------------------------------------------
const coffeeScript18 = `
#print "run workflow script (script 18)"

store('enter script18')
`

const coffeeScript19 = `
#print "run workflow script (script 19)"

store('enter script19')
`

const coffeeScript20 = `
#print "workflow scenario script 20"

on_enter =->
	store('enter script20')

on_exit =->
	store('exit script20')

`

const coffeeScript21 = `
#print "workflow scenario script 21"

on_enter =->
	store('enter script21')

on_exit =->
	store('exit script21')

`

const coffeeScript22 = `
#print "run flow script (script 22)"

store('enter script22')

wf = Workflow
wf.SetScenario('wf_scenario_2')
#print "text after select scenario...."
`

const coffeeScript23 = `
#print "run flow script (script 23)"

store('enter script23')
`

const coffeeScript24 = `
print "run flow script (script 24)"

main =->

`

// bar
// &{foo <nil>}
// foo
// <nil>
const coffeeScript25 = `
"use strict";

So(bar.bar, 'ShouldEqual', 'bar')
So(bar.foo, 'ShouldEqual', '&{foo <nil>}')
So(bar.foo.bar, 'ShouldEqual', 'foo')
So(bar.foo.foo, 'ShouldEqual', '<nil>')

So(bar2.bar, 'ShouldEqual', 'bar')
So(bar2.foo, 'ShouldEqual', '&{foo <nil>}')
So(bar2.foo.bar, 'ShouldEqual', 'foo')
So(bar2.foo.foo, 'ShouldEqual', '<nil>')

So(bar2.bar, 'ShouldEqual', 'bar')
So(bar2.foo, 'ShouldEqual', '&{foo <nil>}')
So(bar2.foo.bar, 'ShouldEqual', 'foo')
So(bar2.foo.foo, 'ShouldEqual', '<nil>')

external('bar', bar)
external('bar2', bar2)
external('bar2', bar2)
`

const coffeeScript26 = `
"use strict";

store(bar2.bar)
store(bar2.foo)
store(bar2.foo.bar)
store(bar2.foo.foo)
`

// test27
// ------------------------------------------------
const coffeeScript27 = `
"use strict";

# get nil var
value = Storage.getByName 'foo'
store value

foo = 
	'bar': 'bar'

value = marshal foo

# save var
Storage.push 'foo', value

# get exist var
value = Storage.getByName 'foo'
store value
foo = unmarshal value
store foo.bar

# update var and save
foo.bar = 'foo'
value = marshal foo
Storage.push 'foo', value

value = Storage.getByName 'foo'
store value
foo = unmarshal value
store foo.bar

list = Storage.search 'bar'
store list

list = Storage.search 'foo'
store list
`

// test...
// ------------------------------------------------
