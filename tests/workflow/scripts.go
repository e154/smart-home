package workflow

const coffeeScript1 = `
main =->
    print "workflow script1"
`

const coffeeScript2 = `
main =->
    print "workflow script2"
`

const coffeeScript3 = `
main =->
    print "workflow scenario script 3"

on_enter =->
    wf = IC.Workflow()

    print "workflow scenario script 3"
    print "enter to " + wf.getName()
    print "description " + wf.getDescription()

    var1 = wf.getVar("var1", 123)

    print "var " + var1


on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName()
    print "description " + wf.getDescription()
`



const coffeeScript4 = `
main =->
    print "workflow scenario script 4"

on_enter =->
    wf = IC.Workflow()

    print "workflow scenario script 4"
    print "enter to " + wf.getName()
    print "description " + wf.getDescription()

    var1 = wf.getVar("var1", 123)

    print "var " + var1


on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName()
    print "description " + wf.getDescription()
`

