package server

import (
	"net"
	"io"
	"log"
)

func netToSerial(netconn net.Conn, serialconn io.ReadWriteCloser, done chan bool) {
	b, err := io.Copy(serialconn, netconn)
	if err != nil {
		log.Printf("Error copying from network to serial: %v", err)
		return
	}
	done <- true
	log.Printf("netToSerial connection closing.  %v bytes written.", b)
	return
}

func serialToNet(netconn net.Conn, serialconn io.ReadWriteCloser, done chan bool) {
	b, err := io.Copy(netconn, serialconn)
	if err != nil {
		log.Printf("Error copying from serial to network: %v", err)
		return
	}
	done <- true
	log.Printf("serialToNet connection closing.  %v bytes written.", b)
	return
}
