package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	graceful "github.com/tylerb/graceful"

	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations"
	"strings"
	"github.com/op/go-logging"
	"fmt"
)

var (
	logServer = logging.MustGetLogger("server")
)

func configureFlags(api *operations.ServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.Logger = logServer.Infof

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.HTMLProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) (err error) {
		_, err = fmt.Fprint(w, data)
		return
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Shortcut helpers for swagger
		if r.URL.Path == "/swagger" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger/", http.StatusFound)
			return
		}

		// Server swagger config file
		if r.URL.Path == "/swagger/swagger.yml" {
			http.ServeFile(w, r, "conf/swagger/server/swagger.yml")
			return
		}

		// Serving ./swagger/
		if strings.Index(r.URL.Path, "/swagger/") == 0 {
			http.StripPrefix("/swagger/", http.FileServer(http.Dir("assets/swagger"))).ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
