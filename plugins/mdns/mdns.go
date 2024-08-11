// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package mdns

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/atomic"

	"github.com/grandcat/zeroconf"
)

type Dns struct {
	server *zeroconf.Server
	isScan *atomic.Bool
}

func NewDns() *Dns {
	return &Dns{
		isScan: atomic.NewBool(false),
	}
}

func (d *Dns) Scan(service, domain string) {
	if !d.isScan.CompareAndSwap(false, true) {
		return
	}

	defer func() {
		d.isScan.Store(false)
	}()

	log.Info("Starting scanning ...")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			j, _ := json.MarshalIndent(entry, " ", " ")
			log.Info(string(j))
		}
		log.Info("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		log.Error(err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}

func (d *Dns) Register(instance, service, domain, ipAddr, host string, port int64, text []string) {
	var err error
	if ipAddr != "" {
		log.Infof("Registering proxy instance(%s) host(%s) service(%s) domain(%s) port(%d) text(%v) on network", instance, host, service, domain, port, text)
		d.server, err = zeroconf.RegisterProxy(instance, service, domain, int(port), host, []string{ipAddr}, text, nil)
	} else {
		log.Infof("Published service instance(%s) service(%s) domain(%s) port(%d) text(%v) on network", instance, service, domain, port, text)
		d.server, err = zeroconf.Register(instance, service, domain, int(port), text, nil)
	}
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func (d *Dns) Shutdown() {
	if d.server != nil {
		d.server.Shutdown()
	}
	d.server = nil
}
