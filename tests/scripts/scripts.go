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
#print 'foo value:', foo

# var bar
#print 'bar value:', bar

# func from external lib
#print 'run external function: ', area(4)

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
#print "run workflow script (script 4)"
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
#print "run workflow script (script 5)"

wf = IC.Workflow()
IC.store(wf.getVar('var'))
wf.setVar('var', 'bar')

`

const coffeeScript6 = `
#print "workflow scenario script 6"

on_enter =->
	wf = IC.Workflow()
	IC.store(wf.getName())
	IC.store(wf.getDescription())
	IC.store(wf.getVar('var'))
	wf.setScenario('wf_scenario_2')
	wf.setVar('var2', 'fee')

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

IC.store('enter script9')

wf = IC.Workflow()
wf.setVar('workflow_var', 'foo')
`
const coffeeScript10 = `
#print "run workflow script (script 10)"

on_enter =->
	IC.store('enter script10')
	wf = IC.Workflow()
	IC.So(wf.getVar('workflow_var'), 'ShouldEqual', 'foo')
	wf.setVar('workflow_var', 'bar')

on_exit =->
	IC.store('exit script10')

`

const coffeeScript11 = `
#print "run flow script (script 11)"

IC.store('enter script11')

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
#print "run flow script (script 12)"

IC.store('enter script12')

IC.So(message.getVar('flow_var'), 'ShouldEqual', 'foo2')
message.setVar('flow_var', 'foo3')

IC.store(flow.getVar('var1'))
flow.setVar('var1', 'bar')
`
const coffeeScript13 = `
#print "run flow script (script 13)"

IC.store('enter script13')

IC.So(message.getVar('flow_var'), 'ShouldEqual', 'foo3')

f = IC.Flow()
IC.store(f.getVar('var1'))
IC.store(flow.getVar('var1'))
`

const coffeeScript14 = `
#print "run flow script (script 14)"

IC.store('enter script14')

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
#print "run flow script (script 15)"

IC.store('enter script15')

IC.So(message.getVar('flow_var'), 'ShouldEqual', 123)
`

const coffeeScript16 = `
#print "run flow script (script 16)"

IC.store('enter script16')

IC.So(message.getVar('flow_var'), 'ShouldEqual', 123)
message.setVar('val', 123)
`

// test8
// ------------------------------------------------
const coffeeScript18 = `
#print "run workflow script (script 18)"

IC.store('enter script18')
`

const coffeeScript19 = `
#print "run workflow script (script 19)"

IC.store('enter script19')
`

const coffeeScript20 = `
#print "workflow scenario script 20"

on_enter =->
	IC.store('enter script20')

on_exit =->
	IC.store('exit script20')

`

const coffeeScript21 = `
#print "workflow scenario script 21"

on_enter =->
	IC.store('enter script21')

on_exit =->
	IC.store('exit script21')

`

const coffeeScript22 = `
#print "run flow script (script 22)"

IC.store('enter script22')

wf = IC.Workflow()
wf.setScenario('wf_scenario_2')
#print "text after select scenario...."
`

const coffeeScript23 = `
#print "run flow script (script 23)"

IC.store('enter script23')
`

const coffeeScript24 = `
print "run flow script (script 24)"

main =->
	tpl = IC.Template.render('template2', {'code':'12345'})
	sms = IC.Notifr.newSMS()
	sms.addPhone('+79133819546')
	sms.setRender(tpl)
	IC.Notifr.send(sms)
	''

`

//bar
//&{foo <nil>}
//foo
//<nil>
const coffeeScript25 = `
"use strict";

IC.So(bar.bar, 'ShouldEqual', 'bar')
IC.So(bar.foo, 'ShouldEqual', '&{foo <nil>}')
IC.So(bar.foo.bar, 'ShouldEqual', 'foo')
IC.So(bar.foo.foo, 'ShouldEqual', '<nil>')

IC.So(IC.bar2.bar, 'ShouldEqual', 'bar')
IC.So(IC.bar2.foo, 'ShouldEqual', '&{foo <nil>}')
IC.So(IC.bar2.foo.bar, 'ShouldEqual', 'foo')
IC.So(IC.bar2.foo.foo, 'ShouldEqual', '<nil>')

IC.So(bar2.bar, 'ShouldEqual', 'bar')
IC.So(bar2.foo, 'ShouldEqual', '&{foo <nil>}')
IC.So(bar2.foo.bar, 'ShouldEqual', 'foo')
IC.So(bar2.foo.foo, 'ShouldEqual', '<nil>')

IC.external('bar', bar)
IC.external('bar2', bar2)
IC.external('IC.bar2', IC.bar2)
`

const coffeeScript26 = `
"use strict";

IC.store(bar2.bar)
IC.store(bar2.foo)
IC.store(bar2.foo.bar)
IC.store(bar2.foo.foo)
`

// test...
// ------------------------------------------------
