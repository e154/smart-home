package scripts

func (j *JavascriptBinding) GetVariable(key string) interface{} {

	j.mu.Lock()
	defer j.mu.Unlock()

	if v, ok := j.pool[key]; ok {
		return v
	}

	return nil
}

func (j *JavascriptBinding) SetVariable(key string, value interface{}) {

	j.mu.Lock()
	defer j.mu.Unlock()

	j.pool[key] = value
}

func (j *JavascriptBinding) setVariablePool(pool map[string]interface{}) {

	j.mu.Lock()
	defer j.mu.Unlock()

	j.pool = pool
}