package onvif

type ClientBind struct {
	client *Client
}

func NewClientBind(client *Client) *ClientBind {
	return &ClientBind{client: client}
}

func (c *ClientBind) ContinuousMove(X, Y float32) {
	c.client.ContinuousMove(X, Y)
}

func (c *ClientBind) StopContinuousMove() {
	c.client.StopContinuousMove()
}
