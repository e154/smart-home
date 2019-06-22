package gate_client

import (
	"encoding/json"
	"github.com/e154/smart-home/system/uuid"
)

type IWsCallback interface {
	onMessage(payload []byte)
	onConnected()
	onClosed()
}

const (
	ClientTypeServer = "server"
)

const (
	Request       = "request"
	Response      = "response"
	StatusSuccess = "success"
	StatusError   = "error"
)

type Message struct {
	Id      uuid.UUID              `json:"id"`
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
	Forward string                 `json:"forward"`
	Status  string                 `json:"status"`
}

func NewMessage(b []byte) (message Message, err error) {

	message = Message{}
	err = json.Unmarshal(b, &message)

	return
}

func (m *Message) Pack() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *Message) Response(payload map[string]interface{}) *Message {
	msg := &Message{
		Id:      m.Id,
		Payload: payload,
		Forward: Response,
	}
	return msg
}

func (m *Message) Success() *Message {
	msg := &Message{
		Id:      m.Id,
		Payload: map[string]interface{}{},
		Forward: Response,
		Status:  StatusSuccess,
	}
	return msg
}

func (m *Message) Error(err error) *Message {
	msg := &Message{
		Id: m.Id,
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
		Forward: Response,
		Status:  StatusError,
	}
	return msg
}
