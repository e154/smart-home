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
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/koron/netx"
	"github.com/spf13/cobra"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

type strlist []string

func (this *strlist) String() string {
	return fmt.Sprint(*this)
}

func (this *strlist) Type() string {
	return "strlist"
}

func (this *strlist) Set(value string) error {
	for _, ip := range strings.Split(value, ",") {
		*this = append(*this, ip)
	}

	return nil
}

var (
	pingmqCmd = &cobra.Command{
		Use:   "pingmq",
		Short: "Pingmq is a program designed to demonstrate the SurgeMQ usage.",
		Long: `Pingmq demonstrates the use of SurgeMQ by pinging a list of hosts, 
publishing the result to any clients subscribed to two topics:
/ping/success/{ip} and /ping/failure/{ip}.`,
	}

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "server starts a SurgeMQ server and publishes to it all the ping results",
	}

	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "client subscribes to the pingmq server and prints out the ping results",
	}

	serverURI      string
	serverQuiet    bool
	serverIPs      strlist
	user, password string

	pingInterval int

	clientURI    string
	clientTopics strlist

	s *gmqtt.Server
	c *service.Client
	p *netx.Pinger

	wg sync.WaitGroup

	done chan struct{}
)

func init() {

	serverCmd.Flags().StringVarP(&serverURI, "uri", "u", "tcp://:1883", "URI to run the server on")
	serverCmd.Flags().BoolVarP(&serverQuiet, "quiet", "q", false, "print out ping results")
	serverCmd.Flags().VarP(&serverIPs, "ping", "p", "Comma separated list of IPv4 addresses to ping")
	serverCmd.Flags().IntVarP(&pingInterval, "interval", "i", 60, "ping interval in seconds")
	serverCmd.Run = server

	clientCmd.Flags().StringVarP(&clientURI, "server", "s", "tcp://127.0.0.1:1883", "PingMQ server to connect to")
	clientCmd.Flags().VarP(&clientTopics, "topic", "t", "Comma separated list of topics to subscribe to")
	clientCmd.Flags().StringVarP(&user, "user", "u", "node1", "user name")
	clientCmd.Flags().StringVarP(&password, "password", "p", "node1", "password")
	clientCmd.Run = client

	pingmqCmd.AddCommand(serverCmd)
	pingmqCmd.AddCommand(clientCmd)

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
			s.Publish(topic, payload, 0, false)
		}

		p.Stop()
		cnt++
	}
}

func server(cmd *cobra.Command, args []string) {

	// Create a new server
	s = gmqtt.NewServer(gmqtt.DefaultConfig)

	log.Printf("Starting server...")
	go func() {
		if _, err := net.Listen("tcp", serverURI); err != nil {
			log.Fatal(err.Error())
		}
	}()
	time.Sleep(300 * time.Millisecond)

	log.Printf("Starting pinger...")
	pinger()
}

func client(cmd *cobra.Command, args []string) {
	// Instantiates a new Client
	c = &service.Client{}

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(fmt.Sprintf("pingmqclient%d%d", os.Getpid(), time.Now().Unix())))
	msg.SetKeepAlive(300)
	if user != "" {
		msg.SetUsernameFlag(true)
		msg.SetUsername([]byte(user))
	}
	if password != "" {
		msg.SetPasswordFlag(true)
		msg.SetPassword([]byte(password))
	}

	// Connects to the remote server at 127.0.0.1 port 1883
	if err := c.Connect(clientURI, msg); err != nil {
		log.Fatal(err)
	}

	// Creates a new SUBSCRIBE message to subscribe to topic "abc"
	submsg := message.NewSubscribeMessage()

	for _, t := range clientTopics {
		submsg.AddTopic([]byte(t), 0)
	}

	c.Subscribe(submsg, nil, onPublish)

	<-done
}

func onPublish(msg *message.PublishMessage) error {
	pr := &netx.PingResult{}
	if err := pr.GobDecode(msg.Payload()); err != nil {
		fmt.Println(string(msg.Payload()))
		return nil
	}

	log.Println(pr)
	return nil
}

func main() {
	pingmqCmd.Execute()
}
