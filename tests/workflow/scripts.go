package workflow

const coffeeScript1 = `
main =->
    wf = IC.Workflow()
    print wf.getName() + " script1"
`

const coffeeScript2 = `
main =->
    wf = IC.Workflow()
    print wf.getName() + " script2"
`

const coffeeScript3 = `
main =->
    print "workflow scenario script 3"

on_enter =->
    wf = IC.Workflow()
    print "enter to " + wf.getName() + " " + wf.getScenarioName()
    wf.setVar("var1", 123)
    wf.setScenario("wf_scenario_2")

on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName() + " " + wf.getScenarioName()
    print "description " + wf.getDescription()
`



const coffeeScript4 = `
main =->
    print "workflow scenario script 4"

on_enter =->
    wf = IC.Workflow()

    print "enter to " + wf.getName() + " " + wf.getScenarioName()
    var1 = wf.getVar("var1")
    print "var: " + var1

on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName() + " " + wf.getScenarioName()
    print "description " + wf.getDescription()
`

