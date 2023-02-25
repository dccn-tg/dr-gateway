package handler

import (
	"github.com/Donders-Institute/dr-gateway/internal/api-server/config"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func init() {

}

// GetPing returns dummy string for health check, including the authentication.
func GetPing(cfg config.Configuration) func(params operations.GetPingParams, principle *models.Principal) middleware.Responder {
	return func(params operations.GetPingParams, principle *models.Principal) middleware.Responder {
		return operations.NewGetPingOK().WithPayload("pong")
	}
}
