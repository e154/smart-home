package server

import (
	"net"
	"fmt"
	"log"
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

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return
	}

	log.Printf("Start server on %s:%d\r\n", addr, port)

	go func(){
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
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

func (s *Server) handleConn(conn net.Conn) {

}

func init() {
	ServerPtr()
}