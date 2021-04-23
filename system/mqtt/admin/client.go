package admin

import (
	"context"
	"github.com/e154/smart-home/common/debug"
	"github.com/golang/protobuf/ptypes/empty"
)

type clientService struct {
	a *Admin
}

func (c *clientService) mustEmbedUnimplementedClientServiceServer() {
	return
}

// List lists clients information which the session is valid in the broker (both connected and disconnected).
func (c *clientService) List(_page, _pageSize uint32) (*ListClientResponse, error) {
	page, pageSize := GetPage(_page, _pageSize)
	clients, total, err := c.a.store.GetClients(page, pageSize)
	debug.Println(clients)
	if err != nil {
		return &ListClientResponse{}, err
	}
	return &ListClientResponse{
		Clients:    clients,
		TotalCount: total,
	}, nil
}

// Get returns the client information for given request client id.
func (c *clientService) Get(ctx context.Context, req *GetClientRequest) (*GetClientResponse, error) {
	if req.ClientId == "" {
		return nil, ErrInvalidArgument("client_id", "")
	}
	client := c.a.store.GetClientByID(req.ClientId)
	if client == nil {
		return nil, ErrNotFound
	}
	return &GetClientResponse{
		Client: client,
	}, nil
}

// Delete force disconnect.
func (c *clientService) Delete(ctx context.Context, req *DeleteClientRequest) (*empty.Empty, error) {
	if req.ClientId == "" {
		return nil, ErrInvalidArgument("client_id", "")
	}
	if req.CleanSession {
		c.a.clientService.TerminateSession(req.ClientId)
	} else {
		c := c.a.clientService.GetClient(req.ClientId)
		if c != nil {
			c.Close()
		}
	}
	return &empty.Empty{}, nil
}
