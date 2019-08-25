package gate_client

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	writeWait = 10 * time.Second
)

type WsClient struct {
	sync.Mutex
	adaptors        *adaptors.Adaptors
	isConnected     bool
	conn            *websocket.Conn
	interrupt       chan struct{}
	enabled         bool
	settings        *Settings
	delta           time.Duration
	cb              IWsCallback
	status          string
	inProgress      bool
	closeProgress   bool
	connectProgress bool
}

func NewWsClient(adaptors *adaptors.Adaptors,
	cb IWsCallback) *WsClient {
	client := &WsClient{
		adaptors:  adaptors,
		interrupt: make(chan struct{}),
		cb:        cb,
	}
	go client.worker()
	return client
}

func (client *WsClient) Close() {
	if !client.isConnected || client.closeProgress {
		return
	}
	client.closeProgress = true
	client.enabled = false
	if client.isConnected {
		log.Info("Close")
		client.status = GateStatusDisabled
		client.interrupt <- struct{}{}
	}
}

func (client *WsClient) Connect(settings *Settings) {
	if client.isConnected || client.connectProgress {
		return
	}

	log.Info("Connect")

	client.enabled = true
	client.connectProgress = true
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
		client.status = GateStatusWait
		client.delta += time.Second
		//log.Debugf("Wait time %v ...", client.delta)
		time.Sleep(client.delta)
		client.connect()
	}
}

func (client *WsClient) connect() {

	if client.isConnected || client.inProgress {
		return
	}

	startTime := time.Now()
	client.inProgress = true

	var loseChan chan struct{}
	loseChan = make(chan struct{})

	var err error
	defer func() {
		if since := time.Since(startTime).Seconds(); since > 10 {
			client.delta = time.Second
			log.Infof("Connect channel closed after %v sec", since)
		}

		client.closeProgress = false
		client.inProgress = false
		client.isConnected = false
		if client.conn != nil {
			_ = client.conn.Close()
		}

		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				return
			}
			if strings.Contains(err.Error(), "bad handshake") {
				return
			}
			log.Debug(err.Error())
		}

	}()

	u := url.URL{Scheme: "ws", Host: client.settings.Address, Path: "ws"}

	requestHeader := http.Header{
		"X-Client-Type": {ClientTypeServer},
	}
	if client.settings.GateServerToken != "" {
		requestHeader.Add("X-API-Key", client.settings.GateServerToken)
		//log.Debugf("X-API-Key: %v", client.settings.GateServerToken)
	}

	if client.conn, _, err = websocket.DefaultDialer.Dial(u.String(), requestHeader); err != nil {
		return
	}
	client.isConnected = true
	client.connectProgress = false
	client.status = GateStatusConnected

	var messageType int
	var message []byte

	go func() {
		for {
			if messageType, message, err = client.conn.ReadMessage(); err != nil {
				//log.Error(err.Error())
				loseChan <- struct{}{}
				break
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

	select {
	case <-client.interrupt:
	case <-loseChan:
		close(loseChan)
	}

	err = client.write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	go client.cb.onClosed()
}

func (client *WsClient) write(opCode int, payload []byte) (err error) {

	client.Lock()
	if !client.isConnected {
		client.Unlock()
		err = ErrGateNotConnected
		return
	}

	if err = client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		client.Unlock()
		return
	}

	err = client.conn.WriteMessage(opCode, payload)
	client.Unlock()

	return
}
