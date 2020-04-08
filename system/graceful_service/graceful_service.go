// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package graceful_service

import (
	"github.com/e154/smart-home/common"
	"os"
	"os/signal"
	"syscall"
)

var (
	log = common.MustGetLogger("graceful_service")
)

// GracefulService ...
type GracefulService struct {
	cfg  *GracefulServiceConfig
	pool *GracefulServicePool
	done chan struct{}
}

// NewGracefulService ...
func NewGracefulService(cfg *GracefulServiceConfig,
	hub *GracefulServicePool) (graceful *GracefulService) {
	graceful = &GracefulService{
		cfg:  cfg,
		pool: hub,
		done: make(chan struct{}, 1),
	}

	log.Info("Graceful shutdown service started")

	return
}

// Wait ...
func (p GracefulService) Wait() {

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-gracefulStop
		p.pool.shutdown()
		p.done <- struct{}{}

	}()

	for {
		select {
		case <-p.done:
			log.Info("Shutdown")
			os.Exit(0)
		}
	}

	close(p.done)
	close(gracefulStop)
}

// Subscribe ...
func (p GracefulService) Subscribe(client IGracefulClient) (id int) {
	id = p.pool.subscribe(client)
	return
}

// Unsubscribe ...
func (p GracefulService) Unsubscribe(id int) {
	p.pool.unsubscribe(id)
	return
}

// Shutdown ...
func (p GracefulService) Shutdown() {
	p.pool.shutdown()
	return
}
