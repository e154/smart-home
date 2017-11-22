package stream

import (
	"encoding/json"
	"github.com/e154/smart-home/api/log"
	"sync"
)

const (
	maxMessageSize = 512
)

type Hub interface {
	Broadcast(message string)
	Subscribe(command string, f func(client *Client, value interface{}))
	UnSubscribe(command string)
	AddClient(client *Client)
}

var instantiated *hub

type hub struct {
	sessions map[*Client]bool
	subscribers map[string]func(client *Client, value interface{})
	sync.Mutex
	broadcast chan string
}

func (h *hub) AddClient(client *Client) {
	log.Infof("new sockjs session established, from ip: %s", client.Ip)

	h.sessions[client] = true
	var closedSession = make(chan struct{})
	for {
		if msg, err := client.Session.Recv(); err == nil {
			h.Recv(client, msg)
			continue
		}
		break
	}

	delete(h.sessions, client)
	close(closedSession)
	log.Infof("sockjs session from ip: %s closed", client.Ip)
}

func (h *hub) Run() {

	for {
		select {
		case m := <-h.broadcast:
			for client := range h.sessions {
				client.Session.Send(m)
			}
		}
	}
}

func (h *hub) Recv(client *Client, message string) {
	re := map[string]interface{}{}
	if err := json.Unmarshal([]byte(message), &re); err != nil {
		client.Notify("error", err.Error())
		return
	}

	for key, value := range re {

		switch key {
		case "client_info":
			client.UpdateInfo(value)

		default:
			for command, f := range h.subscribers {
				if key == command {
					f(client, value)
				}
			}
		}
	}
}

func (h *hub) Send(client *Client, message string) {
	client.Session.Send(message)
}

func (h *hub) Broadcast(message string) {
	h.Lock()
	h.broadcast <- message
	h.Unlock()
}

func (h *hub) Clients() (clients []*Client) {

	clients = []*Client{}
	for c := range h.sessions {
		clients = append(clients, c)
	}

	return
}

func (h *hub) Subscribe(command string, f func(client *Client, value interface{})) {
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
}

func (h *hub) UnSubscribe(command string) {
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
}

func GetHub() Hub {
	if instantiated == nil {
		instantiated = &hub{
			sessions: make(map[*Client]bool),
			broadcast: make(chan string, maxMessageSize),
			subscribers: make(map[string]func(client *Client, value interface{})),
		}

		go instantiated.Run()
	}
	return instantiated
}