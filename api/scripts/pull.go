package scripts

import (
	"sync"
	"github.com/e154/smart-home/api/scripts/bind"
)

var pull *Pull

type Pull struct {
	sync.Mutex
	functions	map[string]interface{}
	structures	map[string]interface{}
}

func PushStruct(name string, s interface{}) {
	pull.Lock()
	defer pull.Unlock()

	pull.structures[name] = s
}

func PushFunctions(name string, s interface{}) {
	pull.Lock()
	defer pull.Unlock()

	pull.functions[name] = s
}

func (p *Pull) GetStruct() map[string]interface{} {
	p.Lock()
	defer p.Unlock()

	return p.structures
}

func (p *Pull) Getfunctions() map[string]interface{} {
	p.Lock()
	defer p.Unlock()

	return p.functions
}

func init() {
	pull = &Pull{
		functions: make(map[string]interface{}),
		structures: make(map[string]interface{}),
	}

	PushStruct("Notifr", &bind.NotifrBind{})
	PushStruct("Log", &bind.LogBind{})
	PushFunctions("Execute", bind.Execute)
}