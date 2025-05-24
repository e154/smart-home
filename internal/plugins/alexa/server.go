// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"

	"github.com/e154/bus"
)

// Server ...
type Server struct {
	engine        *gin.Engine
	isStarted     *atomic.Bool
	addressPort   *string
	server        *http.Server
	skillLock     *sync.Mutex
	skills        map[string]*Skill
	adaptors      *adaptors.Adaptors
	config        Config
	scriptService scripts.ScriptService
	//gate          *client.GateClient
	eventBus bus.Bus
}

// NewServer ...
func NewServer(adaptors *adaptors.Adaptors,
	config Config,
	scriptService scripts.ScriptService,
	//gateClient *gate_client.GateClient,
	eventBus bus.Bus) *Server {
	return &Server{
		isStarted:     atomic.NewBool(false),
		adaptors:      adaptors,
		skillLock:     &sync.Mutex{},
		skills:        make(map[string]*Skill),
		config:        config,
		scriptService: scriptService,
		//gate:          gateClient,
		eventBus: eventBus,
	}
}

// Start ...
func (s *Server) Start() {

	if !s.isStarted.CompareAndSwap(false, true) {
		return
	}

	s.init()

	logger := NewLogger()

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	gin.SetMode(gin.ReleaseMode)

	s.engine = gin.New()
	s.engine.POST("/*any", s.Auth, s.handlerFunc)

	s.server = &http.Server{
		Addr:    s.config.String(),
		Handler: s.engine,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

	//todo fix
	//s.gate.SetAlexaApiEngine(s.engine)

	log.Infof("Serving server at %s", s.config.String())
}

func (s *Server) init() {

	list, err := s.adaptors.AlexaSkill.ListEnabled(context.Background(), 999, 0)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, skill := range list {
		s.AddSkill(skill)
	}
}

// Stop ...
func (s *Server) Stop() {
	if !s.isStarted.CompareAndSwap(true, false) {
		return
	}

	//todo fix
	//s.gate.SetAlexaApiEngine(nil)

	if s.server != nil {
		_ = s.server.Close()
	}
}

func (s *Server) handlerFunc(ctx *gin.Context) {

	log.Info("new request")

	req := &Request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		_ = ctx.AbortWithError(400, err)
		return
	}

	resp := NewResponse()

	switch req.GetRequestType() {
	case "LaunchRequest":
		s.OnLaunchHandler(ctx, req, resp)
	case "IntentRequest":
		s.OnIntentHandle(ctx, req, resp)
	case "SessionEndedRequest":
		s.OnSessionEndedHandler(ctx, req, resp)
	case "AudioPlayer":
		s.OnAudioPlayerHandler(ctx, req, resp)
	default:
		http.Error(ctx.Writer, "Invalid request.", http.StatusBadRequest)
	}

	ctx.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	b, _ := resp.String()
	_, _ = ctx.Writer.Write(b)
}

// OnLaunchHandler ...
func (s *Server) OnLaunchHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnLaunch(ctx, req, resp)
	}
}

// OnIntentHandle ...
func (s *Server) OnIntentHandle(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnIntent(ctx, req, resp)
	}
}

// OnSessionEndedHandler ...
func (s *Server) OnSessionEndedHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnSessionEnded(ctx, req, resp)
	}
}

// OnAudioPlayerHandler ...
func (s *Server) OnAudioPlayerHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	for _, skill := range s.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		//todo check
		//if skill.OnAudioPlayerState != nil {
		//	skill.OnAudioPlayerState(ctx, req, resp)
		//}
	}
}

// Auth ...
func (s Server) Auth(ctx *gin.Context) {

	if os.Getenv("DEV") == "true" {
		return
	}

	if !IsValidAlexaRequest(ctx.Writer, ctx.Request) {
		_ = ctx.AbortWithError(401, apperr.ErrBadRequestParams)
		return
	}
}

// AddSkill ...
func (s *Server) AddSkill(skill *m.AlexaSkill) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if _, ok := s.skills[skill.SkillId]; !ok {
		s.skills[skill.SkillId] = NewSkill(skill, s.adaptors, s.scriptService, s.eventBus)
	}
}

// UpdateSkill ...
func (s *Server) UpdateSkill(skill *m.AlexaSkill) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	s.skills[skill.SkillId] = NewSkill(skill, s.adaptors, s.scriptService, s.eventBus)
}

// DeleteSkill ...
func (s *Server) DeleteSkill(skill *m.AlexaSkill) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if _, ok := s.skills[skill.SkillId]; !ok {
		delete(s.skills, skill.SkillId)
	}
}
