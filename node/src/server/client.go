package server

import (
	"net"
	"bufio"
	"io"
	"time"
	"log"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Send(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) listener(conn net.Conn) {

	for ;; {
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Print(err.Error())
			return
		}

		log.Print(string(line))
	}
}

func (c *Client) echo(conn net.Conn) {

	for ;; {
		_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}