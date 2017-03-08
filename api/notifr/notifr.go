package notifr

import (
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
	"sync"
	"github.com/pkg/errors"
)

var (
	instantiated *Notifr = nil
)

type Notifr struct {
	mu sync.Mutex
	message_queue	[]*models.MessageDeliverie
}

func (n *Notifr) save(msg Message) (err error) {

	var message *models.Message
	if message, err = msg.save(); err != nil {
		log.Error("Notifr:", err.Error())
		return
	}

	// add Message Deliveries
	md := &models.MessageDeliverie{}
	md.Message = message
	md.State = "in_progress"

	if _, err = models.AddMessageDeliverie(md); err != nil {
		log.Error("Notifr:", err.Error())
		return
	}

	n.mu.Lock()
	n.message_queue = append(n.message_queue, md)
	n.mu.Unlock()

	err = n.worker()

	return
}

func (n *Notifr) send(msg *models.Message) (err error) {

	switch msg.Type {
	case "email":
		email := NewEmail()
		email.load(msg)
		err = email.send()
	case "sms":
	case "push":
	default:
		err = errors.New("Notifr: unknown message type")
	}

	return
}

func (n *Notifr) worker() (err error) {

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, md := range n.message_queue {
		if err = n.send(md.Message); err != nil {
			md.Error_system_message = err.Error()
		} else {
			md.State = "succeed"
			md.Error_system_code = ""
			md.Error_system_message = ""
		}

		if err = models.UpdateMessageDeliverieById(md); err != nil {
			return
		}
	}

	n.message_queue = []*models.MessageDeliverie{}

	return
}

func Send(msg Message) {
	go instantiated.save(msg)
}

func Initialize() {
	log.Info("Notifr initialize...")

	if instantiated != nil {
		return
	}

	instantiated = &Notifr{
		message_queue:	[]*models.MessageDeliverie{},
	}

	md, count, err := models.GetAllMessageDeliveriesInProgress()
	if err != nil {
		log.Error("Notifr:", err.Error())
	}

	instantiated.message_queue = md

	if count > 0 {
		instantiated.worker()
	}

	return
}