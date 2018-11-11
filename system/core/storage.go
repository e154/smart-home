package core

import "sync"

type Storage struct {
	sync.Mutex
	pull map[string]interface{}
}

func (s *Storage) GetVar(key string) interface{} {

	s.Lock()
	defer s.Unlock()

	if v, ok := s.pull[key]; ok {
		return v
	}

	return nil
}

func (s *Storage) SetVar(key string, value interface{}) {

	s.Lock()
	defer s.Unlock()

	s.pull[key] = value
}
