package notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type Telegram struct {
	Text string `json:"text"`
}

func NewTelegram(text string) *Telegram {
	return &Telegram{
		Text: text,
	}
}

func (s *Telegram) SetRender(render *m.TemplateRender) {
	s.Text = render.Body
}

func (s *Telegram) Save() (addresses []string, message *m.Message) {

	addresses = []string{""}
	message = &m.Message{
		Type:         m.MessageTypeTelegramNotify,
		TelegramText: common.String(s.Text),
	}
	return
}
