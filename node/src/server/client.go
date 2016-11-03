package server

import (
	"net/rpc"
	"io"
)

type Client struct {
	connection *rpc.Client
}

func (c *Client) listener(conn io.ReadWriteCloser) {
	rpc.ServeConn(conn)
}
