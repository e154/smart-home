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

package admin

import (
	"fmt"
	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/server"
)

var _ server.Plugin = (*Admin)(nil)

const Name = "admin"

func init() {
	server.RegisterPlugin(Name, New)
}

func New(config config.Config) (server.Plugin, error) {
	return &Admin{}, nil
}

type Admin struct {
	statsReader   server.StatsReader
	publisher     server.Publisher
	clientService server.ClientService
	store         *store
}


func (a *Admin) Load(service server.Server) error {
	fmt.Println("------ load")
	//apiRegistrar := service.APIRegistrar()

	//apiRegistrar := service.APIRegistrar()
	//RegisterClientServiceServer(apiRegistrar, &clientService{a: a})
	//RegisterSubscriptionServiceServer(apiRegistrar, &subscriptionService{a: a})
	//RegisterPublishServiceServer(apiRegistrar, &publisher{a: a})
	//err := a.registerHTTP(apiRegistrar)
	//if err != nil {
	//	return err
	//}
	a.statsReader = service.StatsManager()
	a.store = newStore(a.statsReader)
	a.store.subscriptionService = service.SubscriptionService()
	a.publisher = service.Publisher()
	a.clientService = service.ClientService()
	return nil
}

func (a *Admin) Unload() error {
	fmt.Println("------ unload")
	return nil
}

func (a *Admin) Name() string {
	return Name
}

// GetClient ...
func (a *Admin) GetClients(_page, _pageSize uint32) (clients, total interface{}, err error) {
	page, pageSize := GetPage(_page, _pageSize)
	clients, total, err = a.store.GetClients(page, pageSize)
	return
}


// GetClient ...
func (a *Admin) GetClient(_page, _pageSize uint32) (clients, total interface{}, err error) {
	page, pageSize := GetPage(_page, _pageSize)
	clients, total, err = a.store.GetClients(page, pageSize)
	return
}
