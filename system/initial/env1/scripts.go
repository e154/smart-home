package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

func addScripts(adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (script1, script2, script3, script4, script5, script6 *m.Script) {

	// add script
	// ------------------------------------------------
	script1 = &m.Script{
		Lang:        "coffeescript",
		Name:        "turn_on_the_socket",
		Source:      coffeescript1,
		Description: "test1",
	}
	ok, _ := script1.Valid()
	So(ok, ShouldEqual, true)

	engine1, err := scriptService.NewEngine(script1)
	So(err, ShouldBeNil)
	err = engine1.Compile()
	So(err, ShouldBeNil)
	script1Id, err := adaptors.Script.Add(script1)
	So(err, ShouldBeNil)
	script1, err = adaptors.Script.GetById(script1Id)
	So(err, ShouldBeNil)

	script2 = &m.Script{
		Lang:        "coffeescript",
		Name:        "turn_off_the_socket",
		Source:      coffeescript2,
		Description: "script2",
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

	script3 = &m.Script{
		Lang:        "coffeescript",
		Name:        "socket_status_сheck",
		Source:      coffeescript3,
		Description: "script3",
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

	script4 = &m.Script{
		Lang:        "coffeescript",
		Name:        "scenario_weekday",
		Source:      coffeescript4,
		Description: "script4",
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

	script5 = &m.Script{
		Lang:        "coffeescript",
		Name:        "scenario_weekend",
		Source:      coffeescript5,
		Description: "script5",
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

	script6 = &m.Script{
		Lang:        "coffeescript",
		Name:        "task1",
		Source:      coffeescript6,
		Description: "script6",
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

	return
}

const coffeescript1 = `
# Контекст применения: 
# action (действие)
#
# Описание:
# Включение устройства. (частное)
# Не имеет зависимостей, и ни чего не передает наружу
# Должен вызываться в рамках воркера, или действия устройства,
# иначе выдаст ошибку, так как контекст выполнения накладывает 
# некоторые ограничения

fetchStatus =(node, dev)->
    
    FUNCTION = 4
    DEVICE_ADDR = dev.getAddress()
    TOGGLE = 1
    
    COMMAND = [DEVICE_ADDR, FUNCTION, 0, TOGGLE, 0, 0]

    from_node = node.send dev, true, COMMAND
    
    # map element
    element = IC.Map.getElement dev
    
    if from_node.error
        message.setError from_node.error
        element.setState 'ERROR'
        
        # IC.Log.error "#{message.dev.name} - error: #{from_node.error}"
        
        return false
       
    if from_node.result != ""
        result = IC.hex2arr(from_node.result)
        
        # check power
        if result[2] == 1
            element.setState 'ENABLED'
        else
            element.setState 'DISABLED'
        
    from_node.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
     
    fetchStatus(node, dev)

main()
`

const coffeescript2 = `
# Контекст применения: 
# action (действие)
#
# Описание:
# Выключение устройства. (частное)
# Не имеет зависимостей, и ни чего не передает наружу
# Должен вызываться в рамках воркера, или действия устройства,
# иначе выдаст ошибку, так как контекст выполнения накладывает 
# некоторые ограничения

fetchStatus =(node, dev)->
    
    FUNCTION = 4
    DEVICE_ADDR = dev.getAddress()
    TOGGLE = 0
    
    COMMAND = [DEVICE_ADDR, FUNCTION, 0, TOGGLE, 0, 0]
    
    from_node = node.send dev, true, COMMAND
    
    # map element
    element = IC.Map.getElement dev
    
    if from_node.error
        message.setError from_node.error
        element.setState 'ERROR'
        
        # IC.Log.error "#{message.dev.name} - error: #{from_node.error}"
        
        return false
       
    if from_node.result != ""
        result = IC.hex2arr(from_node.result)
        
        # check power
        if result[2] == 1
            element.setState 'ENABLED'
        else
            element.setState 'DISABLED'
        
    from_node.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
     
    fetchStatus(node, dev)

main()
`

const coffeescript3 = `
# Контекст применения: 
# action (действие)
#
# Описание:
# Проверка состояния устройства. (частное)
# Не имеет зависимостей, и ни чего не передает наружу
# Должен вызываться в рамках воркера, или действия устройства,
# иначе выдаст ошибку, так как контекст выполнения накладывает 
# некоторые ограничения

fetchStatus =(node, dev, flow)->
    
    # номер комманнды 
    # 3 - проверка состояния
    FUNCTION = 3

    # получим адрес устройства из контекста запуска
    DEVICE_ADDR = dev.getAddress()
    
    COMMAND = [DEVICE_ADDR, FUNCTION, 0, 0, 0, 5]
    
    # map element
    element = IC.Map.getElement dev
    element.setOptions {text: 'some text'}
    
    print "check status, flow:", flow.getName(), "dev:", DEVICE_ADDR
    
    # запрос состояния устройства
    from_node = node.send dev, true, COMMAND
    
    # if 'Лампа в зале' == dev.getName()
    #     print '---', from_node
    
    if from_node.error
        # if 'Лампа в зале' == dev.getName()
        #     print 'error'
            
        message.setError from_node.error
        element.setState 'ERROR'
        
        # IC.Log.error "#{dev.name} - error: #{from_node.error}"
        return false
       
    if from_node.result != ""
        result = IC.hex2arr(from_node.result)
        
        # check power
        if result[2] == 1
            element.setState 'ENABLED'
        else
            element.setState 'DISABLED'
    
        print 'result', from_node.result, 'dev:', DEVICE_ADDR
    # print 'dev:', DEVICE_ADDR, 'state', element.getState().systemName
    
    from_node.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    flow = IC.Flow()
    
    return if !node || !dev
    
    # номер комманнды 
    # 3 - проверка состояния
    FUNCTION = 3
    
    COMMAND = [FUNCTION, 0, 0, 0, 5]
    print 1
    node.send dev, COMMAND
    print 2

    #fetchStatus(node, dev, flow)
    
main()
`

const coffeescript4 = `
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
    print 'on enter'

on_exit =->
    print 'on exit'
    
main =->
    scenario = 'weekday'
    print 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario
    
`

const coffeescript5 = `
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
    print 'on enter'

on_exit =->
    print 'on exit'
    
main =->
    scenario = 'weekend'
    print 'scenario', scenario
    IC.Workflow().setVar 'scenario', scenario
`

const coffeescript6 = `

main =->
    incoming = IC.hex2arr message.getVar('result')
    print IC.Flow().getName(), 'incoming:', incoming
    print IC.Map.getDeviceState(message.getVar('devId'))
    # IC.Log.warn 'test'
    # print IC.Workflow().getVar 'scenario'
    # scenario = IC.Workflow().getScenario()
    # print scenario
    # if scenario == 'weekday'
    #     IC.Workflow().setScenario 'weekend' 
`