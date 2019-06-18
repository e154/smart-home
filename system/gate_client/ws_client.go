package gate_client

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/gorilla/websocket"
	"net/url"
	"strings"
	"time"
)

const (
	writeWait      = 10 * time.Second
)

type WsClient struct {
	adaptors    *adaptors.Adaptors
	isConnected bool
	conn        *websocket.Conn
	interrupt   chan struct{}
	enabled     bool
	settings    *Settings
	delta       time.Duration
}

func NewWsClient(adaptors *adaptors.Adaptors) *WsClient {
	client := &WsClient{
		adaptors:  adaptors,
		interrupt: make(chan struct{}, 1),
	}
	go client.worker()
	return client
}

func (client *WsClient) GetToken() (token string, err error) {
	return
}

func (client *WsClient) Close() {
	client.enabled = false
	if client.isConnected {
		log.Info("Close")
		client.interrupt <- struct{}{}
	}
}

func (client *WsClient) Connect(settings *Settings) {
	if client.isConnected {
		return
	}

	log.Info("Connect")

	client.enabled = true
	client.settings = settings
}

func (client *WsClient) worker() {
	client.delta = time.Second
	for {
		client.delta += time.Second
		log.Debugf("Wait time %v ...", client.delta)
		time.Sleep(client.delta)
		if !client.enabled {
			continue
		}
		client.connect()
	}
}

func (client *WsClient) connect() {

	if client.isConnected {
		return
	}

	startTime := time.Now()

	var err error
	defer func() {
		if since := time.Since(startTime).Seconds(); since > 10 {
			client.delta = time.Second
			log.Infof("Connect channel closed after %v sec", since)
		}

		client.isConnected = false
		if client.conn != nil {
			_ = client.conn.Close()
		}

		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				return
			}
			log.Error(err.Error())
		}

	}()

	u := url.URL{Scheme: "ws", Host: client.settings.Address, Path: "ws"}
	if client.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil); err != nil {
		return
	}
	client.isConnected = true

	log.Info("Connect successfully")

	var messageType int
	var message []byte

	go func() {
		for {
			if messageType, message, err = client.conn.ReadMessage(); err != nil {
				//log.Error(err.Error())
				client.interrupt <- struct{}{}
				return
			}
			switch messageType {
			case websocket.TextMessage:
				fmt.Printf("recv: %s\n", string(message))
			default:
				log.Errorf("unknown message type(%v)", messageType)
			}
		}
	}()

	<-client.interrupt
	err = client.Write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (client *WsClient) Write(opCode int, payload []byte) error {
	_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return client.conn.WriteMessage(opCode, payload)
}
