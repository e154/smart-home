// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

	"github.com/e154/smart-home/common/apperr"
)

var ProviderManager = NewProviderManager()

type manager struct {
	providerMu   *sync.RWMutex
	providerList map[string]Provider
}

func NewProviderManager() *manager {
	return &manager{
		providerMu:   &sync.RWMutex{},
		providerList: make(map[string]Provider),
	}
}

// AddProvider ...
func (n *manager) AddProvider(name string, provider Provider) {
	n.providerMu.Lock()
	defer n.providerMu.Unlock()

	if _, ok := n.providerList[name]; ok {
		return
	}

	log.Infof("add new notify provider '%s'", name)
	n.providerList[name] = provider
}

// RemoveProvider ...
func (n *manager) RemoveProvider(name string) {
	n.providerMu.Lock()
	defer n.providerMu.Unlock()

	if _, ok := n.providerList[name]; !ok {
		return
	}

	log.Infof("remove notify provider '%s'", name)
	delete(n.providerList, name)
}

// Provider ...
func (n *manager) Provider(name string) (provider Provider, err error) {
	if name == "" {
		err = apperr.ErrProviderIsEmpty
		return
	}

	n.providerMu.RLock()
	defer n.providerMu.RUnlock()

	var ok bool
	if provider, ok = n.providerList[name]; !ok {
		log.Warnf("provider '%s' not found", name)
		err = apperr.ErrNotFound
		return
	}
	return
}
