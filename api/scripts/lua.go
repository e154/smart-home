package scripts

import (
	//"github.com/aarzilli/golua/lua"
	//"github.com/stevedonovan/luar"
)

type Lua struct {
	//state  *lua.State
	engine *Engine
}

func (l *Lua) Init() (err error)  {
	//l.state = luar.Init()
	//l.state.OpenLibs();
	return
}

func (l *Lua) Close() {
	//l.state.Close()
}

func (j *Lua) Compile() (err error) {

	return
}

func (l *Lua) Do() (result string, err error) {

	result = "Lua scripts are not supported yet"
	return
}

func (j *Lua) PushStruct(name string, i interface{}) (int, error) {

	return 0, nil
}

func (j *Lua) PushFunction(name string, i interface{}) (int, error) {

	return 0, nil
}