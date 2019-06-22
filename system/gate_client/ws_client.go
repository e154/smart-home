package gate_client

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	writeWait = 10 * time.Second
)

type WsClient struct {
	adaptors    *adaptors.Adaptors
	isConnected bool
	conn        *websocket.Conn
	interrupt   chan struct{}
	enabled     bool
	settings    *Settings
	delta       time.Duration
	cb          IWsCallback
}

func NewWsClient(adaptors *adaptors.Adaptors,
	cb IWsCallback) *WsClient {
	client := &WsClient{
		adaptors:  adaptors,
		interrupt: make(chan struct{}, 1),
		cb:        cb,
	}
	go client.worker()
	return client
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
		if !client.enabled ||
			client.settings == nil ||
			client.settings.Address == "" {
			time.Sleep(time.Second * 5)
			continue
		}
		client.delta += time.Second
		log.Debugf("Wait time %v ...", client.delta)
		time.Sleep(client.delta)
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

	requestHeader := http.Header{
		"X-Client-Type": {ClientTypeServer},
	}
	if client.settings.GateServerToken != "" {
		requestHeader.Add("X-API-Key", client.settings.GateServerToken)
	}

	if client.conn, _, err = websocket.DefaultDialer.Dial(u.String(), requestHeader); err != nil {
		return
	}
	client.isConnected = true

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
				//fmt.Printf("recv: %s\n", string(message))
				client.cb.onMessage(message)
			default:
				log.Errorf("unknown message type(%v)", messageType)
			}
		}
	}()

	go client.cb.onConnected()

	log.Infof("Connect %v successfully", u.String())

	<-client.interrupt
	err = client.write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	go client.cb.onClosed()
}

func (client *WsClient) write(opCode int, payload []byte) (err error) {

	if !client.isConnected {
		err = ErrGateNotConnected
		return
	}

	if err = client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}

	err = client.conn.WriteMessage(opCode, payload)
	return
}
