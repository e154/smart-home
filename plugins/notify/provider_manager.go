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
