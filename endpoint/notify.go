package endpoint

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/notify"
	"strings"
)

type NotifyEndpoint struct {
	*CommonEndpoint
}

func NewNotifyEndpoint(common *CommonEndpoint) *NotifyEndpoint {
	return &NotifyEndpoint{
		CommonEndpoint: common,
	}
}

func (n *NotifyEndpoint) GetSettings() (cfg *notify.NotifyConfig, err error) {
	cfg = n.notify.GetCfg()
	return
}

func (n *NotifyEndpoint) UpdateSettings(cfg *notify.NotifyConfig) (err error) {
	if err = n.notify.UpdateCfg(cfg); err != nil {
		return
	}

	n.notify.Restart()
	return
}

func (n *NotifyEndpoint) Repeat(id int64) (err error) {

	var msg *m.MessageDelivery
	if msg, err = n.adaptors.MessageDelivery.GetById(id); err != nil {
		return
	}

	n.notify.Repeat(msg)

	return
}

func (n *NotifyEndpoint) Send(params *m.NewNotifrMessage) (err error) {

	var render *m.TemplateRender
	if params.BodyType == "template" && params.Template != nil && params.Params != nil {
		if render, err = n.adaptors.Template.Render(common.StringValue(params.Template), params.Params); err != nil {
			return
		}
	}

	switch params.Type {
	case "email":
		message := notify.NewEmail()
		message.From = common.StringValue(params.EmailFrom)
		message.Subject = common.StringValue(params.EmailSubject)
		message.Body = common.StringValue(params.EmailBody)
		message.To = params.Address

		if render != nil {
			message.SetRender(render)
		}
		n.notify.Send(message)
	case "sms":

		message := notify.NewSMS()
		message.Text = common.StringValue(params.SmsText)

		for _, address := range strings.Split(params.Address, ",") {
			phone := strings.Replace(address, " ", "", -1)
			if phone == "" {
				continue
			}
			message.AddPhone(phone)
		}
		if render != nil {
			message.SetRender(render)
		}
		n.notify.Send(message)

	case "slack":
		message := notify.NewSlackMessage(common.StringValue(params.SlackText), params.Address)
		if render != nil {
			message.SetRender(render)
		}
		n.notify.Send(message)
	case "telegram_notify":
		message := notify.NewTelegram(common.StringValue(params.TelegramText))
		if render != nil {
			message.SetRender(render)
		}
		n.notify.Send(message)
	}

	return
}
