package notify

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/email_service"
	mb "github.com/e154/smart-home/system/messagebird"
	"github.com/e154/smart-home/system/slack"
	"github.com/e154/smart-home/system/telegram"
	tw "github.com/e154/smart-home/system/twilio"
	"time"
)

type Worker struct {
	cfg            *NotifyConfig
	mbClient       *mb.MBClient
	twClient       *tw.TWClient
	emailClient    *email_service.EmailService
	telegramClient *telegram.Telegram
	slackClient    *slack.Slack
	adaptor        *adaptors.Adaptors
	inProcess      bool
	isStarted      bool
}

func NewWorker(cfg *NotifyConfig,
	adaptor *adaptors.Adaptors) *Worker {

	worker := &Worker{
		cfg:     cfg,
		adaptor: adaptor,
	}

	return worker
}

func (n *Worker) Start() {

	if n.isStarted {
		return
	}

	// messagebird
	mbConfig := mb.NewMBClientConfig(n.cfg.MbAccessKey, n.cfg.MbName)
	if mbClient, err := mb.NewMBClient(mbConfig); err == nil {
		n.mbClient = mbClient
	}

	// twilio
	twConfig := tw.NewTWConfig(n.cfg.TWFrom, n.cfg.TWSid, n.cfg.TWAuthToken)
	if twClient, err := tw.NewTWClient(twConfig); err == nil {
		n.twClient = twClient
	}

	// email
	emailConfig := email_service.NewEmailServiceConfig(n.cfg.EmailAuth, n.cfg.EmailPass, n.cfg.EmailSmtp, n.cfg.EmailPort, n.cfg.EmailSender)
	if emailClient, err := email_service.NewEmailService(emailConfig); err == nil {
		n.emailClient = emailClient
	}

	// telegram
	telegramConfig := telegram.NewTelegramConfig(n.cfg.TelegramToken)
	if telegramClient, err := telegram.NewTelegram(telegramConfig); err == nil {
		n.telegramClient = telegramClient
	}

	// slack
	slackConfig := slack.NewSlackConfig(n.cfg.SlackToken, n.cfg.SlackUserName)
	if slackClient, err := slack.NewSlack(slackConfig); err == nil {
		n.slackClient = slackClient
	}

	n.isStarted = true
}

func (n *Worker) Stop() {
	if !n.isStarted {
		return
	}
	if n.telegramClient != nil {
		n.telegramClient.Stop()
	}

	n.isStarted = false
}

func (n *Worker) sendMessageDelivery(msg *m.MessageDelivery) {
	switch msg.Message.Type {
	case m.MessageTypeEmail:
		go n.sendEmail(msg)
	case m.MessageTypeSMS:
		go n.sendSms(msg)
	case m.MessageTypeSlack:
		go n.sendSlack(msg)
	default:
		log.Errorf("unknown message type %v", msg.Message.Type)
	}
}

func (n *Worker) send(msg interface{}) {

	n.inProcess = true

	switch v := msg.(type) {
	case *m.MessageDelivery:
		n.sendMessageDelivery(v)
	default:
		log.Errorf("unknown message type %v", v)
	}

	n.inProcess = false
}

func (n *Worker) sendSms(msg *m.MessageDelivery) {

	text := *msg.Message.SmsText

	if n.twClient != nil {
		msgId, err := n.twClient.SendSMS(msg.Address, text)
		if err != nil {
			n.setError(msg, err)
		}

		ticker := time.NewTicker(time.Second)

		var status string

		i := 0
		for range ticker.C {
			if i > 15 {
				ticker.Stop()
			}
			if status, err = n.twClient.GetStatus(msgId); err != nil {
				n.setError(msg, err)
			}
			if status == tw.StatusDelivered {
				n.setSucceed(msg)
				ticker.Stop()
				return
			}
			i++
		}
	}

	if n.mbClient != nil {
		msgId, err := n.mbClient.SendSMS(msg.Address, text)
		if err != nil {
			n.setError(msg, err)
			return
		}

		ticker := time.NewTicker(time.Second)

		var status string

		i := 0
		for range ticker.C {
			if i > 15 {
				ticker.Stop()
			}
			if status, err = n.mbClient.GetStatus(msgId); err != nil {
				n.setError(msg, err)
				ticker.Stop()
				return
			}
			if status == mb.StatusDelivered {
				n.setSucceed(msg)
				ticker.Stop()
				return
			}
			i++
		}
	}
}

func (n *Worker) sendTelegram(msg *Telegram) {

	if n.telegramClient == nil {
		return
	}

	if err := n.telegramClient.SendMsg(msg.Text, msg.Channel); err != nil {
		log.Error(err.Error())
	}
}

func (n *Worker) sendSlack(msg *m.MessageDelivery) {

	if n.slackClient == nil {
		return
	}

	slackMessage := &slack.SlackMessage{
		Text:    common.StringValue(msg.Message.SlackText),
		Channel: msg.Address,
	}
	if err := n.slackClient.SendMsg(slackMessage); err != nil {
		log.Error(err.Error())
	}
}

func (n *Worker) sendEmail(msg *m.MessageDelivery) {

	if n.emailClient == nil {
		return
	}

	email := &email_service.Email{
		From:    common.StringValue(msg.Message.EmailFrom),
		To:      msg.Address,
		Subject: common.StringValue(msg.Message.EmailSubject),
		Body:    common.StringValue(msg.Message.EmailBody),
	}

	if err := n.emailClient.Send(email); err != nil {
		n.setError(msg, err)
		return
	}

	n.setSucceed(msg)
}

func (n *Worker) setSucceed(msg *m.MessageDelivery) {
	msg.Status = m.MessageStatusSucceed
	_ = n.adaptor.MessageDelivery.SetStatus(msg)
}

func (n *Worker) setError(msg *m.MessageDelivery, err error) {
	msg.Status = m.MessageStatusError
	msg.ErrorMessageBody = common.String(err.Error())
	_ = n.adaptor.MessageDelivery.SetStatus(msg)
}
