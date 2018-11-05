package graceful_service

import (
	"sync"
	"github.com/e154/smart-home/system/uuid"
)

type IGracefulClient interface {
	Shutdown()
}

type GracefulServicePool struct {
	cfg     *GracefulServiceConfig
	m       sync.Mutex
	clients map[string]IGracefulClient
}

func NewGracefulServicePool(cfg *GracefulServiceConfig) *GracefulServicePool {
	return &GracefulServicePool{
		cfg:     cfg,
		clients: make(map[string]IGracefulClient),
	}
}

func (h *GracefulServicePool) subscribe(client IGracefulClient) (id string) {
	h.m.Lock()
	for {
		id = uuid.NewV4().String()
		if _, ok := h.clients[id]; !ok {
			break
		}
	}

	h.clients[id] = client
	h.m.Unlock()

	return
}

func (h *GracefulServicePool) unsubscribe(id string) {
	h.m.Lock()
	if _, ok := h.clients[id]; ok {
		delete(h.clients, id)
	}
	h.m.Unlock()
}

func (h *GracefulServicePool) shutdown() {
	h.m.Lock()
	for _, client := range h.clients {
		client.Shutdown()
	}
	h.m.Unlock()
}
