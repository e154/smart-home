// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package alexa

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"net/http"
	"strings"
	"sync"
)

var (
	log = common.MustGetLogger("alexa")
)

type Alexa struct {
	engine      *gin.Engine
	isStarted   *atomic.Bool
	addressPort *string
	server      *http.Server
	appLock     *sync.Mutex
	apps        []Application
	adaptors    *adaptors.Adaptors
	appConfig   *config.AppConfig
	token       *atomic.String
}

func NewAlexa(adaptors *adaptors.Adaptors,
	appConfig *config.AppConfig) *Alexa {
	return &Alexa{
		isStarted: atomic.NewBool(false),
		adaptors:  adaptors,
		appLock:   &sync.Mutex{},
		apps:      make([]Application, 0),
		appConfig: appConfig,
		token:     atomic.NewString(""),
	}
}

func (a *Alexa) Start() {

	if a.isStarted.Load() {
		return
	}
	a.isStarted.Store(true)

	a.init()

	a.engine = gin.New()
	a.engine.POST("/", a.handlerFunc)

	a.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "127.0.0.1", "3033"),
		Handler: a.engine,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()
}

func (a *Alexa) init() {
	a.appLock.Lock()
	defer a.appLock.Unlock()

	if err := a.getSettings(); err != nil {
		log.Error(err.Error())
	}

	// hello world
	a.apps = append(a.apps, NewHelloWorld("amzn1.ask.skill.e2177905-951c-4e39-b061-0b5cac49bdd9"))
}

func (a *Alexa) Stop() {
	if !a.isStarted.Load() {
		return
	}
	a.isStarted.Store(false)

	if a.server != nil {
		a.server.Close()
	}
}

func (a *Alexa) handlerFunc(ctx *gin.Context) {

	log.Info("new request")

	req := &Request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		ctx.AbortWithError(400, err)
		return
	}

	resp := NewResponse()

	if req.GetRequestType() == "LaunchRequest" {
		fmt.Println("LaunchRequest")
		a.onLaunchHandler(ctx, req, resp)
	} else if req.GetRequestType() == "IntentRequest" {
		fmt.Println("IntentRequest")
		a.onIntentHandle(ctx, req, resp)
	} else if req.GetRequestType() == "SessionEndedRequest" {
		a.onSessionEndedHandler(ctx, req, resp)
	} else if strings.HasPrefix(req.GetRequestType(), "AudioPlayer.") {
		a.onAudioPlayerHandler(ctx, req, resp)
	} else {
		http.Error(ctx.Writer, "Invalid request.", http.StatusBadRequest)
	}

	ctx.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	b, _ := resp.String()
	ctx.Writer.Write(b)
}

func (a *Alexa) onLaunchHandler(ctx *gin.Context, req *Request, resp *Response) {

	for _, app := range a.apps {
		if app.GetAppID() != req.Context.System.Application.ApplicationID {
			return
		}
		if app.OnLaunch != nil {
			app.OnLaunch(ctx, req, resp)
		}
	}
}

func (a *Alexa) onIntentHandle(ctx *gin.Context, req *Request, resp *Response) {

	for _, app := range a.apps {
		if app.GetAppID() != req.Context.System.Application.ApplicationID {
			return
		}
		if app.OnIntent != nil {
			app.OnIntent(ctx, req, resp)
		}
	}
}

func (a *Alexa) onSessionEndedHandler(ctx *gin.Context, req *Request, resp *Response) {

	for _, app := range a.apps {
		if app.GetAppID() != req.Context.System.Application.ApplicationID {
			return
		}
		if app.OnSessionEnded != nil {
			app.OnSessionEnded(ctx, req, resp)
		}
	}
}

func (a *Alexa) onAudioPlayerHandler(ctx *gin.Context, req *Request, resp *Response) {

	for _, app := range a.apps {
		if app.GetAppID() != req.Context.System.Application.ApplicationID {
			return
		}
		if app.OnAudioPlayerState != nil {
			app.OnAudioPlayerState(ctx, req, resp)
		}
	}
}

func (a *Alexa) genAccessToken() {
	a.token.Store(uuid.NewV4().String())
}

func (a *Alexa) getSettings() (err error) {

	var variable *m.Variable
	if variable, err = a.adaptors.Variable.GetByName("alexa_token"); err != nil || variable == nil {
		a.genAccessToken()
		variable = &m.Variable{
			Name:     "alexa_token",
			Value:    a.token.Load(),
			Autoload: false,
		}
		err = a.adaptors.Variable.Add(variable)
	}
	a.token.Store(variable.Value)
	return
}

func (a *Alexa) updateSettings() (err error) {
	err = a.adaptors.Variable.Update(&m.Variable{
		Name:     "alexa_token",
		Value:    a.token.Load(),
		Autoload: false,
	})
	return
}

func (a Alexa) Auth(ctx *gin.Context) {

	accessToken := ctx.Request.URL.Query().Get("token")

	if accessToken == "" || accessToken != a.token.Load() {
		ctx.AbortWithError(401, errors.New("access token invalid"))
		return
	}
	if !IsValidAlexaRequest(ctx.Writer, ctx.Request) {
		ctx.AbortWithError(401, errors.New("invalid request"))
		return
	}
}
