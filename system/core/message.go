package core

func NewMessage() (m *Message) {
	m = &Message{
		storage: NewStorage(),
	}
	return
}

type Message struct {
	Error     string
	storage   Storage
	Success   bool
	Direction bool
	Mqtt      bool
}

func (m *Message) clearError() {
	m.Error = ""
}

func (m *Message) SetError(err string) {
	m.Error = err
}

func (m *Message) Setdir(d bool) {
	m.Direction = d
}

func (m *Message) Ok() {
	m.Success = true
}

func (m *Message) Clear() {
	m.storage.pull = make(map[string]interface{})
	m.Error = ""
}

func (m *Message) Copy() (msg *Message) {
	msg = NewMessage()
	for k, v := range m.storage.pull {
		msg.storage.SetVar(k, v)
	}
	return
}

func (m *Message) GetVar(key string) (value interface{}) {
	return m.storage.GetVar(key)
}

func (m *Message) SetVar(key string, value interface{}) {
	m.storage.SetVar(key, value)
}
