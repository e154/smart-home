package graceful_service

import (
	"sync"
)

type IGracefulClient interface {
	Shutdown()
}

type GracefulServicePool struct {
	cfg     *GracefulServiceConfig
	m       sync.Mutex
	clients map[int]IGracefulClient
}

func NewGracefulServicePool(cfg *GracefulServiceConfig) *GracefulServicePool {
	return &GracefulServicePool{
		cfg:     cfg,
		clients: make(map[int]IGracefulClient),
	}
}

func (h *GracefulServicePool) subscribe(client IGracefulClient) (id int) {
	h.m.Lock()
	id = len(h.clients)
	h.clients[id] = client
	h.m.Unlock()
	return
}

func (h *GracefulServicePool) unsubscribe(id int) {
	h.m.Lock()
	if _, ok := h.clients[id]; ok {
		delete(h.clients, id)
	}
	h.m.Unlock()
}

func (h *GracefulServicePool) shutdown() {
	h.m.Lock()
	i := len(h.clients)
	for ;i>=0;i-- {
		client := h.clients[i]
		if client != nil {
			client.Shutdown()
		}
	}
	h.m.Unlock()
}
