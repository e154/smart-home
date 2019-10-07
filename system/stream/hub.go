package stream

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

type Hub struct {
	sessions    map[*Client]bool
	subscribers map[string]func(client *Client, msg Message)
	gateClient  BroadcastClient
	sync.Mutex
	broadcast chan []byte
	interrupt chan os.Signal
}

func NewHub() *Hub {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	hub := &Hub{
		sessions:    make(map[*Client]bool),
		broadcast:   make(chan []byte, maxMessageSize),
		subscribers: make(map[string]func(client *Client, msg Message)),
		interrupt:   interrupt,
	}
	go hub.Run()

	return hub
}

func (h *Hub) AddClient(client *Client) {

	defer func() {
		delete(h.sessions, client)
		log.Infof("websocket session from ip: %s closed", client.Ip)
	}()

	h.Lock()
	h.sessions[client] = true
	h.Unlock()

	switch client.ConnType {
	//case SOCKJS:
	//	log.Infof("new sockjs session established, from ip: %s", client.Ip)
	//
	//	for {
	//		if msg, err := client.Session.Recv(); err == nil {
	//			h.Recv(client, []byte(msg))
	//			continue
	//		}
	//		break
	//	}
	case WEBSOCK:
		log.Infof("new websocket xsession established, from ip: %s", client.Ip)

		//client.Connect.SetReadLimit(maxMessageSize)
		client.Connect.SetReadDeadline(time.Now().Add(pongWait))
		client.Connect.SetPongHandler(func(string) error {
			client.Connect.SetReadDeadline(time.Now().Add(pongWait));
			return nil
		})
		for {
			op, r, err := client.Connect.NextReader()
			if err != nil {
				break
			}
			switch op {
			case websocket.TextMessage:
				message, err := ioutil.ReadAll(r)
				if err != nil {
					break
				}
				h.Recv(client, message)
			}
		}
	default:
		log.Warningf("unknown conn type %s", client.ConnType)
	}
}

func (h *Hub) Run() {

	for {
		select {
		case m := <-h.broadcast:
			h.Lock()
			for client := range h.sessions {
				client.Send <- m
			}
			h.Unlock()
		case <-h.interrupt:
			//fmt.Println("Close websocket client session")
			h.Lock()
			for client := range h.sessions {
				client.Close()
				delete(h.sessions, client)
			}
			h.Unlock()
		}

	}
}

func (h *Hub) Recv(client *Client, b []byte) {

	//fmt.Printf("client(%v), message(%v)\n", client, string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch msg.Command {
	case "client_info":
		client.UpdateInfo(msg)

	default:
		for command, f := range h.subscribers {

			if msg.Command == command {
				f(client, msg)
			}
		}
	}
}

func (h *Hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *Hub) Broadcast(message []byte) {
	h.Lock()
	h.broadcast <- message
	h.Unlock()

	go func() {
		if h.gateClient != nil {
			h.gateClient.Broadcast(message)
		}
	}()
}

func (h *Hub) Clients() (clients []*Client) {

	clients = []*Client{}
	for c := range h.sessions {
		clients = append(clients, c)
	}

	return
}

func (h *Hub) Subscribe(command string, f func(client *Client, msg Message)) {
	log.Infof("subscribe %s", command)
	h.Lock()
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
	h.Unlock()
}

func (h *Hub) UnSubscribe(command string) {
	h.Lock()
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.Unlock()
}

func (h *Hub) Subscriber(command string) (f func(client *Client, msg Message)) {
	h.Lock()
	if h.subscribers[command] != nil {
		f = h.subscribers[command]
	}
	h.Unlock()
	return
}