package scripts

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/log"
)

type Magic interface {
	Init() error
	Do() (string, error)
	Compile() error
	PushStruct(string, interface{}) (int, error)
	PushFunction(string, interface{}) (int, error)
	EvalString(string) (error)
	Close()
}

func New(s *models.Script) (engine *Engine, err error) {

	engine = &Engine{
		model:s,
		buf: make([]string, 0),
	}

	switch s.Lang {
	case "ts":
		engine.script = &Javascript{engine:engine}
	case "javascript":
		engine.script = &Javascript{engine:engine}
	case "coffeescript":
		engine.script = &Javascript{engine:engine}
	default:
		err = errors.New("undefined language")

	}

	if err == nil {
		log.Infof("Add script: %s (%s)", s.Name, s.Lang)
	}

	engine.script.Init()

	return
}

type Engine struct {
	model  *models.Script
	script Magic
	buf    []string
	IsRun  bool
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

func (s *Engine) EvalString(str string) error {
	return s.script.EvalString(str)
}

func (s *Engine) Close() {
	s.script.Close()
}

func (s *Engine) DoFull() (res string, err error) {
	if s.IsRun {
		return
	}

	s.IsRun = true
	var result string
	result, err = s.script.Do()
	for _, r := range s.buf {
		res += r + "\n"
	}

	res += result + "\n"

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Do() (result string, err error) {
	if s.IsRun {
		return
	}

	s.IsRun = true
	result, err = s.script.Do()

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Print(v ...interface{}) {
	beego.Debug(v...)
	s.buf = append(s.buf, fmt.Sprint(v...))
}

func (s *Engine) Get() Magic {
	return s.script
}