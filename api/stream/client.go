package stream

import (
	"encoding/json"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Client struct {
	Session  sockjs.Session
	//User *rbac.User
	Ip string
	Referer string
	UserAgent string
	Width int
	Height int
	Cookie bool
	Language string
	Platform string
	Location string
	Href string
}

func (c *Client) UpdateInfo(info interface{}) {
	v, ok := info.(map[string]interface{})
	if !ok {
		return
	}

	width, ok := v["width"].(float64)
	if ok {
		c.Width = int(width)
	}

	if height, ok := v["height"].(float64); ok {
		c.Height = int(height)
	}

	if cookie, ok := v["cookie"].(bool); ok {
		c.Cookie = cookie
	}

	if language, ok := v["language"].(string); ok {
		c.Language = language
	}

	if platform, ok := v["platform"].(string); ok {
		c.Platform = platform
	}

	if location, ok := v["location"].(string); ok {
		c.Location = location
	}

	if href, ok := v["href"].(string); ok {
		c.Href = href
	}
}

func (c *Client) Notify(t, b string) {
	msg, _ := json.Marshal(&map[string]interface{}{"type": "notify", "value": &map[string]interface{}{"type": t, "body": b}})
	c.Session.Send(string(msg))
}

func (c *Client) Send(msg string) {
	c.Session.Send(msg)
}