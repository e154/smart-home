package notifr

import (
	"github.com/e154/smart-home/api/log"
)

var (
	instantiated *Notifr = nil
)

type Notifr struct {

}

func (n *Notifr) send(msg Message) error {
	return msg.send()
}

func Send(msg Message) error {
	return instantiated.send(msg)
}

func Initialize() {
	log.Info("Notifr initialize...")

	if instantiated != nil {
		return
	}

	instantiated = &Notifr{}

	return
}