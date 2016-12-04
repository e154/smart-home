package scripts

import (
	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Lua struct {
	State *lua.State
	engine	*Engine
}

func (l *Lua) Init() (err error)  {
	l.State = luar.Init()
	l.State.OpenLibs();
	return
}

func (l *Lua) Close() {
	l.State.Close()
}

func (j *Lua) Compile() (err error) {

	return
}

func (l *Lua) Do() (result string, err error) {

	result = "Lua scripts are not supported yet"
	return
}

func (j *Lua) Reg() (err error) {

	return
}

