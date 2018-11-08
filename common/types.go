package common

type ScriptLang string

const (
	ScriptLangLua = ScriptLang("lua")
	ScriptLangCoffee = ScriptLang("coffeescript")
	ScriptLangJavascript = ScriptLang("javascript")
)