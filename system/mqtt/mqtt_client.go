package mqtt

import (
	"github.com/surgemq/surgemq/service"
	"github.com/surgemq/message"
	"fmt"
	"os"
	"time"
)

type Client struct {
	qos        byte
	topic, uri string
	client     *service.Client
	onComplete service.OnCompleteFunc
	onPublish  service.OnPublishFunc
}

func NewClient(uri, topic string,
	onComplete service.OnCompleteFunc,
	onPublish service.OnPublishFunc) (client *Client, err error) {

	var qos byte = 0x0

	// Instantiates a new Client
	c := &service.Client{}

	client = &Client{
		qos:        qos,
		topic:      topic,
		client:     c,
		uri:        uri,
		onComplete: onComplete,
		onPublish:  onPublish,
	}

	return
}

func (c *Client) Connect() (err error) {

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetWillQos(c.qos)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(fmt.Sprintf("mqclient%d%d", os.Getpid(), time.Now().Unix())))
	msg.SetKeepAlive(300)
	//msg.SetWillTopic([]byte("will"))
	//msg.SetWillMessage([]byte("send me home"))
	//msg.SetUsername([]byte("surgemq"))
	//msg.SetPassword([]byte("verysecret"))

	// Connects to the remote server at 127.0.0.1 port 1883
	if err = c.client.Connect(c.uri, msg); err != nil {
		return
	}

	// Creates a new SUBSCRIBE message to subscribe to topic "topic"
	submsg := message.NewSubscribeMessage()
	if err = submsg.AddTopic([]byte(c.topic), c.qos); err != nil {
		return
	}

	// Subscribes to the topic by sending the message. The first nil in the function
	// call is a OnCompleteFunc that should handle the SUBACK message from the server.
	// Nil means we are ignoring the SUBACK messages. The second nil should be a
	// OnPublishFunc that handles any messages send to the client because of this
	// subscription. Nil means we are ignoring any PUBLISH messages for this topic.
	if err = c.client.Subscribe(submsg, c.onComplete, c.onPublish); err != nil {
		return
	}

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
	if err = pubmsg.SetTopic([]byte(c.topic)); err != nil {
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
