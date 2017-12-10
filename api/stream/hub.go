package stream

import (
	"encoding/json"
	"github.com/e154/smart-home/api/log"
	"sync"
	"time"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"os/signal"
)

const (
	writeWait = 10 * time.Second
	maxMessageSize = 512
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Hub interface {
	Broadcast(message []byte)
	Subscribe(command string, f func(client *Client, value interface{}))
	UnSubscribe(command string)
	AddClient(client *Client)
}

type hub struct {
	sessions    map[*Client]bool
	subscribers map[string]func(client *Client, value interface{})
	sync.Mutex
	broadcast   chan []byte
	interrupt	chan os.Signal
}

var (
	instantiated *hub
)

func (h *hub) AddClient(client *Client) {

	defer func(){
		delete(h.sessions, client)
		log.Infof("websocket session from ip: %s closed", client.Ip)
	}()

	h.sessions[client] = true
	switch client.ConnType {
	case SOCKJS:
		log.Infof("new sockjs session established, from ip: %s", client.Ip)

		for {
			if msg, err := client.Session.Recv(); err == nil {
				h.Recv(client, []byte(msg))
				continue
			}
			break
		}
	case WEBSOCK:
		log.Infof("new websocket xsession established, from ip: %s", client.Ip)

		//client.Connect.SetReadLimit(maxMessageSize)
		client.Connect.SetReadDeadline(time.Now().Add(pongWait))
		client.Connect.SetPongHandler(func(string) error {
			client.Connect.SetReadDeadline(time.Now().Add(pongWait)); return nil
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

	}
}

func (h *hub) Run() {



	for {
		select {
		case m := <-h.broadcast:
			for client := range h.sessions {
				client.Send <- m
			}
		case <- h.interrupt:
			//fmt.Println("Close websocket client session")
			for client := range h.sessions {
				client.Close()
				delete(h.sessions, client)
			}
		}

	}
}

func (h *hub) Recv(client *Client, message []byte) {

	re := map[string]interface{}{}
	if err := json.Unmarshal(message, &re); err != nil {
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

func (h *hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *hub) Broadcast(message []byte) {
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if instantiated == nil {
		instantiated = &hub{
			sessions: make(map[*Client]bool),
			broadcast: make(chan []byte, maxMessageSize),
			subscribers: make(map[string]func(client *Client, value interface{})),
			interrupt: interrupt,
		}

		go instantiated.Run()
	}
	return instantiated
}