package mqtt

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"os"
	"time"
)

type Client struct {
	qos           byte
	topic, uri    string
	client        *service.Client
	authenticator *Authenticator
}

func NewClient(uri, topic string,
	authenticator *Authenticator) (client *Client, err error) {

	var qos byte = 0x0

	// Instantiates a new Client
	c := &service.Client{}

	client = &Client{
		qos:           qos,
		topic:         topic,
		client:        c,
		uri:           uri,
		authenticator: authenticator,
	}

	return
}

func (c *Client) Connect() (err error) {

	log.Debug("Connect")

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetWillQos(c.qos)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(fmt.Sprintf("mqclient%d%d", os.Getpid(), time.Now().Unix())))
	msg.SetKeepAlive(300)
	msg.SetUsername([]byte("local"))
	msg.SetPassword([]byte(c.authenticator.LocalClientUuid()))

	// Connects to the remote server at 127.0.0.1 port 1883
	if err = c.client.Connect(c.uri, msg); err != nil {
		return
	}

	return
}

func (c *Client) Subscribe(topic string, onComplete service.OnCompleteFunc, onPublish service.OnPublishFunc) (err error) {

	log.Debugf("node subscribe to %s", topic)
	submsg := message.NewSubscribeMessage()
	if err = submsg.AddTopic([]byte(topic), c.qos); err != nil {
		return
	}

	err = c.client.Subscribe(submsg, onComplete, onPublish)

	return
}

// Disconnects from the server
func (c *Client) Disconnect() {
	if c.client != nil {
		c.client.Disconnect()
	}
}

func (c *Client) Publish(v []byte) (err error) {

	// Creates a new PUBLISH message with the appropriate contents for publishing
	pubmsg := message.NewPublishMessage()
	if err = pubmsg.SetTopic([]byte(c.topic + "/req")); err != nil {
		return
	}
	pubmsg.SetPayload(v)
	if err = pubmsg.SetQoS(c.qos); err != nil {
		return
	}

	// Publishes to the server by sending the message
	c.client.Publish(pubmsg, nil)

	return
}