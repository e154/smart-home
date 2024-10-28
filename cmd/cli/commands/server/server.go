// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

// The following commands will run pingmq as a server, pinging the 8.8.8.0/28 CIDR
// block, and publishing the results to /ping/success/{ip} and /ping/failure/{ip}
// topics every 30 seconds. `sudo` is needed because we are using RAW sockets and
// that requires root privilege.
//
//	$ go build
//	$ sudo ./pingmq server -p 8.8.8.0/28 -i 30
package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/mqtt"

	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/koron/netx"
	"github.com/spf13/cobra"
)

var (
	// Server ...
	Server = &cobra.Command{
		Use:   "server",
		Short: "server starts a SurgeMQ server and publishes to it all the ping results",
	}

	serverURI    string
	serverQuiet  bool
	serverIPs    strlist
	pingInterval int

	s mqtt.GMqttServer

	p *netx.Pinger
)

type strlist []string

// String ...
func (s *strlist) String() string {
	return fmt.Sprint(*s)
}

// Type ...
func (s *strlist) Type() string {
	return "strlist"
}

// Set ...
func (s *strlist) Set(value string) error {
	for _, ip := range strings.Split(value, ",") {
		*s = append(*s, ip)
	}

	return nil
}

func init() {
	Server.Flags().StringVarP(&serverURI, "uri", "u", "0.0.0.0:1883", "URI to run the server on")
	Server.Flags().BoolVarP(&serverQuiet, "quiet", "q", false, "print out ping results")
	Server.Flags().VarP(&serverIPs, "ping", "p", "Comma separated list of IPv4 addresses to ping")
	Server.Flags().IntVarP(&pingInterval, "interval", "i", 60, "ping interval in seconds")
	Server.Run = serv

}

func serv(cmd *cobra.Command, args []string) {

	log.Printf("Starting server...")
	go func() {
		ln, err := net.Listen("tcp", serverURI)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		options := []server.Options{
			server.WithTCPListener(ln),
		}

		// Create a new server
		s = server.New(options...)

		if err = s.Run(); err != nil {
			log.Println(err.Error())
		}

	}()
	time.Sleep(300 * time.Millisecond)

	log.Printf("Starting pinger...")

	pinger()
}

func pinger() {

	p = &netx.Pinger{}
	if err := p.AddIPs(serverIPs); err != nil {
		log.Fatal(err)
	}

	cnt := 0
	tick := time.NewTicker(time.Duration(pingInterval) * time.Second)

	for {
		if cnt != 0 {
			<-tick.C
		}

		res, err := p.Start()
		if err != nil {
			log.Fatal(err)
		}

		for pr := range res {
			if !serverQuiet {
				log.Println(pr)
			}

			var topic string

			// Creates a new PUBLISH message with the appropriate contents for publishing
			if pr.Err != nil {
				topic = fmt.Sprintf("/ping/failure/%s", pr.Src)
			} else {
				topic = fmt.Sprintf("/ping/success/%s", pr.Src)
			}

			payload, err := pr.GobEncode()
			if err != nil {
				log.Printf("pinger: Error from GobEncode: %v\n", err)
				continue
			}

			// Publishes to the server
			s.Publisher().Publish(&gmqtt.Message{
				Topic:    topic,
				Payload:  payload,
				QoS:      0,
				Retained: true,
			})
		}

		p.Stop()
		cnt++
	}
}
