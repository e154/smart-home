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

// The following command will run pingmq as a client, subscribing to /ping/failure/+
// topic and receiving any failed ping attempts.
//
//   $ ./pingmq client -t /ping/failure/+
//   8.8.8.6: Request timed out for seq 1
//
// The following command will run pingmq as a client, subscribing to /ping/failure/+
// topic and receiving any failed ping attempts.
//
//   $ ./pingmq client -t /ping/success/+
//   8 bytes from 8.8.8.8: seq=1 ttl=56 tos=32 time=21.753711ms
//
// One can also subscribe to a specific IP by using the following command.
//
//   $ ./pingmq client -t /ping/+/8.8.8.8
//   8 bytes from 8.8.8.8: seq=1 ttl=56 tos=32 time=21.753711ms
//
package client

import (
	"fmt"
	"github.com/e154/smart-home/system/mqtt_client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/koron/netx"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var (
	// Client ...
	Client = &cobra.Command{
		Use:   "client",
		Short: "client subscribes to the pingmq server and prints out the ping results",
	}

	clientURI      string
	clientTopics   strlist
	user, password string

	done chan struct{}
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
	Client.Flags().StringVarP(&clientURI, "server", "s", "tcp://127.0.0.1:1883", "PingMQ server to connect to")
	Client.Flags().VarP(&clientTopics, "topic", "t", "Comma separated list of topics to subscribe to")
	Client.Flags().StringVarP(&user, "user", "u", "node1", "user name")
	Client.Flags().StringVarP(&password, "password", "p", "node1", "password")
	Client.Run = client

	done = make(chan struct{})
}

func onPublish(i mqtt.Client, msg mqtt.Message) {

	pr := &netx.PingResult{}
	if err := pr.GobDecode(msg.Payload()); err != nil {
		fmt.Println(string(msg.Payload()))
		return
	}

	log.Println(pr)
}

func client(cmd *cobra.Command, args []string) {

	cfg := &mqtt_client.Config{
		CleanSession: true,
		Broker:       clientURI,
		KeepAlive:    300,
		Username:     user,
		Password:     password,
		ClientID:     fmt.Sprintf("pingmqclient%d%d", os.Getpid(), time.Now().Unix()),
	}
	c, err := mqtt_client.NewClient(cfg)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if err = c.Connect(); err != nil {
		log.Fatalln(err.Error())
		return
	}

	for _, t := range clientTopics {
		if err := c.Subscribe(t, 0, onPublish); err != nil {
			log.Fatalln(err.Error())
		}
	}

	<-done
}
