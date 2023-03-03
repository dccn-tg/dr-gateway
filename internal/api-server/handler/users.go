package handler

import (
	"github.com/Donders-Institute/dr-gateway/pkg/dr"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// GetUsers returns all users.
func GetUsers(ucache *UsersCache) func(params operations.GetUsersParams) middleware.Responder {
	return func(params operations.GetUsersParams) middleware.Responder {

		users := []*models.ResponseBodyUserMetadata{}

		for _, u := range ucache.GetUsers() {
			users = append(users, makeResponseBodyUserMetadata(u))
		}

		return operations.NewGetUsersOK().WithPayload(&models.ResponseBodyUsers{
			Users: users,
		})
	}
}

func makeResponseBodyUserMetadata(u *dr.DRUser) *models.ResponseBodyUserMetadata {
	return &models.ResponseBodyUserMetadata{
		DisplayName:         &u.DisplayName,
		Email:               strfmt.Email(u.Email),
		IdentityProvider:    &u.IdentityProvider,
		OrganisationalUnits: u.OrganizationalUnits,
	}
}
