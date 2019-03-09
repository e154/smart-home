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

type DeviceType string

type PrototypeType string

const (
	PrototypeTypeText   = PrototypeType("text")
	PrototypeTypeImage  = PrototypeType("image")
	PrototypeTypeDevice = PrototypeType("device")
)

type LogLevel string

const (
	LogLevelEmergency = LogLevel("Emergency")
	LogLevelAlert     = LogLevel("Alert")
	LogLevelCritical  = LogLevel("Critical")
	LogLevelError     = LogLevel("Error")
	LogLevelWarning   = LogLevel("Warning")
	LogLevelNotice    = LogLevel("Notice")
	LogLevelInfo      = LogLevel("Info")
	LogLevelDebug     = LogLevel("Debug")
)
