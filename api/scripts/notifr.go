package scripts

import (
	"github.com/e154/smart-home/api/notifr"
)

type Notifr struct {

}

func (n *Notifr) NewEmail() *notifr.Email {
	return notifr.NewEmail()
}

func (n *Notifr) NewSms() *notifr.Sms {
	return notifr.NewSms()
}

func (n *Notifr) NewPush() *notifr.Push {
	return notifr.NewPush()
}

func (n *Notifr) Send(msg notifr.Message) {
	notifr.Send(msg)
}