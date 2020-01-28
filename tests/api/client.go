package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	invalidToken1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	invalidToken2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODI4MjQzMTUsImlhdCI6MTU4MDE0NTkxNSwiaXNzIjoic2VydmVyIiwibmJmIjoxNTgwMTQ1OTE1LCJ1c2VySWQiOjF9.p9jcO7pu6afExwNkwF6F2y-mK3eJZOQWubcs4BhAQw2"
)

var (
	accessToken string
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

func (c *Client) GetAccessList() *httptest.ResponseRecorder {
	return c.req("GET", "/api/v1/access_list", nil)
}

// device
func (c *Client) NewDevice(device interface{}) *httptest.ResponseRecorder {
	return c.req("POST", "/api/v1/device", device)
}

func (c *Client) GetDevice(deviceId int64) *httptest.ResponseRecorder {
	return c.req("GET", fmt.Sprintf("/api/v1/device/%d", deviceId), nil)
}

func (c *Client) UpdateDevice(deviceId int64, device interface{}) *httptest.ResponseRecorder {
	return c.req("PUT", fmt.Sprintf("/api/v1/device/%d", deviceId), device)
}

func (c *Client) DeleteDevice(deviceId int64) *httptest.ResponseRecorder {
	return c.req("DELETE", fmt.Sprintf("/api/v1/device/%d", deviceId), nil)
}

func (c *Client) GetDeviceList(limit, offset int, order, sort string) *httptest.ResponseRecorder {
	uri, _ := url.Parse("/api/v1/devices")
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))
	if order != "" {
		params.Add("order", order)
	}
	if sort != "" {
		params.Add("sort_by", sort)
	}
	uri.RawQuery = params.Encode()
	return c.req("GET", uri.String(), nil)
}
