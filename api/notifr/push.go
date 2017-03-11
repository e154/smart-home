package notifr

import "github.com/e154/smart-home/api/models"

type Push struct {
	To		string
	Body		string
}

func NewPush() (push *Push) {
	push = &Push{}
	return
}

func (e *Push) save() (to string, msg *models.Message, err error) {

	msg = &models.Message{}
	msg.Ui_text = e.Body
	to = e.To

	_, err = models.AddMessage(msg)

	return
}

func (e *Push) load(md *models.MessageDeliverie) {

	e.Body = md.Message.Ui_text
	e.To = md.Address

	return
}

func (e *Push) send() (err error) {

	// send push

	return
}