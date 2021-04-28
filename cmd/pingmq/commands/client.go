package commands

import (
	"fmt"
	"github.com/e154/smart-home/system/mqtt_client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/koron/netx"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "client subscribes to the pingmq server and prints out the ping results",
	}

	clientURI    string
	clientTopics strlist
	user, password string

	done chan struct{}
)

func init() {
	clientCmd.Flags().StringVarP(&clientURI, "server", "s", "tcp://127.0.0.1:1883", "PingMQ server to connect to")
	clientCmd.Flags().VarP(&clientTopics, "topic", "t", "Comma separated list of topics to subscribe to")
	clientCmd.Flags().StringVarP(&user, "user", "u", "node1", "user name")
	clientCmd.Flags().StringVarP(&password, "password", "p", "node1", "password")
	clientCmd.Run = client
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
