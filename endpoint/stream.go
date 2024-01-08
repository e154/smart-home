// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/e154/smart-home/system/stream"
)

var (
	upgrader = websocket.Upgrader{}
)

type StreamEndpoint struct {
	*CommonEndpoint
	stream *stream.Stream
}

// NewStreamEndpoint ...
func NewStreamEndpoint(common *CommonEndpoint, stream *stream.Stream) *StreamEndpoint {
	return &StreamEndpoint{
		CommonEndpoint: common,
		stream:         stream,
	}
}

func (s *StreamEndpoint) Subscribe(ctx echo.Context, currentUser *m.User) error {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	s.stream.NewConnection(ws, currentUser)

	return nil
}
