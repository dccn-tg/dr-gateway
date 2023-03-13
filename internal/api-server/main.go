package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Donders-Institute/dr-gateway/internal/api-server/config"
	"github.com/Donders-Institute/dr-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/restapi"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/golang-jwt/jwt/v4"
	"github.com/s12v/go-jwks"
	"github.com/square/go-jose"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	optsPort    *int
	optsVerbose *bool
	configFile  *string
)

func init() {
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	configFile = flag.String("c", os.Getenv("DR_GATEWAY_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/dr-gateway.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

func usage() {
	fmt.Printf("\nAPI server of the DR gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", *configFile)
	}

	// initialize Cache
	ctx, cancel := context.WithCancel(context.Background())

	// collections cache
	ccache := handler.CollectionsCache{
		Config:  cfg,
		Context: ctx,
	}
	ccache.Init()

	// users cache
	ucache := handler.UsersCache{
		Config:  cfg,
		Context: ctx,
	}
	ucache.Init()

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	api := operations.NewDrGatewayAPI(swaggerSpec)
	api.UseRedoc()
	server := restapi.NewServer(api)

	// actions to take when the main program exists.
	defer func() {
		// stop all background services of the context.
		cancel()

		// stop API server.
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalf("%s", err)
		}
	}()

	server.Port = *optsPort
	server.ListenLimit = 10
	server.TLSListenLimit = 10

	// authentication with api key.
	api.APIKeyHeaderAuth = func(token string) (*models.Principal, error) {

		if token != cfg.ApiKey {
			return nil, errors.New(401, "incorrect api key auth")
		}

		// there is no user information attached, set the Principal as empty string.
		Principal := models.Principal("")
		return &Principal, nil
	}

	// authentication with username/password.
	api.BasicAuthAuth = func(username, password string) (*models.Principal, error) {

		pass, ok := cfg.Auth[username]

		if !ok || pass != password {
			return nil, errors.New(401, "incorrect username/password")
		}

		// there is login user information attached, set the Principal as the username.
		Principal := models.Principal(username)
		return &Principal, nil
	}

	// authentication with oauth2 token.
	api.Oauth2Auth = func(tokenStr string, scopes []string) (*models.Principal, error) {

		// custom claims data structure, this should match the
		// data structure expected from the authentication server.
		type IDServerClaims struct {
			Scope    []string `json:"scope"`
			Audience []string `json:"aud"`
			ClientID string   `json:"client_id"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(tokenStr, &IDServerClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New(401, "unexpected signing method: %v", token.Header["alg"])
			}

			// get public key from the auth server
			// TODO: discover jwks endpoint using oidc client.
			jwksSource := jwks.NewWebSource(cfg.JwksEndpoint)
			jwksClient := jwks.NewDefaultClient(
				jwksSource,
				time.Hour,    // Refresh keys every 1 hour
				12*time.Hour, // Expire keys after 12 hours
			)

			var jwk *jose.JSONWebKey
			jwk, err := jwksClient.GetEncryptionKey(token.Header["kid"].(string))
			if err != nil {
				return nil, errors.New(401, "cannot retrieve encryption key: %s", err)
			}

			return jwk.Key, nil
		})

		if err != nil {
			return nil, errors.New(401, "invalid token: %s", err)
		}

		// check token scope
		claims, ok := token.Claims.(*IDServerClaims)
		if !ok {
			return nil, errors.New(401, "cannot get claims from the token")
		}

		inScope := func(target string) bool {
			for _, s := range claims.Scope {
				if s == target {
					return true
				}
			}
			return false
		}

		for _, scope := range scopes {
			if !inScope(scope) {
				return nil, errors.New(401, "token not in scope: %s", scope)
			}
		}

		Principal := models.Principal(claims.ClientID)
		return &Principal, nil
	}

	// associate handler functions with implementations
	api.GetPingHandler = operations.GetPingHandlerFunc(handler.GetPing(cfg))
	api.GetMetricsHandler = operations.GetMetricsHandlerFunc(handler.GetMetrics(&ucache, &ccache))

	api.GetCollectionsHandler = operations.GetCollectionsHandlerFunc(handler.GetCollections(&ccache))
	api.GetCollectionsOuIDHandler = operations.GetCollectionsOuIDHandlerFunc(handler.GetCollectionsOfOu(&ccache))
	api.GetCollectionsProjectIDHandler = operations.GetCollectionsProjectIDHandlerFunc(handler.GetCollectionsOfProject(&ccache))

	api.GetUsersHandler = operations.GetUsersHandlerFunc(handler.GetUsers(&ucache))
	api.GetUsersOuIDHandler = operations.GetUsersOuIDHandlerFunc(handler.GetUsersOfOu(&ucache))
	api.GetUsersSearchHandler = operations.GetUsersSearchHandlerFunc(handler.SearchUsers(&ucache))

	// configure API
	server.ConfigureAPI()

	// Start API server
	if err := server.Serve(); err != nil {
		log.Fatalf("%s", err)
	}
}
