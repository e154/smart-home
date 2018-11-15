package core

func NewMessage() (m *Message) {
	m = &Message{}
	m.pull = make(map[string]interface{})
	return
}

type Message struct {
	Error string
	Storage
}

func (m *Message) clearError() {
	m.Error = ""
}

func (m *Message) SetError(err string) {
	m.Error = err
}

func (m *Message) Clear() {
	m.pull = make(map[string]interface{})
	m.Error = ""
}
