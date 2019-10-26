package core

import (
	"sync"
)

type Storage struct {
	mx   *sync.Mutex
	pull map[string]interface{}
}

func NewStorage() Storage {
	return Storage{
		mx:   &sync.Mutex{},
		pull: make(map[string]interface{}),
	}
}

func (s *Storage) GetVar(key string) (value interface{}) {

	s.mx.Lock()
	if v, ok := s.pull[key]; ok {
		value = v
	} else {
		value = nil
	}
	s.mx.Unlock()
	return
}

func (s *Storage) SetVar(key string, value interface{}) {

	s.mx.Lock()
	s.pull[key] = value
	s.mx.Unlock()
}
