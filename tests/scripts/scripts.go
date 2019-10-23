package scripts

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
	"coffeeScript15": coffeeScript15,
	"coffeeScript16": coffeeScript16,
}

// test1
// ------------------------------------------------
const coffeeScript1 = `
"use strict";

r = IC.ExecuteSync "data/scripts/ping.sh", "google.com"

if r.out == 'ok'
    print "site is available ^^"

IC.store(r.out)
`

// test2
// ------------------------------------------------
const coffeeScript2 = `
"use strict";

require('data/scripts/util.js')
require('data/scripts/vars.js')

# var foo
print 'foo value:', foo

# var bar
print 'bar value:', bar

# func from external lib
print 'run external function: ', area(4)

console.log months[0]

IC.store(foo + '-' + bar + '-' + months[0])
`

// test3
// ------------------------------------------------
const coffeeScript3 = `
on_enter =->
	IC.store('on_enter')	

on_exit =->
	IC.store('on_exit')
`

// test5,6
// ------------------------------------------------
const coffeeScript4 = `
print "run workflow script (script 4)"
wf = IC.Workflow()

# name
name = wf.getName()
IC.store(name)

IC.store(wf.getDescription())

wf.setVar('var', 'foo')
IC.store(wf.getVar('var'))
IC.store(wf.getScenario())
IC.store(wf.getScenarioName())
`

const coffeeScript5 = `
print "run workflow script (script 5)"

wf = IC.Workflow()
IC.store(wf.getVar('var'))
wf.setVar('var', 'bar')

`

const coffeeScript6 = `
print "workflow scenario script 6"

on_enter =->
	wf = IC.Workflow()
	IC.store(wf.getName())
	IC.store(wf.getDescription())
	IC.store(wf.getVar('var'))
	#wf.setScenario('wf_scenario_2')
	wf.setVar('var2', 'fee')

on_exit =->
	IC.store('exit from wf_scenario_1')
`

const coffeeScript7 = `
on_enter =->
	wf = IC.Workflow()
	IC.store('enter to wf_scenario_2')

on_exit =->
	wf = IC.Workflow()

`

const coffeeScript8 = `
print "run flow script (script 8)"

wf = IC.Workflow()
IC.store(wf.getName())
IC.store(wf.getDescription())
IC.store(wf.getVar('var2'))
#wf.setScenario('wf_scenario_2')

`

// test7
// ------------------------------------------------
const coffeeScript9 = `
print "run workflow script (script 9)"

wf = IC.Workflow()
wf.setVar('workflow_var', 'foo')
`
const coffeeScript10 = `
on_enter =->
	wf = IC.Workflow()
	IC.So(wf.getVar('workflow_var'), 'ShouldEqual', 'foo')
	wf.setVar('workflow_var', 'bar')
	
on_exit =->
	
`

const coffeeScript11 = `
print "run flow script (script 11)"

wf = IC.Workflow()
flow = IC.Flow()

IC.So(wf.getScenario(), 'ShouldEqual', 'wf_scenario_1')
IC.So(wf.getVar('workflow_var'), 'ShouldEqual', 'bar')
IC.So(flow.getName(), 'ShouldEqual', 'flow1')
IC.So(flow.getDescription(), 'ShouldEqual', 'flow1 desc')
IC.So(IC.CurrentNode(), 'ShouldEqual', null)
IC.So(IC.CurrentDevice(), 'ShouldEqual', null)
IC.So(flow.node(), 'ShouldEqual', null)

IC.store(flow.getVar('var1'))
flow.setVar('var1', 'foo')
IC.store(flow.getVar('var1'))

IC.So(message.getVar('flow_var'), 'ShouldEqual', null)
message.setVar('flow_var', 'foo2')
IC.So(message.getVar('flow_var'), 'ShouldEqual', 'foo2')
IC.So(message.getVar('val'), 'ShouldEqual', 1)
`
const coffeeScript12 = `
print "run flow script (script 12)"

IC.So(message.getVar('flow_var'), 'ShouldEqual', 'foo2')
message.setVar('flow_var', 'foo3')

IC.store(flow.getVar('var1'))
flow.setVar('var1', 'bar')
`
const coffeeScript13 = `
print "run flow script (script 13)"

IC.So(message.getVar('flow_var'), 'ShouldEqual', 'foo3')

f = IC.Flow()
IC.store(f.getVar('var1'))
IC.store(flow.getVar('var1'))
`

const coffeeScript14 = `
print "run flow script (script 14)"

wf = IC.Workflow()
flow = IC.Flow()

IC.So(wf.getScenario(), 'ShouldEqual', 'wf_scenario_1')
IC.So(wf.getVar('workflow_var'), 'ShouldEqual', 'bar')
IC.So(flow.getName(), 'ShouldEqual', 'flow2')
IC.So(flow.getDescription(), 'ShouldEqual', 'flow2 desc')
IC.So(IC.CurrentNode(), 'ShouldEqual', null)
IC.So(IC.CurrentDevice(), 'ShouldEqual', null)
IC.So(flow.node(), 'ShouldEqual', null)
IC.So(flow.getVar('var1'), 'ShouldEqual', null)
IC.So(message.getVar('flow_var'), 'ShouldEqual', null)

message.setVar('flow_var', 123)
IC.So(message.getVar('val'), 'ShouldEqual', 2)
`

const coffeeScript15 = `
print "run flow script (script 15)"

IC.So(message.getVar('flow_var'), 'ShouldEqual', 123)
`

const coffeeScript16 = `
print "run flow script (script 16)"

IC.So(message.getVar('flow_var'), 'ShouldEqual', 123)
message.setVar('val', 123)
`

// test...
// ------------------------------------------------
