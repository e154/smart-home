package notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"strings"
)

type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"template"`
}

func NewEmail() (email *Email) {
	return &Email{}
}

func (e *Email) SetRender(render *m.TemplateRender) {
	e.Subject = render.Subject
	e.Body = render.Body
}

func (e *Email) Save() (addresses []string, message *m.Message) {
	e.To = strings.Replace(e.To, " ", "", -1)
	addresses = strings.Split(e.To, ",")
	message = &m.Message{
		Type:         m.MessageTypeEmail,
		EmailFrom:    common.String(e.From),
		EmailSubject: common.String(e.Subject),
		EmailBody:    common.String(e.Body),
	}
	return
}
