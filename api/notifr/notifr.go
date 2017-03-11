package notifr

import (
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
	"sync"
	"github.com/pkg/errors"
	"strings"
	"github.com/astaxie/beego/orm"
)

// email := notifr.NewEmail()
// email.To = "Alice <support@e154.ru>, Алексей <alll80@mail.ru>"
// email.Template = "password_reset"
// email.Params["user:name:last"] = "last"
// email.Params["user:name:first"] = "first"
// email.Params["site:name"] = "Smart home"
// email.Params["user:one-time-login-url"] = ""
// notifr.Send(email)
//
// email = notifr.NewEmail()
// email.Subject = "test message"
// email.Body = "test message body"
// email.To = "Alice <support@e154.ru>"
// email.Params["user:name:last"] = "last"
// email.Params["user:name:first"] = "first"
// email.Params["site:name"] = "Smart home"
// email.Params["user:one-time-login-url"] = ""
// notifr.Send(email)
// email := notifr.NewEmail()
// email.To = "Alice <support@e154.ru>, Алексей <alll80@mail.ru>"
// email.Template = "password_reset"
// email.Params["user:name:last"] = "last"
// email.Params["user:name:first"] = "first"
// email.Params["site:name"] = "Smart home"
// email.Params["user:one-time-login-url"] = ""
// notifr.Send(email)
//
// email = notifr.NewEmail()
// email.Subject = "test message"
// email.Body = "test message body"
// email.To = "Alice <support@e154.ru>"
// email.Params["user:name:last"] = "last"
// email.Params["user:name:first"] = "first"
// email.Params["site:name"] = "Smart home"
// email.Params["user:one-time-login-url"] = ""
// notifr.Send(email)

var (
	instantiated *Notifr = nil
)

type Notifr struct {
	mu sync.Mutex
	message_queue	[]*models.MessageDeliverie
}

func (n *Notifr) save(msg Message) (err error) {

	var message *models.Message
	var to string
	if to, message, err = msg.save(); err != nil {
		log.Error("Notifr:", err.Error())
		return
	}

	// add Message Deliveries
	var mds []*models.MessageDeliverie
	addresses := strings.Split(to, ",")
	for _, address := range addresses {
		md := &models.MessageDeliverie{}
		md.Message = message
		md.State = "in_progress"
		md.Address = strings.TrimSpace(address)
		mds = append(mds, md)
	}

	if _, _errors := models.AddMessageDeliverieMultiple(mds); len(_errors) != 0 {
		for _, err := range _errors {
			log.Error("Notifr:", err.Error())
		}
		return
	}

	n.mu.Lock()
	for _, md := range mds {
		n.message_queue = append(n.message_queue, md)
	}
	n.mu.Unlock()

	err = n.worker()

	return
}

func (n *Notifr) send(md *models.MessageDeliverie) (err error) {

	switch md.Message.Type {
	case "email":
		email := NewEmail()
		email.load(md)
		err = email.send()
	case "sms":
		sms := NewSms()
		sms.load(md)
		err = sms.send()
	case "push":
		push := NewPush()
		push.load(md)
		err = push.send()
	default:
		err = errors.New("Notifr: unknown message type")
	}

	return
}

func (n *Notifr) repeat(id int64) (err error) {

	var md *models.MessageDeliverie
	if md, err = models.GetMessageDeliverieById(id); err != nil {
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(md, "Message"); err != nil {
		return
	}

	if err = n.send(md); err != nil {
		md.Error_system_message = err.Error()
		md.State = "error"
	}

	if err = models.UpdateMessageDeliverieById(md); err != nil {
		return
	}

	return
}

func (n *Notifr) worker() (err error) {

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, md := range n.message_queue {
		if err = n.send(md); err != nil {
			md.Error_system_message = err.Error()
			md.State = "error"
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

func RepeatById(id int64) {
	go func() {
		if err := instantiated.repeat(id); err != nil {
			log.Error("Notifr:", err.Error())
		}
	}()
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