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
}

func NewWorker(cfg *NotifyConfig) *Worker {

	worker := &Worker{
		cfg: cfg,
	}

	// messagebird
	mbConfig := mb.NewMBClientConfig(cfg.MbAccessKey, cfg.MbName)
	if mbClient, err := mb.NewMBClient(mbConfig); err == nil {
		worker.mbClient = mbClient
	}

	// twilio
	twConfig := tw.NewTWConfig(cfg.TWFrom, cfg.TWSid, cfg.TWAuthToken)
	if twClient, err := tw.NewTWClient(twConfig); err == nil {
		worker.twClient = twClient
	}

	// email
	emailConfig := email_service.NewEmailServiceConfig(cfg.EmailAuth, cfg.EmailPass, cfg.EmailSmtp, cfg.EmailPort, cfg.EmailSender)
	if emailClient, err := email_service.NewEmailService(emailConfig); err == nil {
		worker.emailClient = emailClient
	}

	// telegram
	telegramClient := telegram.NewTelegramConfig(cfg.TelegramToken)
	if telegramClient, err := telegram.NewTelegram(telegramClient); err == nil {
		worker.telegramClient = telegramClient
	}

	return worker
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
