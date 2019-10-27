package notify

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/op/go-logging"
	"time"
)

var (
	log = logging.MustGetLogger("notify")
)

const (
	notifyVarName = "notify"
)

type Notify struct {
	adaptor     *adaptors.Adaptors
	cfg         *NotifyConfig
	appCfg      *config.AppConfig
	stat        *NotifyStat
	isStarted   bool
	stopPrecess bool
	ticker      *time.Ticker
	workers     []*Worker
	queue       chan interface{}
	stopQueue   chan struct{}
}

func NewNotify(
	adaptor *adaptors.Adaptors,
	appCfg *config.AppConfig,
	graceful *graceful_service.GracefulService, ) *Notify {

	notify := &Notify{
		adaptor:   adaptor,
		appCfg:    appCfg,
		queue:     make(chan interface{}),
		stopQueue: make(chan struct{}),
	}

	graceful.Subscribe(notify)

	notify.Start()

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
	n.getCfg()

	if n.cfg == nil {
		n.initCfg()
	}

	n.isStarted = true

	// workers
	n.workers = []*Worker{
		NewWorker(n.cfg),
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

	log.Infof("Notify service started")
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

	n.workers = make([]*Worker, 0)

	n.isStarted = false

	log.Infof("Notify service stopped")
}

func (n *Notify) Restart() {
	n.stop()
	n.Start()
}

func (n Notify) Send(msg interface{}) {
	if !n.isStarted && n.stopPrecess {
		return
	}

	n.queue <- msg
}

func (n *Notify) getCfg() {

	v, err := n.adaptor.Variable.GetByName(notifyVarName)
	if err != nil {
		log.Error(err.Error())
		return
	}

	n.cfg = &NotifyConfig{}
	if err = json.Unmarshal([]byte(v.Value), n.cfg); err != nil {
		log.Error(err.Error())
	}
}

func (n *Notify) initCfg() {
	log.Infof("init config")

	n.cfg = &NotifyConfig{}
	n.UpdateCfg(n.cfg)
}

func (n *Notify) UpdateCfg(cfg *NotifyConfig) {

	b, err := json.Marshal(cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	variable := &m.Variable{
		Name:     notifyVarName,
		Value:    string(b),
		Autoload: false,
	}

	if err = n.adaptor.Variable.Update(variable); err != nil {
		log.Error(err.Error())
	}
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
