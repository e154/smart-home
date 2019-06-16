package gate_client

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/gorilla/websocket"
	"net/url"
)

type WsClient struct {
	adaptors    *adaptors.Adaptors
	isConnected bool
}

func NewWsClient(adaptors *adaptors.Adaptors) *WsClient {
	return &WsClient{
		adaptors: adaptors,
	}
}

func (client *WsClient) GetToken() (token string, err error) {
	return
}

func (client *WsClient) Connect(settings *Settings) {

	if client.isConnected {
		return
	}

	go func() {
		if err := client.connect(settings); err != nil {
			log.Error(err.Error())
		}
	}()
}

func (client *WsClient) connect(settings *Settings) (err error) {

	u := url.URL{Scheme: "ws", Host: settings.Address}

	var c *websocket.Conn
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}

	done := make(chan struct{})

	defer func() {
		close(done)
		client.isConnected = false
		_ = c.Close()
	}()

	client.isConnected = true

	var message []byte
	go func() {
		defer close(done)
		for {
			if _, message, err = c.ReadMessage(); err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s", message)
		}
	}()

	return
}
