package notifr

import "github.com/e154/smart-home/api/models"

type Sms struct {
	To		string //"+79139131232,+79139131232,+79139131232"
	Body		string
}

func NewSms() (sms *Sms) {
	sms = &Sms{}
	return
}

func (e *Sms) save() (to string, msg *models.Message, err error) {

	msg = &models.Message{}
	msg.Sms_text = e.Body
	to = e.To

	_, err = models.AddMessage(msg)

	return
}

func (e *Sms) load(md *models.MessageDeliverie) {

	e.Body = md.Message.Sms_text
	e.To = md.Address

	return
}

func (e *Sms) send() (err error) {

	// send sms

	return
}