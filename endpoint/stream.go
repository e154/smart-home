package endpoint

import (
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

func (s *StreamEndpoint) Subscribe(server api.StreamService_SubscribeServer) error {

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

	if user, err = s.adaptors.User.GetById(claims.UserId); err != nil {
		return err
	}

	return s.stream.NewConnection(server, user)
}
