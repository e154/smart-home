package notifr

import "github.com/e154/smart-home/api/log"

var (
	notifr *Notifr
)

type Notifr struct {

}

func NewNotifr() (notifr *Notifr) {

	notifr = &Notifr{}

	return
}

func (n *Notifr) SendEmail() {

}

func (n *Notifr) SendSMS() {

}

func (n *Notifr) SendPush() {

}

func Initialize() {
	log.Info("Notifr initialize...")

	if notifr != nil {
		return
	}

	notifr = NewNotifr()

	return
}