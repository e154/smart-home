package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

type Client struct {
	engine    *gin.Engine
	token     string
	basicAuth string
}

func NewClient(engine *gin.Engine) *Client {
	return &Client{
		engine: engine,
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) BasicAuth(login, pass string) {
	c.basicAuth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, pass)))
}

func (c *Client) req(meth, uri string, data interface{}) (w *httptest.ResponseRecorder) {
	var request *http.Request
	if data != nil {
		b, _ := json.Marshal(data)
		request, _ = http.NewRequest(meth, uri, bytes.NewBuffer(b))
	} else {
		request, _ = http.NewRequest(meth, uri, nil)
	}
	request.Header.Add("accept", "application/json")
	if c.basicAuth != "" {
		request.Header.Set("authorization", "Basic "+c.basicAuth)
	}
	if c.token != "" {
		request.Header.Set("Authorization", c.token)
	}
	w = httptest.NewRecorder()
	c.engine.ServeHTTP(w, request)
	return
}

// auth
func (c *Client) Signin() *httptest.ResponseRecorder {
	return c.req("POST", "/api/v1/signin", nil)
}

func (c *Client) Signout() *httptest.ResponseRecorder {
	return c.req("POST", "/api/v1/signout", nil)
}

// device
func (c *Client) NewDevice(device interface{}) *httptest.ResponseRecorder {
	return c.req("POST", "/api/v1/device", device)
}

func (c *Client) GetDevice(deviceId int64) *httptest.ResponseRecorder {
	return c.req("GET", fmt.Sprintf("/api/v1/device/%d", deviceId), nil)
}
