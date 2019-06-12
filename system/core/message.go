package core

func NewMessage() (m *Message) {
	m = &Message{}
	m.pull = make(map[string]interface{})
	return
}

type Message struct {
	Error string
	Storage
	Success   bool
	Direction bool
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
	m.pull = make(map[string]interface{})
	m.Error = ""
}

func (m *Message) Copy() (msg *Message) {
	msg = NewMessage()
	for k, v := range m.pull {
		msg.SetVar(k, v)
	}
	return
}
