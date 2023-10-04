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
	"context"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/stream"
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

func (s *StreamEndpoint) Subscribe(ctx context.Context, server api.StreamService_SubscribeServer) error {

	var user *m.User
	request, err := server.Recv()
	if err != nil {
		return err
	}

	claims, err := s.jwtManager.Verify(request.GetAccessToken())
	if err != nil {
		log.Error(err.Error())
		return apperr.ErrTokenIsDeprecated
	}

	if user, err = s.adaptors.User.GetById(ctx, claims.UserId); err != nil {
		return err
	}

	return s.stream.NewConnection(server, user)
}
