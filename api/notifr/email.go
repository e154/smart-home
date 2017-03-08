package notifr

import (
	"github.com/e154/smart-home/api/models"
	"errors"
	"net/mail"
	"strings"
	"github.com/astaxie/beego"
	"fmt"
	"encoding/base64"
	"net/smtp"
	"bytes"
	"io"
)

var (
	email_user string
	email_password string
	email_smtp string
	email_from_name string
	email_from_address string
)

type Email struct {
	To       string //"Alice <alice@example.com>, Bob <bob@example.com>, Eve <eve@example.com>"
	Subject  string
	Body     string
	Template string
	Params   map[string]interface{}
	rendered *models.EmailRender
}

func NewEmail() (email *Email) {
	email = &Email{
		Params: make(map[string]interface{}),
	}

	return
}

func (e *Email) render() (err error) {

	if e.Template == "" {
		err = errors.New("template field is nil")
		return
	}

	params := map[string]interface {}{

	}

	for k, n := range e.Params {
		params[k] = n
	}

	var rendered *models.EmailRender
	if rendered, err = models.EmailTemplateRender(e.Template, params); err != nil {
		return
	}

	e.rendered = rendered

	return
}

func (e *Email) save() (msg *models.Message, err error) {

	if e.rendered == nil {
		if err = e.render(); err != nil {
			return
		}
	}

	msg = &models.Message{}
	msg.Type = "email"
	msg.EmailBody = e.rendered.Body
	msg.EmailTitle = e.rendered.Subject
	msg.To = e.To

	_, err = models.AddMessage(msg)

	return
}

func (e *Email) load(msg *models.Message) {

	e.rendered = &models.EmailRender{
		Subject:msg.EmailTitle,
		Body:msg.EmailBody,
	}

	e.To = msg.To
}

func (e *Email) send() (err error) {

	if e.rendered == nil {
		if err = e.render(); err != nil {
			return
		}
	}

	from := mail.Address{email_from_name, email_from_address}
	var to []*mail.Address
	if to, err = e.parseAddressList(); err != nil {
		return
	}

	encodeRFC2047 := func(String string) string{
		// use mail's rfc2047 to encode any string
		addr := mail.Address{String, ""}
		return strings.Trim(addr.String(), " <>")
	}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = e.To
	header["Subject"] = encodeRFC2047(e.rendered.Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(e.rendered.Body))

	var addresses []string
	for _, v := range to {
		addresses = append(addresses, v.Address)
	}

	if email_user != "" && email_password != "" {

		auth := smtp.PlainAuth(
			"",
			email_user,
			email_password,
			email_smtp,
		)

		err = smtp.SendMail(
			email_smtp + ":25",
			auth,
			from.Address,
			addresses,
			[]byte(message),
		)

	} else {

		// unsecure mode
		for _, address := range addresses {
			var client *smtp.Client
			if client, err = smtp.Dial(email_smtp + ":25"); err != nil { return err }
			defer client.Close()
			if err = client.Mail(from.Address); err != nil { return err }

			if err = client.Rcpt(address); err != nil { return err }

			var w io.WriteCloser
			if w, err = client.Data(); err != nil { return err }
			defer w.Close()

			buf := bytes.NewBufferString(message)
			if _, err = buf.WriteTo(w); err != nil { return err }

			w.Close()
			client.Close()
		}

	}

	return
}

func (e *Email) parseAddressList() (emails []*mail.Address, err error) {

	if emails, err = mail.ParseAddressList(e.To); err != nil {
		return
	}

	return
}

func init()  {
	email_user = beego.AppConfig.String("email_user")
	email_password = beego.AppConfig.String("email_password")
	email_smtp = beego.AppConfig.String("email_smtp")
	email_from_name = beego.AppConfig.String("email_from_name")
	email_from_address = beego.AppConfig.String("email_from_address")
}