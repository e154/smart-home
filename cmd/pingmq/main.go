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

// The following commands will run pingmq as a server, pinging the 8.8.8.0/28 CIDR
// block, and publishing the results to /ping/success/{ip} and /ping/failure/{ip}
// topics every 30 seconds. `sudo` is needed because we are using RAW sockets and
// that requires root privilege.
//
//   $ go build
//   $ sudo ./pingmq server -p 8.8.8.0/28 -i 30
//
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
package main

import (
	"github.com/e154/smart-home/cmd/pingmq/commands"
)

func main() {
	commands.Pingmq.Execute()
}
