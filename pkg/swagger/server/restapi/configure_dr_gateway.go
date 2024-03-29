// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/models"
	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name DrGateway --spec ../../swagger.yaml --principal models.Principal --exclude-main

func configureFlags(api *operations.DrGatewayAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DrGatewayAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.Oauth2Auth == nil {
		api.Oauth2Auth = func(token string, scopes []string) (*models.Principal, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (oauth2) has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.GetCollectionsHandler == nil {
		api.GetCollectionsHandler = operations.GetCollectionsHandlerFunc(func(params operations.GetCollectionsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetCollections has not yet been implemented")
		})
	}
	if api.GetCollectionsOuIDHandler == nil {
		api.GetCollectionsOuIDHandler = operations.GetCollectionsOuIDHandlerFunc(func(params operations.GetCollectionsOuIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetCollectionsOuID has not yet been implemented")
		})
	}
	if api.GetCollectionsProjectIDHandler == nil {
		api.GetCollectionsProjectIDHandler = operations.GetCollectionsProjectIDHandlerFunc(func(params operations.GetCollectionsProjectIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetCollectionsProjectID has not yet been implemented")
		})
	}
	if api.GetMetricsHandler == nil {
		api.GetMetricsHandler = operations.GetMetricsHandlerFunc(func(params operations.GetMetricsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetMetrics has not yet been implemented")
		})
	}
	if api.GetPingHandler == nil {
		api.GetPingHandler = operations.GetPingHandlerFunc(func(params operations.GetPingParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPing has not yet been implemented")
		})
	}
	if api.GetUsersHandler == nil {
		api.GetUsersHandler = operations.GetUsersHandlerFunc(func(params operations.GetUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUsers has not yet been implemented")
		})
	}
	if api.GetUsersOuIDHandler == nil {
		api.GetUsersOuIDHandler = operations.GetUsersOuIDHandlerFunc(func(params operations.GetUsersOuIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUsersOuID has not yet been implemented")
		})
	}
	if api.GetUsersSearchHandler == nil {
		api.GetUsersSearchHandler = operations.GetUsersSearchHandlerFunc(func(params operations.GetUsersSearchParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUsersSearch has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
