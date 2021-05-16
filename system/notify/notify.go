// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package notify

import (
	"context"
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/fx"
	"time"
)

var (
	log = common.MustGetLogger("notify")
)

type Notify interface {
	Shutdown() error
	Start() (err error)
	Restart()
	GetCfg() *Config
	UpdateCfg(cfg *Config) error
	Stat() *NotifyStat
	Repeat(msg *m.MessageDelivery)
	Send(msg interface{})
}

// notify ...
type notify struct {
	adaptor       *adaptors.Adaptors
	cfg           *Config
	appCfg        *config.AppConfig
	stat          *NotifyStat
	isStarted     bool
	stopPrecess   bool
	scriptService scripts.ScriptService
	ticker        *time.Ticker
	workers       []*Worker
	queue         chan interface{}
	stopQueue     chan struct{}
}

// NewNotify ...
func NewNotify(
	lc fx.Lifecycle,
	adaptor *adaptors.Adaptors,
	appCfg *config.AppConfig,
	scriptService scripts.ScriptService) Notify {

	notify := &notify{
		adaptor:   adaptor,
		appCfg:    appCfg,
		queue:     make(chan interface{}),
		stopQueue: make(chan struct{}),
		cfg:       NewConfig(adaptor),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return notify.Start()
		},
		OnStop: func(ctx context.Context) (err error) {
			return notify.Shutdown()
		},
	})

	scriptService.PushStruct("Notifr", &NotifyBind{
		notify: notify,
	})

	scriptService.PushStruct("Template", &TemplateBind{
		adaptor: adaptor,
	})

	return notify
}

// Shutdown ...
func (n *notify) Shutdown() error {
	n.stop()
	close(n.stopQueue)
	return nil
}

// Start ...
func (n *notify) Start() (err error) {

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
		for {
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

	return
}

func (n *notify) stop() {
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

// Restart ...
func (n *notify) Restart() {
	n.stop()
	n.Start()
}

// Send ...
func (n notify) Send(msg interface{}) {
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

func (n *notify) save(t IMessage) {

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

func (n *notify) read() {

	messageDeliveries, _, err := n.adaptor.MessageDelivery.GetAllUncompleted(99, 0)
	if err != nil {
		log.Error(err.Error())
	}

	for _, msg := range messageDeliveries {
		n.queue <- msg
	}
}

func (n *notify) getCfg() {

	v, err := n.adaptor.Variable.GetByName(notifyVarName)
	if err != nil {
		log.Error(err.Error())
		return
	}

	n.cfg = &Config{}
	if err = json.Unmarshal([]byte(v.Value), n.cfg); err != nil {
		log.Error(err.Error())
	}
}

// GetCfg ...
func (n *notify) GetCfg() *Config {
	return n.cfg
}

// UpdateCfg ...
func (n *notify) UpdateCfg(cfg *Config) error {
	cfg.adaptor = n.adaptor
	n.cfg = cfg
	return n.cfg.Update()
}

// Stat ...
func (n *notify) Stat() *NotifyStat {
	return n.stat
}

func (n *notify) updateStat() {

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

// Repeat ...
func (n *notify) Repeat(msg *m.MessageDelivery) {
	if !n.isStarted && n.stopPrecess {
		return
	}

	msg.Status = m.MessageStatusInProgress
	_ = n.adaptor.MessageDelivery.SetStatus(msg)

	n.queue <- msg
}
