package server

import (
	"net"
	"fmt"
	"log"
	"net/rpc"
)

// Singleton
var instantiated *Server = nil

func ServerPtr() *Server {
	if instantiated == nil {
		instantiated = &Server{
			clients: make(map[*Client]bool),
		}
	}
	return instantiated
}

type Server struct {
	listener net.Listener
	clients map[*Client]bool
}

func (s *Server) Start(addr string, port int) (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return
	}

	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	log.Printf("Start server on %s:%d\r\n", addr, port)

	go func() {
		for {
			conn, _ := s.listener.Accept()
			go s.AddClient(conn)
		}
	}()

	return
}

func (s *Server) AddClient(conn net.Conn) {

	addr := conn.RemoteAddr().String()
	log.Printf("New client %s\r\n", addr)

	client := &Client{}
	s.clients[client] = true
	defer func() {
		conn.Close()
		delete(s.clients, client)
		log.Printf("Connection %s closed", addr)
	}()

	client.listener(conn)
}

func init() {
	ServerPtr()
	rpc.Register(&Modbus{})
	rpc.Register(&Node{})
}