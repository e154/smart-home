package scripts

import (
	"sync"
)

type Pull struct {
	sync.Mutex
	functions  map[string]interface{}
	structures map[string]interface{}
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
