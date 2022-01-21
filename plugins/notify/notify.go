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
	"sync"
	"time"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

const (
	queueSize = 30
)

// Notify ...
type Notify interface {
	Shutdown() error
	Start() (err error)
	Stat() *Stat
	Repeat(msg m.MessageDelivery)
	Send(msg Message)
	AddProvider(name string, provider Provider)
	RemoveProvider(name string)
	Provider(name string) (provider Provider, err error)
}

// notify ...
type notify struct {
	adaptor       *adaptors.Adaptors
	stat          *Stat
	isStarted     *atomic.Bool
	scriptService scripts.ScriptService
	workers       []*Worker
	queue         chan m.MessageDelivery
	providerMu    *sync.RWMutex
	providerList  map[string]Provider
}

// NewNotify ...
func NewNotify(
	adaptor *adaptors.Adaptors,
	scriptService scripts.ScriptService) Notify {

	n := &notify{
		adaptor:      adaptor,
		isStarted:    atomic.NewBool(false),
		queue:        make(chan m.MessageDelivery, queueSize),
		providerMu:   &sync.RWMutex{},
		providerList: make(map[string]Provider),
	}

	// workers
	n.workers = []*Worker{
		NewWorker(adaptor),
		NewWorker(adaptor),
		NewWorker(adaptor),
	}

	scriptService.PushStruct("notifr", NewNotifyBind(n))
	scriptService.PushStruct("template", NewTemplateBind(adaptor))

	return n
}

// Shutdown ...
func (n *notify) Shutdown() error {
	n.stop()
	return nil
}

// Start ...
func (n *notify) Start() (err error) {

	if n.isStarted.Load() {
		return
	}
	n.isStarted.Store(true)

	n.updateStat()

	go func() {

		defer func() {
			n.workers = make([]*Worker, 0)
		}()

		for {
			var worker *Worker
			for _, w := range n.workers {
				if w.InWork() {
					continue
				}
				worker = w
			}
			if worker == nil {
				time.Sleep(time.Millisecond * 500)
				continue
			}

			for event := range n.queue {
				provider, err := n.Provider(event.Message.Type)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				worker.send(event, provider)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	n.read()

	return
}

func (n *notify) stop() {

	close(n.queue)
	n.isStarted.Store(false)
}

// Send ...
func (n notify) Send(msg Message) {

	if !n.isStarted.Load() {
		return
	}

	n.save(msg)
}

func (n *notify) save(event Message) {

	provider, err := n.Provider(event.Type)
	if err != nil {
		log.Error(err.Error())
		return
	}

	addresses, message := provider.Save(event)

	for _, address := range addresses {
		messageDelivery := m.MessageDelivery{
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

// Stat ...
func (n *notify) Stat() *Stat {
	return n.stat
}

func (n *notify) updateStat() {

	stat := &Stat{
		Workers: len(n.workers),
	}

	if stat.Workers == 0 {
		return
	}

	n.stat = stat
}

// Repeat ...
func (n *notify) Repeat(msg m.MessageDelivery) {
	if !n.isStarted.Load() {
		return
	}

	msg.Status = m.MessageStatusInProgress
	_ = n.adaptor.MessageDelivery.SetStatus(msg)

	n.queue <- msg
}

// AddProvider ...
func (n *notify) AddProvider(name string, provider Provider) {
	n.providerMu.Lock()
	defer n.providerMu.Unlock()

	if _, ok := n.providerList[name]; ok {
		return
	}

	log.Infof("add new notify provider '%s'", name)
	n.providerList[name] = provider
}

// RemoveProvider ...
func (n *notify) RemoveProvider(name string) {
	n.providerMu.Lock()
	defer n.providerMu.Unlock()

	if _, ok := n.providerList[name]; !ok {
		return
	}

	log.Infof("remove notify provider '%s'", name)
	delete(n.providerList, name)
}

// Provider ...
func (n *notify) Provider(name string) (provider Provider, err error) {
	if name == "" {
		err = common.ErrProviderIsEmpty
		return
	}

	n.providerMu.RLock()
	defer n.providerMu.RUnlock()

	var ok bool
	if provider, ok = n.providerList[name]; !ok {
		log.Warnf("provider '%s' not found", name)
		err = common.ErrNotFound
		return
	}
	return
}
