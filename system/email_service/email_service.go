package email_service

import (
	"errors"
	"github.com/op/go-logging"
	"gopkg.in/gomail.v2"
)

var (
	log = logging.MustGetLogger("email")
)

type EmailService struct {
	cfg *EmailServiceConfig
}

func NewEmailService(cfg *EmailServiceConfig) (*EmailService, error) {

	if cfg.Auth == "" || cfg.Pass == "" || cfg.Smtp == "" || cfg.Port == 0 ||
		cfg.Sender == "" {
		return nil, errors.New("bad parameters")
	}

	client := &EmailService{
		cfg: cfg,
	}
	return client, nil
}

func (e EmailService) Send(email *Email) error {

	email.From = e.cfg.Sender

	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":     {email.From},
		"Reply-To": {email.From},
		"To":       {email.To},
		"Subject":  {email.Subject},
	})

	m.SetBody("text/html", email.Body)

	d := gomail.NewPlainDialer(e.cfg.Smtp, e.cfg.Port, e.cfg.Auth, e.cfg.Pass)
	if err := d.DialAndSend(m); err != nil {
		return errors.New(err.Error())
	}

	log.Debug("Sent email '" + email.Subject + "' to:" + email.To)

	return nil
}
