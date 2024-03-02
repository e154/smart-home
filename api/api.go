// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This libraryc is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	echopprof "github.com/hiko1129/echo-pprof"
	echoCacheMiddleware "github.com/kenshin579/echo-http-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/controllers"
	"github.com/e154/smart-home/api/stub"
	publicAssets "github.com/e154/smart-home/build"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/rbac"
)

var (
	log = logger.MustGetLogger("api")
)

// Api ...
type Api struct {
	controllers *controllers.Controllers
	echoFilter  *rbac.EchoAccessFilter
	echo        *echo.Echo
	cfg         Config
	certPublic  string
	certKey     string
	adaptors    *adaptors.Adaptors
	eventBus    bus.Bus
	httpServer  http.Server
	tlsServer   http.Server
	tlsStarted  *atomic.Bool
}

// NewApi ...
func NewApi(controllers *controllers.Controllers,
	echoFilter *rbac.EchoAccessFilter,
	cfg Config,
	eventBus bus.Bus,
	adaptors *adaptors.Adaptors) (api *Api) {
	api = &Api{
		controllers: controllers,
		echoFilter:  echoFilter,
		cfg:         cfg,
		adaptors:    adaptors,
		eventBus:    eventBus,
		tlsStarted:  atomic.NewBool(false),
	}
	return
}

// Start ...
func (a *Api) Start() (err error) {

	// HTTP
	a.echo = echo.New()
	a.echo.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Skipper: middleware.DefaultSkipper,
		Limit:   "128M",
	}))
	a.echo.Use(controllers.NewMiddlewareContextValue)
	a.echo.Use(middleware.Recover())

	if a.cfg.Debug {
		var format = `INFO	api/v1	[${method}] ${uri} ${status} ${latency_human} ${error}` + "\n"

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
		a.echo.Use(echoCacheMiddleware.CacheWithConfig(echoCacheMiddleware.CacheConfig{
			Store: echoCacheMiddleware.NewCacheMemoryStoreWithConfig(echoCacheMiddleware.CacheMemoryStoreConfig{
				Capacity:  5,
				Algorithm: echoCacheMiddleware.LFU,
			}),
			Expiration: 10 * time.Second,
		}))

	}

	a.registerHandlers()

	go a.startTlsServer()
	go a.startServer()

	a.eventBus.Subscribe("system/models/variables/+", a.eventHandler, false)
	a.eventBus.Publish("system/services/api", events.EventServiceStarted{Service: "Api"})

	return nil
}

// Shutdown ...
func (a *Api) Shutdown(ctx context.Context) (err error) {
	a.httpServer.Shutdown(ctx)
	a.tlsServer.Shutdown(ctx)
	if a.echo != nil {
		err = a.echo.Shutdown(ctx)
	}
	a.eventBus.Unsubscribe("system/models/variables/+", a.eventHandler)
	a.eventBus.Publish("system/services/api", events.EventServiceStopped{Service: "Api"})

	return
}

func (a *Api) startServer() {
	log.Infof("HTTP Server started at :%d", a.cfg.HttpPort)
	a.httpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.HttpPort),
		Handler: a.echo,
	}
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("error when starting HTTP server: %w", err)
	} else {
		log.Info("HTTP server stopped serving requests")
	}
}

func (a *Api) startTlsServer() {
	if !a.tlsStarted.CompareAndSwap(false, true) {
		return
	}
	defer a.tlsStarted.Store(false)

	a.getCerts()
	if a.certPublic == "" || a.certKey == "" {
		return
	}

	// Generate a key pair from your pem-encoded cert and key ([]byte).
	cert, err := tls.X509KeyPair([]byte(a.certPublic), []byte(a.certKey))
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("HTTPS Server started at :%d", a.cfg.HttpsPort)

	a.tlsServer = http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.HttpsPort),
		Handler: a.echo,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	if err = a.tlsServer.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		log.Errorf("error when starting HTTPS server: %w", err)
	} else {
		log.Info("HTTPS server stopped serving requests")
	}
}

// CustomMatcher ...
func (a *Api) CustomMatcher(key string) (string, bool) {
	switch key {
	case "X-Api-Key":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func (a *Api) registerHandlers() {

	// Swagger
	if a.cfg.Swagger {
		var contentHandler = echo.WrapHandler(http.FileServer(http.FS(SwaggerAssets)))
		a.echo.GET("/swagger-ui", contentHandler)
		a.echo.GET("/swagger-ui/*", contentHandler)
		a.echo.GET("/api.swagger3.yaml", contentHandler)
	}

	var typedocHandler = echo.WrapHandler(http.FileServer(http.FS(TypedocAssets)))
	a.echo.GET("/typedoc", typedocHandler)
	a.echo.GET("/typedoc/*", typedocHandler)

	wrapper := stub.ServerInterfaceWrapper{
		Handler: a.controllers,
	}

	v1 := a.echo.Group("/v1")
	v1.GET("/access_list", a.echoFilter.Auth(wrapper.AuthServiceAccessList))
	v1.POST("/action", a.echoFilter.Auth(wrapper.ActionServiceAddAction))
	v1.DELETE("/action/:id", a.echoFilter.Auth(wrapper.ActionServiceDeleteAction))
	v1.GET("/action/:id", a.echoFilter.Auth(wrapper.ActionServiceGetActionById))
	v1.PUT("/action/:id", a.echoFilter.Auth(wrapper.ActionServiceUpdateAction))
	v1.GET("/actions", a.echoFilter.Auth(wrapper.ActionServiceGetActionList))
	v1.GET("/actions/search", a.echoFilter.Auth(wrapper.ActionServiceSearchAction))
	v1.POST("/area", a.echoFilter.Auth(wrapper.AreaServiceAddArea))
	v1.DELETE("/area/:id", a.echoFilter.Auth(wrapper.AreaServiceDeleteArea))
	v1.GET("/area/:id", a.echoFilter.Auth(wrapper.AreaServiceGetAreaById))
	v1.PUT("/area/:id", a.echoFilter.Auth(wrapper.AreaServiceUpdateArea))
	v1.GET("/areas", a.echoFilter.Auth(wrapper.AreaServiceGetAreaList))
	v1.GET("/areas/search", a.echoFilter.Auth(wrapper.AreaServiceSearchArea))
	v1.GET("/backups", a.echoFilter.Auth(wrapper.BackupServiceGetBackupList))
	v1.POST("/backups", a.echoFilter.Auth(wrapper.BackupServiceNewBackup))
	v1.POST("/backup/upload", a.echoFilter.Auth(wrapper.BackupServiceUploadBackup))
	v1.POST("/backup/apply", a.echoFilter.Auth(wrapper.BackupServiceApplyState))
	v1.POST("/backup/rollback", a.echoFilter.Auth(wrapper.BackupServiceRevertState))
	v1.PUT("/backup/:name", a.echoFilter.Auth(wrapper.BackupServiceRestoreBackup))
	v1.DELETE("/backup/:name", a.echoFilter.Auth(wrapper.BackupServiceDeleteBackup))
	v1.POST("/condition", a.echoFilter.Auth(wrapper.ConditionServiceAddCondition))
	v1.DELETE("/condition/:id", a.echoFilter.Auth(wrapper.ConditionServiceDeleteCondition))
	v1.GET("/condition/:id", a.echoFilter.Auth(wrapper.ConditionServiceGetConditionById))
	v1.PUT("/condition/:id", a.echoFilter.Auth(wrapper.ConditionServiceUpdateCondition))
	v1.GET("/conditions", a.echoFilter.Auth(wrapper.ConditionServiceGetConditionList))
	v1.GET("/conditions/search", a.echoFilter.Auth(wrapper.ConditionServiceSearchCondition))
	v1.POST("/dashboard", a.echoFilter.Auth(wrapper.DashboardServiceAddDashboard))
	v1.DELETE("/dashboard/:id", a.echoFilter.Auth(wrapper.DashboardServiceDeleteDashboard))
	v1.GET("/dashboard/:id", a.echoFilter.Auth(wrapper.DashboardServiceGetDashboardById))
	v1.PUT("/dashboard/:id", a.echoFilter.Auth(wrapper.DashboardServiceUpdateDashboard))
	v1.POST("/dashboard_card", a.echoFilter.Auth(wrapper.DashboardCardServiceAddDashboardCard))
	v1.POST("/dashboard_card/import", a.echoFilter.Auth(wrapper.DashboardCardServiceImportDashboardCard))
	v1.DELETE("/dashboard_card/:id", a.echoFilter.Auth(wrapper.DashboardCardServiceDeleteDashboardCard))
	v1.GET("/dashboard_card/:id", a.echoFilter.Auth(wrapper.DashboardCardServiceGetDashboardCardById))
	v1.PUT("/dashboard_card/:id", a.echoFilter.Auth(wrapper.DashboardCardServiceUpdateDashboardCard))
	v1.POST("/dashboard_card_item", a.echoFilter.Auth(wrapper.DashboardCardItemServiceAddDashboardCardItem))
	v1.DELETE("/dashboard_card_item/:id", a.echoFilter.Auth(wrapper.DashboardCardItemServiceDeleteDashboardCardItem))
	v1.GET("/dashboard_card_item/:id", a.echoFilter.Auth(wrapper.DashboardCardItemServiceGetDashboardCardItemById))
	v1.PUT("/dashboard_card_item/:id", a.echoFilter.Auth(wrapper.DashboardCardItemServiceUpdateDashboardCardItem))
	v1.GET("/dashboard_card_items", a.echoFilter.Auth(wrapper.DashboardCardItemServiceGetDashboardCardItemList))
	v1.GET("/dashboard_cards", a.echoFilter.Auth(wrapper.DashboardCardServiceGetDashboardCardList))
	v1.POST("/dashboard_tab", a.echoFilter.Auth(wrapper.DashboardTabServiceAddDashboardTab))
	v1.DELETE("/dashboard_tab/:id", a.echoFilter.Auth(wrapper.DashboardTabServiceDeleteDashboardTab))
	v1.GET("/dashboard_tab/:id", a.echoFilter.Auth(wrapper.DashboardTabServiceGetDashboardTabById))
	v1.PUT("/dashboard_tab/:id", a.echoFilter.Auth(wrapper.DashboardTabServiceUpdateDashboardTab))
	v1.POST("/dashboard_tabs/import", a.echoFilter.Auth(wrapper.DashboardTabServiceImportDashboardTab))
	v1.GET("/dashboard_tabs", a.echoFilter.Auth(wrapper.DashboardTabServiceGetDashboardTabList))
	v1.GET("/dashboards", a.echoFilter.Auth(wrapper.DashboardServiceGetDashboardList))
	v1.POST("/dashboards/import", a.echoFilter.Auth(wrapper.DashboardServiceImportDashboard))
	v1.GET("/dashboards/search", a.echoFilter.Auth(wrapper.DashboardServiceSearchDashboard))
	v1.POST("/developer_tools/automation/call_action", a.echoFilter.Auth(wrapper.DeveloperToolsServiceCallAction))
	v1.POST("/developer_tools/automation/call_trigger", a.echoFilter.Auth(wrapper.DeveloperToolsServiceCallTrigger))
	v1.GET("/developer_tools/bus/state", a.echoFilter.Auth(wrapper.DeveloperToolsServiceGetEventBusStateList))
	v1.POST("/developer_tools/entity/reload", a.echoFilter.Auth(wrapper.DeveloperToolsServiceReloadEntity))
	v1.POST("/developer_tools/entity/set_state", a.echoFilter.Auth(wrapper.DeveloperToolsServiceEntitySetState))
	v1.GET("/entities", a.echoFilter.Auth(wrapper.EntityServiceGetEntityList))
	v1.POST("/entities/import", a.echoFilter.Auth(wrapper.EntityServiceImportEntity))
	v1.POST("/entity", a.echoFilter.Auth(wrapper.EntityServiceAddEntity))
	v1.GET("/entity/search", a.echoFilter.Auth(wrapper.EntityServiceSearchEntity))
	v1.DELETE("/entity/:id", a.echoFilter.Auth(wrapper.EntityServiceDeleteEntity))
	v1.GET("/entity/:id", a.echoFilter.Auth(wrapper.EntityServiceGetEntity))
	v1.PUT("/entity/:id", a.echoFilter.Auth(wrapper.EntityServiceUpdateEntity))
	v1.POST("/entity/:id/disable", a.echoFilter.Auth(wrapper.EntityServiceDisabledEntity))
	v1.POST("/entity/:id/enable", a.echoFilter.Auth(wrapper.EntityServiceEnabledEntity))
	v1.GET("/entity_storage", a.echoFilter.Auth(wrapper.EntityStorageServiceGetEntityStorageList))
	v1.GET("/entities/statistic", a.echoFilter.Auth(wrapper.EntityServiceGetStatistic))
	v1.POST("/image", a.echoFilter.Auth(wrapper.ImageServiceAddImage))
	v1.POST("/image/upload", a.echoFilter.Auth(wrapper.ImageServiceUploadImage))
	v1.DELETE("/image/:id", a.echoFilter.Auth(wrapper.ImageServiceDeleteImageById))
	v1.GET("/image/:id", a.echoFilter.Auth(wrapper.ImageServiceGetImageById))
	v1.PUT("/image/:id", a.echoFilter.Auth(wrapper.ImageServiceUpdateImageById))
	v1.GET("/images", a.echoFilter.Auth(wrapper.ImageServiceGetImageList))
	v1.GET("/images/filter_list", a.echoFilter.Auth(wrapper.ImageServiceGetImageFilterList))
	v1.GET("/images/filtered", a.echoFilter.Auth(wrapper.ImageServiceGetImageListByDate))
	v1.POST("/interact/entity/call_action", a.echoFilter.Auth(wrapper.InteractServiceEntityCallAction))
	v1.GET("/logs", a.echoFilter.Auth(wrapper.LogServiceGetLogList))
	v1.GET("/message_delivery", a.echoFilter.Auth(wrapper.MessageDeliveryServiceGetMessageDeliveryList))
	v1.GET("/metric", a.echoFilter.Auth(wrapper.MetricServiceGetMetric))
	v1.GET("/mqtt/client/:id", a.echoFilter.Auth(wrapper.MqttServiceGetClientById))
	v1.GET("/mqtt/clients", a.echoFilter.Auth(wrapper.MqttServiceGetClientList))
	v1.GET("/mqtt/subscriptions", a.echoFilter.Auth(wrapper.MqttServiceGetSubscriptionList))
	v1.POST("/password_reset", a.echoFilter.Auth(wrapper.AuthServicePasswordReset))
	v1.GET("/plugin/:name", a.echoFilter.Auth(wrapper.PluginServiceGetPlugin))
	v1.POST("/plugin/:name/disable", a.echoFilter.Auth(wrapper.PluginServiceDisablePlugin))
	v1.POST("/plugin/:name/enable", a.echoFilter.Auth(wrapper.PluginServiceEnablePlugin))
	v1.PUT("/plugin/:name/settings", a.echoFilter.Auth(wrapper.PluginServiceUpdatePluginSettings))
	v1.GET("/plugins", a.echoFilter.Auth(wrapper.PluginServiceGetPluginList))
	v1.GET("/plugins/search", a.echoFilter.Auth(wrapper.PluginServiceSearchPlugin))
	v1.GET("/plugin/:name/readme", a.echoFilter.Auth(wrapper.PluginServiceGetPluginReadme))
	v1.POST("/role", a.echoFilter.Auth(wrapper.RoleServiceAddRole))
	v1.DELETE("/role/:name", a.echoFilter.Auth(wrapper.RoleServiceDeleteRoleByName))
	v1.GET("/role/:name", a.echoFilter.Auth(wrapper.RoleServiceGetRoleByName))
	v1.PUT("/role/:name", a.echoFilter.Auth(wrapper.RoleServiceUpdateRoleByName))
	v1.GET("/role/:name/access_list", a.echoFilter.Auth(wrapper.RoleServiceGetRoleAccessList))
	v1.PUT("/role/:name/access_list", a.echoFilter.Auth(wrapper.RoleServiceUpdateRoleAccessList))
	v1.GET("/roles", a.echoFilter.Auth(wrapper.RoleServiceGetRoleList))
	v1.GET("/roles/search", a.echoFilter.Auth(wrapper.RoleServiceSearchRoleByName))
	v1.POST("/script", a.echoFilter.Auth(wrapper.ScriptServiceAddScript))
	v1.POST("/script/exec_src", a.echoFilter.Auth(wrapper.ScriptServiceExecSrcScriptById))
	v1.DELETE("/script/:id", a.echoFilter.Auth(wrapper.ScriptServiceDeleteScriptById))
	v1.GET("/script/:id", a.echoFilter.Auth(wrapper.ScriptServiceGetScriptById))
	v1.GET("/script/:id/compiled", a.echoFilter.Auth(wrapper.ScriptServiceGetCompiledScriptById))
	v1.PUT("/script/:id", a.echoFilter.Auth(wrapper.ScriptServiceUpdateScriptById))
	v1.POST("/script/:id/copy", a.echoFilter.Auth(wrapper.ScriptServiceCopyScriptById))
	v1.POST("/script/:id/exec", a.echoFilter.Auth(wrapper.ScriptServiceExecScriptById))
	v1.GET("/scripts", a.echoFilter.Auth(wrapper.ScriptServiceGetScriptList))
	v1.GET("/scripts/search", a.echoFilter.Auth(wrapper.ScriptServiceSearchScript))
	v1.GET("/scripts/statistic", a.echoFilter.Auth(wrapper.ScriptServiceGetStatistic))
	v1.GET("/tags/search", a.echoFilter.Auth(wrapper.TagServiceSearchTag))
	v1.GET("/tags", a.echoFilter.Auth(wrapper.TagServiceGetTagList))
	v1.DELETE("/tag/:id", a.echoFilter.Auth(wrapper.TagServiceDeleteTagById))
	v1.GET("/tag/:id", a.echoFilter.Auth(wrapper.TagServiceGetTagById))
	v1.PUT("/tag/:id", a.echoFilter.Auth(wrapper.TagServiceUpdateTagById))
	v1.POST("/signin", wrapper.AuthServiceSignin)
	v1.POST("/signout", a.echoFilter.Auth(wrapper.AuthServiceSignout))
	v1.POST("/task", a.echoFilter.Auth(wrapper.AutomationServiceAddTask))
	v1.DELETE("/task/:id", a.echoFilter.Auth(wrapper.AutomationServiceDeleteTask))
	v1.GET("/task/:id", a.echoFilter.Auth(wrapper.AutomationServiceGetTask))
	v1.PUT("/task/:id", a.echoFilter.Auth(wrapper.AutomationServiceUpdateTask))
	v1.POST("/task/:id/disable", a.echoFilter.Auth(wrapper.AutomationServiceDisableTask))
	v1.POST("/task/:id/enable", a.echoFilter.Auth(wrapper.AutomationServiceEnableTask))
	v1.GET("/tasks", a.echoFilter.Auth(wrapper.AutomationServiceGetTaskList))
	v1.POST("/tasks/import", a.echoFilter.Auth(wrapper.AutomationServiceImportTask))
	v1.POST("/trigger", a.echoFilter.Auth(wrapper.TriggerServiceAddTrigger))
	v1.DELETE("/trigger/:id", a.echoFilter.Auth(wrapper.TriggerServiceDeleteTrigger))
	v1.GET("/trigger/:id", a.echoFilter.Auth(wrapper.TriggerServiceGetTriggerById))
	v1.PUT("/trigger/:id", a.echoFilter.Auth(wrapper.TriggerServiceUpdateTrigger))
	v1.GET("/triggers", a.echoFilter.Auth(wrapper.TriggerServiceGetTriggerList))
	v1.GET("/triggers/search", a.echoFilter.Auth(wrapper.TriggerServiceSearchTrigger))
	v1.POST("/triggers/:id/disable", a.echoFilter.Auth(wrapper.TriggerServiceDisableTrigger))
	v1.POST("/triggers/:id/enable", a.echoFilter.Auth(wrapper.TriggerServiceEnableTrigger))
	v1.POST("/user", a.echoFilter.Auth(wrapper.UserServiceAddUser))
	v1.DELETE("/user/:id", a.echoFilter.Auth(wrapper.UserServiceDeleteUserById))
	v1.GET("/user/:id", a.echoFilter.Auth(wrapper.UserServiceGetUserById))
	v1.PUT("/user/:id", a.echoFilter.Auth(wrapper.UserServiceUpdateUserById))
	v1.GET("/users", a.echoFilter.Auth(wrapper.UserServiceGetUserList))
	v1.POST("/variable", a.echoFilter.Auth(wrapper.VariableServiceAddVariable))
	v1.DELETE("/variable/:name", a.echoFilter.Auth(wrapper.VariableServiceDeleteVariable))
	v1.GET("/variable/:name", a.echoFilter.Auth(wrapper.VariableServiceGetVariableByName))
	v1.PUT("/variable/:name", a.echoFilter.Auth(wrapper.VariableServiceUpdateVariable))
	v1.GET("/variables", a.echoFilter.Auth(wrapper.VariableServiceGetVariableList))
	v1.GET("/variables/search", a.echoFilter.Auth(wrapper.VariableServiceSearchVariable))
	v1.GET("/zigbee2mqtt/bridge", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceGetBridgeList))
	v1.POST("/zigbee2mqtt/bridge", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceAddZigbee2mqttBridge))
	v1.DELETE("/zigbee2mqtt/bridge/:id", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceDeleteBridgeById))
	v1.GET("/zigbee2mqtt/bridge/:id", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceGetZigbee2mqttBridge))
	v1.PUT("/zigbee2mqtt/bridge/:id/bridge", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceUpdateBridgeById))
	v1.GET("/zigbee2mqtt/bridge/:id/devices", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceDeviceList))
	v1.GET("/zigbee2mqtt/bridge/:id/networkmap", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceNetworkmap))
	v1.POST("/zigbee2mqtt/bridge/:id/networkmap", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceUpdateNetworkmap))
	v1.POST("/zigbee2mqtt/bridge/:id/reset", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceResetBridgeById))
	v1.POST("/zigbee2mqtt/device_ban", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceDeviceBan))
	v1.POST("/zigbee2mqtt/device_rename", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceDeviceRename))
	v1.POST("/zigbee2mqtt/device_whitelist", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceDeviceWhitelist))
	v1.GET("/zigbee2mqtt/search_device", a.echoFilter.Auth(wrapper.Zigbee2mqttServiceSearchDevice))
	v1.GET("/ws", a.echoFilter.Auth(wrapper.StreamServiceSubscribe))

	// static files
	a.echo.GET("/", echo.WrapHandler(a.controllers.Index(publicAssets.F)))
	a.echo.GET("/*", echo.WrapHandler(http.FileServer(http.FS(publicAssets.F))))
	a.echo.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(publicAssets.F))))
	fileServer := http.FileServer(http.Dir("./data/file_storage"))
	a.echo.Any("/upload/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/upload/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		fileServer.ServeHTTP(w, r)
	})))
	staticServer := http.FileServer(http.Dir("./data/static"))
	a.echo.Any("/static/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/static/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		staticServer.ServeHTTP(w, r)
	})))
	snapshotServer := http.FileServer(http.Dir("./snapshots"))
	a.echo.GET("/snapshots/*", a.echoFilter.Auth(echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/snapshots/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		snapshotServer.ServeHTTP(w, r)
	}))))
	// webdav
	webdav := echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//r.RequestURI = strings.ReplaceAll(r.RequestURI, "/webdav/", "/")
		//r.URL, _ = r.URL.Parse(r.RequestURI)
		a.controllers.Webdav(w, r)
	}))
	a.echo.Any("/webdav", webdav)
	a.echo.Any("/webdav/*", webdav)

	// media
	a.echo.Any("/stream/:entity_id/channel/:channel/mse", a.echoFilter.Auth(a.controllers.StreamMSE)) //Auth
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/init.mp4", a.controllers.Media.StreamHLSLLInit)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/index.m3u8", a.controllers.Media.StreamHLSLLM3U8)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/segment/:segment/:any", a.controllers.Media.StreamHLSLLM4Segment)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/fragment/:segment/:fragment/:any", a.controllers.Media.StreamHLSLLM4Fragment)

	// Cors
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodHead},
	}))

}

func (a *Api) Echo() *echo.Echo {
	return a.echo
}

func (a *Api) getCerts() {
	certPublicVar, err := a.adaptors.Variable.GetByName(context.Background(), "certPublic")
	if err != nil {
		return
	}
	a.certPublic = strings.TrimSpace(certPublicVar.Value)
	certKeyVar, err := a.adaptors.Variable.GetByName(context.Background(), "certKey")
	if err != nil {
		return
	}
	a.certKey = strings.TrimSpace(certKeyVar.Value)
}

func (a *Api) eventHandler(_ string, message interface{}) {
	switch v := message.(type) {
	case events.EventUpdatedVariableModel:
		switch v.Name {
		case "certPublic", "certKey":
			log.Infof("updated settings name %s", v.Name)
			a.tlsServer.Shutdown(context.Background())
			a.startTlsServer()
		}
	}
}
