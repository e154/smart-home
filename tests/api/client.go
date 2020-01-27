package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

type Client struct {
	engine *gin.Engine
	token  string
}

func NewClient(engine *gin.Engine) *Client {
	return &Client{
		engine: engine,
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Signin(login, pass string) (w *httptest.ResponseRecorder) {
	request, _ := http.NewRequest("POST", "/api/v1/signin", nil)
	request.Header.Add("accept", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, pass)))
	request.Header.Set("authorization", "Basic "+auth)
	w = httptest.NewRecorder()
	c.engine.ServeHTTP(w, request)
	return
}

func (c *Client) Signout() (w *httptest.ResponseRecorder) {
	request, _ := http.NewRequest("POST", "/api/v1/signout", nil)
	request.Header.Add("accept", "application/json")
	request.Header.Set("Authorization", c.token)
	w = httptest.NewRecorder()
	c.engine.ServeHTTP(w, request)
	return
}
