package notify

import (
	"github.com/e154/smart-home/system/email_service"
	mb "github.com/e154/smart-home/system/messagebird"
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
	inProcess      bool
	isStarted      bool
}

func NewWorker(cfg *NotifyConfig) *Worker {

	worker := &Worker{
		cfg: cfg,
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
	telegramClient := telegram.NewTelegramConfig(n.cfg.TelegramToken)
	if telegramClient, err := telegram.NewTelegram(telegramClient); err == nil {
		n.telegramClient = telegramClient
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

func (n *Worker) send(msg interface{}) {

	n.inProcess = true

	switch v := msg.(type) {
	case *SMS:
		n.sendSms(v)
	default:
		log.Errorf("unknown message type %v", v)
	}

	n.inProcess = false
}

func (n *Worker) sendSms(msg *SMS) {

	msgId, err := n.twClient.SendSMS(msg.Phone, msg.Text)
	if err != nil {
		log.Error(err.Error())
	}

	time.Sleep(25 * time.Second)

	var status string
	if status, err = n.twClient.GetStatus(msgId); err != nil {
		log.Error(err.Error())
	}

	if status == tw.StatusDelivered {
		return
	}

	if n.mbClient != nil {
		if _, err := n.mbClient.SendSMS(msg.Phone, msg.Text); err != nil {
			log.Error(err.Error())
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

func (n *Worker) sendEmail(msg *Email) {

	if n.emailClient == nil {
		return
	}

	email := &email_service.Email{
		From:     msg.From,
		To:       msg.To,
		Subject:  msg.Subject,
		Template: msg.Template,
		Data:     msg.Data,
	}

	if err := n.emailClient.Send(email); err != nil {
		log.Error(err.Error())
	}
}
