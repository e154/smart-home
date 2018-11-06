package server

import (
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/api/server_v1/stub/restapi"
	"github.com/e154/smart-home/api/server_v1/controllers"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/index"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/node"
)

var (
	log = logging.MustGetLogger("server")
)

type Server struct {
	RestApiServer *restapi.Server
	Config        *ServerConfig
	Controllers   *controllers.Controllers
	api           *operations.ServerAPI
}

func (s *Server) Start() {
	if err := s.RestApiServer.Serve(); err != nil {
		log.Error(err.Error())
		return
	}
}

func (s *Server) Shutdown() {
	if s.RestApiServer != nil {
		s.RestApiServer.Shutdown()
	}
}

func NewServer(cfg *ServerConfig, ctrls *controllers.Controllers) (newServer *Server) {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	api := operations.NewServerAPI(swaggerSpec)

	newServer = &Server{
		Config:      cfg,
		Controllers: ctrls,
		api:         api,
	}

	newServer.setControllers()

	server := restapi.NewServer(api)

	//defer server.Shutdown()
	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Server API"
	parser.LongDescription = "Server API"

	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	server.ConfigureAPI()

	server.Host = cfg.Host
	server.Port = cfg.Port

	newServer.RestApiServer = server

	return
}

func (s *Server) setControllers() {

	// index
	s.api.IndexIndexHandler = index.IndexHandlerFunc(s.Controllers.Index.ControllerIndex)

	// node
	s.api.NodeAddNodeHandler = node.AddNodeHandlerFunc(s.Controllers.Node.AddNode)
	s.api.NodeGetNodeByIDHandler = node.GetNodeByIDHandlerFunc(s.Controllers.Node.GetNodeById)
	s.api.NodeGetNodeListHandler = node.GetNodeListHandlerFunc(s.Controllers.Node.GetNodeList)
	s.api.NodeUpdateNodeHandler = node.UpdateNodeHandlerFunc(s.Controllers.Node.UpdateNode)
	s.api.NodeDeleteNodeByIDHandler = node.DeleteNodeByIDHandlerFunc(s.Controllers.Node.DeleteNodeById)

	//
}
