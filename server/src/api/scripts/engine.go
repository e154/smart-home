package scripts

import (
	"../models"
	"errors"
	"github.com/astaxie/beego"
)

type Magic interface {
	Init() error
	Do() (string, error)
	Compile() error
	PushStruct(string, interface{}) (int, error)
	PushFunction(string, interface{}) (int, error)
	Close()
}

func New(s *models.Script) (engine *Engine, err error) {

	engine = &Engine{
		model:s,
		buf: make([]string, 0),
	}

	switch s.Lang {
	case "lua":
		engine.script = &Lua{engine:engine}
	case "javascript":
		engine.script = &Javascript{engine:engine}
	case "coffeescript":
		engine.script = &Javascript{engine:engine}
	default:
		err = errors.New("undefined language")

	}

	engine.script.Init()

	return
}

type Engine struct {
	model  *models.Script
	script Magic
	buf	[]string
}

func (s *Engine) Compile() error {
	return s.script.Compile()
}

func (s *Engine) Update() (err error) {

	if err = s.Compile(); err != nil {
		return
	}

	if err = models.UpdateScriptById(s.model); err != nil {
		return
	}

	return
}

func (s *Engine) PushStruct(name string, i interface{}) (int, error) {
	return s.script.PushStruct(name, i)
}

func (s *Engine) PushFunction(name string, i interface{}) (int, error) {
	return s.script.PushFunction(name, i)
}

func (s *Engine) Close() {
	s.script.Close()
}

func (s *Engine) Do() (res string, err error) {
	var result string
	result, err = s.script.Do()
	for _, r := range s.buf {
		res += r + "\n"
	}

	res += result + "\n"

	return
}

func (s *Engine) Print(str string) {
	//TODO remove
	beego.Debug(str)
	s.buf = append(s.buf, str)
}
