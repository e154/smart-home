package notify

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/scripts"
	"github.com/op/go-logging"
	"time"
)

var (
	log = logging.MustGetLogger("notify")
)

type Notify struct {
	adaptor       *adaptors.Adaptors
	cfg           *NotifyConfig
	appCfg        *config.AppConfig
	stat          *NotifyStat
	isStarted     bool
	stopPrecess   bool
	scriptService *scripts.ScriptService
	ticker        *time.Ticker
	workers       []*Worker
	queue         chan interface{}
	stopQueue     chan struct{}
}

func NewNotify(
	adaptor *adaptors.Adaptors,
	appCfg *config.AppConfig,
	graceful *graceful_service.GracefulService,
	scriptService *scripts.ScriptService) *Notify {

	notify := &Notify{
		adaptor:   adaptor,
		appCfg:    appCfg,
		queue:     make(chan interface{}),
		stopQueue: make(chan struct{}),
		cfg:       NewNotifyConfig(adaptor),
	}

	graceful.Subscribe(notify)

	notify.Start()

	scriptService.PushStruct("Notifr", &NotifyBind{
		notify: notify,
	})

	scriptService.PushStruct("Template", &TemplateBind{
		adaptor: adaptor,
	})

	return notify
}

func (n *Notify) Shutdown() {
	n.stop()
	close(n.stopQueue)
}

func (n *Notify) Start() {

	if n.isStarted {
		return
	}

	// update config
	n.cfg.Get()

	n.isStarted = true

	// workers
	n.workers = []*Worker{
		NewWorker(n.cfg, n.adaptor),
	}

	// stats
	n.ticker = time.NewTicker(time.Second * 5)

	n.updateStat()

	//...
	go func() {
		for ; ; {
			var worker *Worker
			for _, w := range n.workers {
				if w.inProcess {
					continue
				}
				worker = w
			}
			if worker == nil {
				time.Sleep(time.Millisecond * 500)
				continue
			}

			select {
			case msg := <-n.queue:
				worker.send(msg)
			case <-n.stopQueue:
				return
			}
		}
	}()

	for _, worker := range n.workers {
		worker.Start()
	}

	n.read()

	log.Infof("Notifr service started")
}

func (n *Notify) stop() {
	if n.stopPrecess {
		return
	}

	n.stopQueue <- struct{}{}

	n.stopPrecess = true
	defer func() {
		n.stopPrecess = false
	}()

	//...

	if n.ticker != nil {
		n.ticker.Stop()
	}
	n.ticker = nil

	for _, worker := range n.workers {
		worker.Stop()
	}
	n.workers = make([]*Worker, 0)

	n.isStarted = false

	log.Infof("Notifr service stopped")
}

func (n *Notify) Restart() {
	n.stop()
	n.Start()
}

func (n Notify) Send(msg interface{}) {
	if !n.isStarted && n.stopPrecess {
		return
	}

	switch v := msg.(type) {
	case IMessage:
		n.save(v)
	default:
		log.Errorf("unknown message type %v", v)
	}
}

func (n *Notify) save(t IMessage) {

	addresses, message := t.Save()

	messageId, err := n.adaptor.Message.Add(message)
	if err != nil {
		log.Error(err.Error())
		return
	}
	message.Id = messageId

	for _, address := range addresses {
		messageDelivery := &m.MessageDelivery{
			Message:   message,
			MessageId: message.Id,
			Status:    m.MessageStatusInProgress,
			Address:   address,
		}
		if messageDelivery.Id, err = n.adaptor.MessageDelivery.Add(messageDelivery); err != nil {
			log.Error(err.Error())
		}

		n.queue <- messageDelivery
	}
}

func (n *Notify) read() {

	messageDeliveries, _, err := n.adaptor.MessageDelivery.GetAllUncompleted(99, 0)
	if err != nil {
		log.Error(err.Error())
	}

	for _, msg := range messageDeliveries {
		n.queue <- msg
	}
}

func (n *Notify) GetCfg() *NotifyConfig {
	return n.cfg
}

func (n *Notify) UpdateCfg(cfg *NotifyConfig) error {
	return n.cfg.Update()
}

func (n *Notify) Stat() *NotifyStat {
	return n.stat
}

func (n *Notify) updateStat() {

	stat := &NotifyStat{
		Workers: len(n.workers),
	}

	if stat.Workers == 0 {
		return
	}

	worker := n.workers[0]

	// messagebird balance
	if worker.mbClient != nil {
		if mbBalance, err := worker.mbClient.Balance(); err == nil {
			stat.MbBalance = mbBalance.Amount
		}
	}

	// twilio balance
	if worker.twClient != nil {
		if twBalance, err := worker.twClient.Balance(); err == nil {
			stat.TwBalance = twBalance
		} else {
			log.Error(err.Error())
		}
	}

	n.stat = stat
}

func (n *Notify) Repeat(msg *m.MessageDelivery) {
	if !n.isStarted && n.stopPrecess {
		return
	}

	msg.Status = m.MessageStatusInProgress
	_ = n.adaptor.MessageDelivery.SetStatus(msg)

	n.queue <- msg
}
