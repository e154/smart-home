package stream

import (
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
)

type ConnectType string

const (
	WEBSOCK = ConnectType("websock")
)

var (
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:  func(r *http.Request) bool {
			return true
		},
	}
)

func (w *StreamService) Ws(ctx *gin.Context) {

	// CORS
	ctx.Writer.Header().Del("Access-Control-Allow-Credentials")

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("Failed to set websocket upgrade: %v", err)
		return
	}
	if _, ok := err.(websocket.HandshakeError); ok {
		ctx.AbortWithError(400, errors.New("not a websocket handshake"))
		return
	} else if err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	client := &Client{
		ConnType:  WEBSOCK,
		Connect:   conn,
		Ip:        ctx.ClientIP(),
		Referer:   ctx.Request.Referer(),
		UserAgent: ctx.Request.UserAgent(),
		Send:      make(chan []byte),
	}

	go client.WritePump()
	w.Hub.AddClient(client)
}
