// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package server

import (
	"context"
	"net/http"

	echopprof "github.com/hiko1129/echo-pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/e154/smart-home/api"
	publicAssets "github.com/e154/smart-home/build"
	"github.com/e154/smart-home/system/gate/server/controllers"
	"github.com/e154/smart-home/system/gate/server/wsp"
)

type Server struct {
	controllers *controllers.Controllers
	echo        *echo.Echo
	proxy       *wsp.Server
	cfg         *Config
}

func NewServer(cfg *Config, proxy *wsp.Server) *Server {
	return &Server{
		controllers: controllers.NewControllers(cfg.FullAddress(), "gate"),
		proxy:       proxy,
		cfg:         cfg,
	}
}

// Start ...
func (a *Server) Start() (err error) {

	// HTTP
	a.echo = echo.New()
	a.echo.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Skipper: middleware.DefaultSkipper,
		Limit:   "128M",
	}))
	a.echo.Use(middleware.Recover())

	if a.cfg.Debug {
		var format = `INFO	gate	[${method}] ${uri} ${status} ${latency_human} ${error}` + "\n"

		log.Info("debug enabled")
		DefaultLoggerConfig := middleware.LoggerConfig{
			Skipper:          middleware.DefaultSkipper,
			Format:           format,
			CustomTimeFormat: "2006-01-02 15:04:05.00000",
		}
		a.echo.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
		a.echo.Debug = true
	}

	if a.cfg.Pprof {
		// automatically add routers for net/http/pprof
		// e.g. /debug/pprof, /debug/pprof/heap, etc.
		log.Info("pprof enabled")
		echopprof.Wrap(a.echo)

		prefix := "/debug/pprof"
		group := a.echo.Group(prefix)
		echopprof.WrapGroup(prefix, group)
	}

	a.echo.HideBanner = true
	a.echo.HidePort = true

	if a.cfg.Gzip {
		a.echo.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
		a.echo.Use(middleware.Decompress())
	}

	a.registerHandlers()

	var https bool
	if a.cfg.Https {
		if a.cfg.Domain == "" {
			log.Warnf("domain settings is empty")
		} else {
			https = true
			a.echo.Pre(middleware.HTTPSRedirect())
			// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
			a.echo.AutoTLSManager.Cache = autocert.DirCache("./conf")
			a.echo.AutoTLSManager.HostPolicy = autocert.HostWhitelist(a.cfg.Domain)
		}
	}

	go func() {
		var err error
		if https {
			err = a.echo.StartAutoTLS(a.cfg.String())
		} else {
			err = a.echo.Start(a.cfg.String())
		}
		if err.Error() != "http: Server closed" {
			log.Error(err.Error())
		}
	}()

	log.Infof("server started at %s", a.cfg.String())

	return
}

// Shutdown ...
func (a *Server) Shutdown(ctx context.Context) (err error) {

	if a.echo != nil {
		err = a.echo.Shutdown(ctx)
	}

	return
}

func (a *Server) registerHandlers() {

	// static files
	a.echo.GET("/", echo.WrapHandler(a.controllers.Index(publicAssets.F)))
	a.echo.GET("/public/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(publicAssets.F)))))
	var contentHandler = echo.WrapHandler(http.FileServer(http.FS(api.SwaggerAssets)))
	a.echo.GET("/swagger-ui", contentHandler)
	a.echo.GET("/swagger-ui/*", contentHandler)
	a.echo.GET("/api.swagger3.yaml", contentHandler)

	// proxy
	a.echo.Any("/v1/*", a.proxyHandler)
	a.echo.Any("/upload/*", a.proxyHandler)
	a.echo.Any("/static/*", a.proxyHandler)
	a.echo.Any("/snapshots/*", a.proxyHandler)
	a.echo.GET("/v1/ws", func(c echo.Context) error {
		a.proxy.Ws(c.Response(), c.Request())
		return nil
	})

	// internal
	a.echo.GET("/gate/register", func(c echo.Context) error {
		a.proxy.Register(c.Response(), c.Request())
		return nil
	})

	// Cors
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodHead},
	}))
}

func (a *Server) proxyHandler(c echo.Context) error {
	a.proxy.Request(c.Response(), c.Request())
	return nil
}
