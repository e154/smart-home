package notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func NewSlackMessage(channel, text string) *SlackMessage {
	return &SlackMessage{Channel: channel, Text: text}
}

func (s *SlackMessage) SetRender(render *m.TemplateRender) {
	s.Text = render.Body
}

func (s *SlackMessage) Save() (addresses []string, message *m.Message) {

	addresses = []string{s.Channel}
	message = &m.Message{
		Type:      m.MessageTypeSlack,
		SlackText: common.String(s.Text),
	}
	return
}
