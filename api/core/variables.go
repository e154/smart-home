package core

import "sync"

type Variables struct {
	pool map[string]interface{}
	mu   sync.Mutex
}

func NewVariablePool() *Variables {
	return &Variables{
		pool: make(map[string]interface{}),
	}
}

func (b *Variables) GetVariable(key string) interface{} {

	b.mu.Lock()
	defer b.mu.Unlock()

	if v, ok := b.pool[key]; ok {
		return v
	}

	return nil
}

func (b *Variables) SetVariable(key string, value interface{}) {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.pool[key] = value
}