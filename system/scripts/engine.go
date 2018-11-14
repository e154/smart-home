package scripts

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"io/ioutil"
)

type Magic interface {
	Init() error
	Do() (string, error)
	DoCustom(string) (string, error)
	Compile() error
	PushStruct(string, interface{}) (int, error)
	PushFunction(string, interface{}) (int, error)
	EvalString(string) (error)
	Close()
}

type Engine struct {
	model  *m.Script
	script Magic
	buf    []string
	IsRun  bool
	pull   *Pull
}

func (s *Engine) Compile() error {
	return s.script.Compile()
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

func (s *Engine) Do() (string, error) {
	return s.script.Do()
}

func (s *Engine) DoCustom(f string) (result string, err error) {

	if s.IsRun {
		return
	}

	s.IsRun = true
	result, err = s.script.DoCustom(f)

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Print(v ...interface{}) {
	fmt.Println(v...)
	s.buf = append(s.buf, fmt.Sprint(v...))
}

func (s *Engine) Get() Magic {
	return s.script
}

func (s *Engine) File(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}