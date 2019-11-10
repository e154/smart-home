package notify

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"strings"
)

type SMS struct {
	phones []string
	Text   string `json:"text"`
}

func NewSMS() (sms *SMS) {
	return &SMS{}
}

func (s *SMS) SetRender(render *m.TemplateRender) {
	s.Text = render.Body
}

func (s *SMS) AddPhone(phone string) {
	if !strings.Contains(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}
	s.phones = append(s.phones, phone)
}

func (s *SMS) Save() (addresses []string, message *m.Message) {

	addresses = s.phones
	message = &m.Message{
		Type:    m.MessageTypeSMS,
		SmsText: common.String(s.Text),
	}
	return
}
