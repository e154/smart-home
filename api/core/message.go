package core

func NewMessage() *Message {
	return &Message{
		Pull: make(map[string]interface{}),
	}
}
//TODO refactor message system
type Message struct {
	Error       	string
	Pull			map[string]interface{}
}
//TODO refactor message system
func (m *Message) clearError() {
	m.Error = ""
}
//TODO refactor message system
func (m *Message) SetError(err string) {
	m.Error = err
}

func (m *Message) SetVar(key string, val interface{}) {
	m.Pull[key] = val
}

func (m *Message) GetVar(key string) interface{} {
	if value, ok := m.Pull[key]; ok {
		return value
	}

	return nil
}

func (m *Message) Clear() {
	m.Pull = make(map[string]interface{})
	m.Error = ""
}