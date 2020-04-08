// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package stream

import (
	"errors"
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	log        = common.MustGetLogger("stream")
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// StreamService ...
type StreamService struct {
	Hub *Hub
}

// NewStreamService ...
func NewStreamService(hub *Hub) *StreamService {
	return &StreamService{
		Hub: hub,
	}
}

// Broadcast ...
func (s *StreamService) Broadcast(message []byte) {
	s.Hub.Broadcast(message)
}

// Subscribe ...
func (s *StreamService) Subscribe(command string, f func(client IStreamClient, msg Message)) {
	s.Hub.Subscribe(command, f)
}

// UnSubscribe ...
func (s *StreamService) UnSubscribe(command string) {
	s.Hub.UnSubscribe(command)
}

// Ws ...
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

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	client := &Client{
		Connect:   conn,
		Ip:        ctx.ClientIP(),
		Referer:   ctx.Request.Referer(),
		UserAgent: ctx.Request.UserAgent(),
		Send:      make(chan []byte),
	}

	go client.WritePump()
	w.Hub.AddClient(client)
}
