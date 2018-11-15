package common

type ScriptLang string

const (
	ScriptLangTs         = ScriptLang("ts")
	ScriptLangCoffee     = ScriptLang("coffeescript")
	ScriptLangJavascript = ScriptLang("javascript")
)

type FlowElementsPrototypeType string

const (
	FlowElementsPrototypeDefault        = FlowElementsPrototypeType("default")
	FlowElementsPrototypeMessageHandler = FlowElementsPrototypeType("MessageHandler")
	FlowElementsPrototypeMessageEmitter = FlowElementsPrototypeType("MessageEmitter")
	FlowElementsPrototypeTask           = FlowElementsPrototypeType("Task")
	FlowElementsPrototypeGateway        = FlowElementsPrototypeType("Gateway")
	FlowElementsPrototypeFlow           = FlowElementsPrototypeType("Flow")
)

type StatusType string

const (
	Enabled  = StatusType("enabled")
	Disabled = StatusType("disabled")
)
