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
	"github.com/DrmagicE/gmqtt"
	"github.com/koron/netx"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	Pingmq = &cobra.Command{
		Use:   "pingmq",
		Short: "Pingmq is a program designed to demonstrate the SurgeMQ usage.",
		Long: `Pingmq demonstrates the use of SurgeMQ by pinging a list of hosts, 
publishing the result to any clients subscribed to two topics:
/ping/success/{ip} and /ping/failure/{ip}.`,
	}

	p *netx.Pinger
)

func init() {

	Pingmq.AddCommand(serverCmd)
	Pingmq.AddCommand(clientCmd)

	done = make(chan struct{})
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
