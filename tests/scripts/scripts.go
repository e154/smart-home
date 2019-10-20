package scripts

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

// test5
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
print "run workflow script (script 8)"

wf = IC.Workflow()
IC.store(wf.getName())
IC.store(wf.getDescription())
IC.store(wf.getVar('var2'))
wf.setScenario('wf_scenario_2')

`


// test...
// ------------------------------------------------
