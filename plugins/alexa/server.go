// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"net/http"
	"os"
	"sync"
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
	gate          *gate_client.GateClient
	eventBus      event_bus.EventBus
}

// NewServer ...
func NewServer(adaptors *adaptors.Adaptors,
	config Config,
	scriptService scripts.ScriptService,
	gateClient *gate_client.GateClient,
	eventBus event_bus.EventBus) *Server {
	return &Server{
		isStarted:     atomic.NewBool(false),
		adaptors:      adaptors,
		skillLock:     &sync.Mutex{},
		skills:        make(map[string]*Skill),
		config:        config,
		scriptService: scriptService,
		gate:          gateClient,
		eventBus:      eventBus,
	}
}

// Start ...
func (s *Server) Start() {

	if s.isStarted.Load() {
		return
	}
	s.isStarted.Store(true)

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

	s.gate.SetAlexaApiEngine(s.engine)

	log.Infof("Serving server at %s", s.config.String())
}

func (s *Server) init() {

	list, err := s.adaptors.AlexaSkill.ListEnabled(999, 0)
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
	if !s.isStarted.Load() {
		return
	}
	s.isStarted.Store(false)

	s.gate.SetAlexaApiEngine(nil)

	if s.server != nil {
		s.server.Close()
	}
}

func (s *Server) handlerFunc(ctx *gin.Context) {

	log.Info("new request")

	req := &Request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		ctx.AbortWithError(400, err)
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
	ctx.Writer.Write(b)
}

func (s *Server) OnLaunchHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnLaunch(ctx, req, resp)
	}
}

func (s *Server) OnIntentHandle(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnIntent(ctx, req, resp)
	}
}

func (s *Server) OnSessionEndedHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	if skill, ok := s.skills[req.Context.System.Application.ApplicationID]; ok {
		skill.OnSessionEnded(ctx, req, resp)
	}
}

func (s *Server) OnAudioPlayerHandler(ctx *gin.Context, req *Request, resp *Response) {
	s.skillLock.Lock()
	defer s.skillLock.Unlock()

	for _, skill := range s.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		if skill.OnAudioPlayerState != nil {
			skill.OnAudioPlayerState(ctx, req, resp)
		}
	}
}

// Auth ...
func (s Server) Auth(ctx *gin.Context) {

	if os.Getenv("DEV") == "true" {
		return
	}

	if !IsValidAlexaRequest(ctx.Writer, ctx.Request) {
		ctx.AbortWithError(401, errors.New("invalid request"))
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