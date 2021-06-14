// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package commands

import (
	"fmt"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"strings"
	"time"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "server starts a SurgeMQ server and publishes to it all the ping results",
	}

	serverURI    string
	serverQuiet  bool
	serverIPs    strlist
	pingInterval int

	s mqtt.GMqttServer
)

type strlist []string

// String ...
func (this *strlist) String() string {
	return fmt.Sprint(*this)
}

// Type ...
func (this *strlist) Type() string {
	return "strlist"
}

// Set ...
func (this *strlist) Set(value string) error {
	for _, ip := range strings.Split(value, ",") {
		*this = append(*this, ip)
	}

	return nil
}

func init() {
	serverCmd.Flags().StringVarP(&serverURI, "uri", "u", "0.0.0.0:1883", "URI to run the server on")
	serverCmd.Flags().BoolVarP(&serverQuiet, "quiet", "q", false, "print out ping results")
	serverCmd.Flags().VarP(&serverIPs, "ping", "p", "Comma separated list of IPv4 addresses to ping")
	serverCmd.Flags().IntVarP(&pingInterval, "interval", "i", 60, "ping interval in seconds")
	serverCmd.Run = serv

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
