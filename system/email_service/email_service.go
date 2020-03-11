// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package email_service

import (
	"errors"
	"github.com/e154/smart-home/system/logging"
	"gopkg.in/gomail.v2"
)

var (
	log = common.MustGetLogger("email")
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
