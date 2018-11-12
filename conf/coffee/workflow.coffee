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

main =->
    print "workflow"

on_enter =->
    wf = IC.Workflow()

    print "enter to " + wf.getName()
    print "description " + wf.getDescription()

    wf.setVar("var1", 123)

    print "scenario " + wf.getScenario()
    wf.setScenario("test2")


on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName()
    print "description " + wf.getDescription()

